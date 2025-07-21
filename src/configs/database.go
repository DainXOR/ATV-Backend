package configs

import (
	"cmp"
	"context"
	"dainxor/atv/logger"
	"dainxor/atv/models"
	"dainxor/atv/types"
	"dainxor/atv/utils"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"

	"fmt"
	"os"
)

type mongoType struct {
	client        *mongo.Client
	db            *mongo.Database
	disconectFunc func()
}
type gormType struct {
	db *gorm.DB
}
type dbTypes struct {
}

func (dbTypes) Postgres() string {
	return "POSTGRES"
}
func (dbTypes) MongoDB() string {
	return "MONGO"
}
func (dbTypes) SQLite() string {
	return "SQLITE"
}
func (dbTypes) Default() string {
	return DB.Types().SQLite()
}

type db struct {
	dbType string

	dbName           string
	connectionString string
}

var DB db
var mongoT mongoType
var gormT gormType

func (db) Type() string {
	return DB.dbType
}
func (db) Types() dbTypes {
	return dbTypes{}
}

func (db) Gorm() *gormType {
	return &gormT
}
func (db) Mongo() *mongoType {
	return &mongoT
}
func (mongoType) Disconnect() {
	if DB.Mongo().disconectFunc != nil {
		DB.Mongo().disconectFunc()
	} else {
		logger.Warning("MongoDB disconect function is nil, nothing to do")
	}
}

func init() {
	DB.envInit()
}
func ReloadDBEnv() {
	DB.Close()
	DB.envInit()
}
func (db) envInit() error {
	var exist bool
	DB.dbType, exist = os.LookupEnv("DB_TYPE")
	if !exist {
		logger.Warning("DB_TYPE not found, using default: ", DB.Types().Default())
		DB.dbType = DB.Types().Default()
	}

	DB.loadDBConfig()
	DB.connectDB()
	return DB.CreateDatabase()
}

func (db) GormDB() *gorm.DB {
	return DB.Gorm().db
}

func (db) MongoDB() *mongo.Database {
	return DB.Mongo().db
}
func (db) In(name string) *mongo.Collection {
	return DB.MongoDB().Collection(name)
}
func (db) From(v models.DBModelInterface) *mongo.Collection {
	// Use reflection to get the collection name from the struct
	collectionName := v.TableName()
	return DB.In(collectionName)
}
func (db) Context() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}

func (db) FindOne(filter any, result models.DBModelInterface) error {
	ctx, cancel := DB.Context()
	defer cancel()

	return DB.From(result).FindOne(ctx, filter).Decode(result)
}
func (db) FindAll(filter any, result any) error {
	logger.Lava("0.1.1", "This mf should be refactored to use []models.DBModelInterface instead of any for the result")
	ctx, cancel := DB.Context()
	defer cancel()

	eType, err := utils.SliceType(result)
	if err != nil {
		logger.Fatal("This function ONLY works with slices or pointers to slices")
	}

	iType, ok := eType.(models.DBModelInterface)
	if !ok {
		logger.Fatal("Result type does NOT IMPLEMENT TableName method")
	}

	cursor, err := DB.From(iType).Find(ctx, filter)
	if err != nil {
		logger.Error("Failed to find documents:", err)
		return err
	}
	defer cursor.Close(ctx)

	return cursor.All(ctx, result)
}

func (db) InsertOne(document models.DBModelInterface) (*mongo.InsertOneResult, error) {
	ctx, cancel := DB.Context()
	defer cancel()

	return DB.From(document).InsertOne(ctx, document)
}
func (db) UpdateOne(filter any, update any, result models.DBModelInterface) types.Result[mongo.UpdateResult] {
	ctx, cancel := DB.Context()
	defer cancel()

	updateResult, err := DB.From(result).UpdateOne(ctx, filter, update)
	if err != nil {
		logger.Error("Failed to update document:", err)
		return types.ResultErr[mongo.UpdateResult](err)
	}

	err = DB.FindOne(filter, result)
	if err != nil {
		logger.Error("Failed to find updated document:", err)
		return types.ResultErr[mongo.UpdateResult](err)
	}
	return types.ResultOk(*updateResult)
}
func (db) PatchOne(filter any, update any, result models.DBModelInterface) types.Result[mongo.UpdateResult] {
	ctx, cancel := DB.Context()
	defer cancel()

	updateResult, err := DB.From(result).UpdateOne(ctx, filter, update)
	if err != nil {
		logger.Error("Failed to update document:", err)
		return types.ResultErr[mongo.UpdateResult](err)
	}

	err = DB.FindOne(filter, result)

	return types.ResultOf(*updateResult, err, err != nil)
}

// LoadDBConfig loads the database configuration from environment variables
// and sets the default values if not found. It also sets the database type.
func (db) loadDBConfig() {
	useTesting, exist := os.LookupEnv("DB_TESTING")

	if exist || useTesting == "TRUE" {
		logger.Debug("Using testing database")
		DB.connectionString = os.Getenv("CONECTION_STRING_TEST")
		DB.dbName = os.Getenv("DB_NAME_TEST")
	} else {
		logger.Debug("Using production database")
		DB.connectionString = os.Getenv("CONECTION_STRING")
		DB.dbName = os.Getenv("DB_NAME")
	}
}
func (db) connectDB() {
	switch DB.Type() {
	case DB.Types().Postgres():
		logger.Debug("Using Postgres database")
		DB.ConnectPostgresEnv()

	case DB.Types().MongoDB():
		logger.Debug("Using MongoDB database")
		DB.ConnectMongoDBEnv()

	case DB.Types().SQLite():
		logger.Debug("Using SQLite database")
		DB.ConnectSQLiteEnv()
	default:
		logger.Debug("Using default database:", DB.Types().Default())
		DB.dbType = DB.Types().Default()
	}

	logger.Debug("Database connection established")
}

// ConnectPostgresEnv connects to the Postgres database using environment variables
// It checks for the testing environment and uses the appropriate database credentials
func (db) ConnectPostgresEnv() {
	DB.ConnectPostgres(DB.connectionString)
}

// ConnectPostgres connects to the Postgres database using the provided credentials
// It uses the gorm library to establish the connection
func (db) ConnectPostgres(connectionString string) {
	var err error
	dsn := connectionString

	logger.Debug("Connecting to database: ", dsn)
	DB.Gorm().db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		logger.Fatal(err)
	}
}

// ConnectSQLiteEnv connects to the SQLite database using environment variables
// It checks for the database name in the environment variables and uses a default if not found
func (db) ConnectSQLiteEnv() {
	DB.ConnectSQLite(cmp.Or(DB.dbName, "atvsqlite.db"))
}

// ConnectSQLite connects to the SQLite database using the provided database name
// It uses the gorm library to establish the connection
func (db) ConnectSQLite(dbname string) {
	var err error
	DB.Gorm().db, err = gorm.Open(sqlite.Open(dbname), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}
}

// ConnectMongoDBEnv connects to the MongoDB database using environment variables
// It checks for the testing environment and uses the appropriate database credentials
func (db) ConnectMongoDBEnv() {
	DB.ConnectMongoDB(DB.dbName, DB.connectionString)
}

// Connects to the MongoDB database
func (db) ConnectMongoDB(dbName, conectionString string) {
	logger.Debug("Connecting to MongoDB: ", conectionString)
	clientOpts := options.Client().ApplyURI(conectionString)
	client, err := mongo.Connect(clientOpts)
	if err != nil {
		logger.Fatal(err)
	}

	ctx, cancel := DB.Context()
	defer cancel()
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		logger.Fatal(err)
	}

	DB.Mongo().client = client
	DB.Mongo().db = client.Database(dbName)

	DB.Mongo().disconectFunc = func() {
		if err = DB.Mongo().client.Disconnect(context.Background()); err != nil {
			logger.Fatal("Error disconnecting from MongoDB: ", err)
		}
	}
}

// CreateDatabase creates the database based on the database type
// It checks if the database already exists and creates it if not
func (db) CreateDatabase() error {
	logger.Debug("Creating database")

	switch DB.Type() {
	case DB.Types().Postgres():
		logger.Debug("Creating Postgres database")
		err := DB.CreatePostgresDatabase()

		if err != nil {
			logger.Error("Error creating Postgres database: ", err)
			return err
		}

	case DB.Types().MongoDB():
		logger.Debug("Creating MongoDB database")
		DB.CreateMongoDBDatabase()

	case DB.Types().SQLite():
		logger.Debug("Creating SQLite database")
		DB.CreateSQLiteDatabase()
	}

	return nil
}

// CreatePostgresDatabase creates the Postgres database if it does not exist
// It uses a raw SQL query to check for the existence of the database
func (db) CreatePostgresDatabase() error {
	var exists bool
	checkQuery := "SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = ?)"
	if err := DB.GormDB().Raw(checkQuery, DB.dbName).Scan(&exists).Error; err != nil {
		logger.Error("Error checking database existence: ", err)
		return err
	}

	if exists {
		logger.Debug("Database", DB.dbName, "already exists")
		return nil
	}

	// Create the new database
	createQuery := fmt.Sprintf("CREATE DATABASE \"%s\"", DB.dbName)
	if err := DB.GormDB().Exec(createQuery).Error; err != nil {
		logger.Error("failed to create database", DB.dbName, ": ", err)
		return fmt.Errorf("failed to create database '%s': %w", DB.dbName, err)
	}

	logger.Debug("Database", DB.dbName, "created successfully")
	return nil
}

// CreateMongoDBDatabase creates the MongoDB database if it does not exist
// MongoDB does not require explicit database creation, it is created on first use
func (db) CreateMongoDBDatabase() {
	logger.Info("MongoDB database created automatically")
}

// CreateSQLiteDatabase creates the SQLite database if it does not exist
// SQLite databases are created automatically when you open a connection to a non-existent database file
func (db) CreateSQLiteDatabase() {
	// SQLite databases are created automatically when you open a connection to a non-existent database file.
	// So, no explicit creation is needed.
	logger.Info("SQLite database created automatically")
}

// Migrate performs database migrations for the provided models
// It uses the gorm library to automatically migrate the models to the database
func (db) Migrate(models ...any) {
	logger.Info("Starting migrations")

	if DB.dbType == DB.Types().Postgres() {
		for _, model := range models {
			err := DB.GormDB().AutoMigrate(model)

			if err != nil {
				logger.Error("Error migrating model: ", err)
				logger.Error("Model: ", utils.StructToString(model))
				logger.Fatal("Migration failed")
			}
		}

		logger.Info("Migrations completed")

	} else {
		logger.Info("No migrations needed for SQLite or MongoDB")
	}
}

// Close closes the database connection
func closeGormDB() {
	sqlDB, err := DB.GormDB().DB()
	if err != nil {
		logger.Error("Error getting SQL DB: ", err)
		return
	}
	sqlDB.Close()
	logger.Info("GormDB connection closed")
}
func (db) ClosePostgres() {
	closeGormDB()
}
func (db) CloseSQLite() {
	closeGormDB()
}
func (db) CloseMongoDB() {
	DB.Mongo().Disconnect()
	logger.Info("MongoDB connection closed")
}

func (db) Close() {
	switch DB.Type() {
	case DB.Types().Postgres():
		DB.ClosePostgres()

	case DB.Types().MongoDB():
		DB.CloseMongoDB()

	case DB.Types().SQLite():
		DB.CloseSQLite()

	default:
		logger.Warning("Unknown database type, no specific close method available")
	}
}
