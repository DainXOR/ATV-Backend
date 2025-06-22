package configs

import (
	"cmp"
	"context"
	"dainxor/atv/logger"
	"dainxor/atv/models"
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

	connectionString string

	user string
	pass string
	name string
	host string
	port string
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

func (gormType) DB() *gorm.DB {
	return DB.Gorm().db
}

func (mongoType) DB() *mongo.Database {
	return DB.Mongo().db
}
func (mongoType) In(name string) *mongo.Collection {
	return DB.Mongo().DB().Collection(name)
}
func (mongoType) From(v models.DBModelInterface) *mongo.Collection {
	// Use reflection to get the collection name from the struct
	collectionName := v.TableName()
	return DB.Mongo().In(collectionName)
}
func (mongoType) Context() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 2*time.Second)
}
func (mongoType) Disconnect() {
	if DB.Mongo().disconectFunc != nil {
		DB.Mongo().disconectFunc()
	} else {
		logger.Warning("MongoDB disconect function is nil, nothing to do")
	}
}

func (mongoType) FindOne(filter any, result models.DBModelInterface) error {
	ctx, cancel := DB.Mongo().Context()
	defer cancel()

	return DB.Mongo().From(result).FindOne(ctx, filter).Decode(result)
}
func (mongoType) FindAll(filter any, result any) error {
	logger.Lava(1, "This mf should be refactored to use models.DBModelInterface instead of any in the result")
	ctx, cancel := DB.Mongo().Context()
	defer cancel()

	eType, err := utils.SliceType(result)
	if err != nil {
		logger.Fatal("This function ONLY works with slices or pointers to slices")
	}

	iType, ok := eType.(interface{ TableName() string })
	if !ok {
		logger.Fatal("Result type does NOT IMPLEMENT TableName method")
	}

	cursor, err := DB.Mongo().From(iType).Find(ctx, filter)
	if err != nil {
		logger.Error("Failed to find documents:", err)
		return err
	}
	defer cursor.Close(ctx)

	return cursor.All(ctx, result)
}

func (mongoType) InsertOne(document models.DBModelInterface) (*mongo.InsertOneResult, error) {
	ctx, cancel := DB.Mongo().Context()
	defer cancel()

	return DB.Mongo().From(document).InsertOne(ctx, document)
}

func (mongoType) PatchOne(filter any, update any, result models.DBModelInterface) error {
	ctx, cancel := DB.Mongo().Context()
	defer cancel()

	updateResult, err := DB.Mongo().From(result).UpdateOne(ctx, filter, update)
	if err != nil {
		logger.Error("Failed to update document:", err)
		return err
	}
	if updateResult.MatchedCount == 0 {
		logger.Warning("No documents matched the filter for update")
		return nil
	}
	return DB.Mongo().FindOne(filter, result)
}

// LoadDBConfig loads the database configuration from environment variables
// and sets the default values if not found. It also sets the database type.
func (db) loadDBConfig() {
	useTesting, exist := os.LookupEnv("DB_TESTING")
	logger.Lava(1, "Remove individual database environment variables.")
	envUser := "DB_USER"
	envPass := "DB_PASSWORD"
	envName := "DB_NAME"
	envHost := "DB_HOST"
	envPort := "DB_PORT"

	if exist || useTesting == "TRUE" {
		logger.Debug("Using testing database")
		envUser += "_TEST"
		envPass += "_TEST"
		envName += "_TEST"
		envHost += "_TEST"
		envPort += "_TEST"

		DB.connectionString = os.Getenv("CONECTION_STRING_TEST")
	} else {
		logger.Debug("Using production database")
		DB.connectionString = os.Getenv("CONECTION_STRING")
	}

	dbType, exist := os.LookupEnv("DB_TYPE")
	if exist {
		DB.dbType = dbType
	} else {
		logger.Warning("DB_TYPE not found, using default: ", DB.Types().Default())
		DB.dbType = DB.Types().Default()
	}

	v := logger.Lava(1, "Loading database configuration from individual environment variables")
	v.LavaStart()
	host, exist := os.LookupEnv(envHost)
	if exist {
		DB.host = host
	} else {
		logger.Warning(envHost, "not found, using default")
		DB.host = "localhost"
	}

	user, exist := os.LookupEnv(envUser)
	if exist {
		DB.user = user
	} else {
		logger.Warning(envUser, "not found, using default")
		DB.user = "postgres"
	}

	pass, exist := os.LookupEnv(envPass)
	if exist {
		DB.pass = pass
	} else {
		logger.Warning(envPass, "not found, using default")
		DB.pass = ""
	}

	name, exist := os.LookupEnv(envName)
	if exist {
		DB.name = name
	} else {
		logger.Warning(envName, "not found, using default")
		DB.name = "atv-test"
	}

	port, exist := os.LookupEnv(envPort)
	if exist {
		DB.port = port
	} else {
		logger.Warning(envPort, "not found, using default")
		DB.port = "5432"
	}
	v.LavaEnd()
}

// EnvInit initializes the database connection based on the environment variables
// It checks for the database type and connects to the appropriate database
func (db) EnvInit() error {
	dbType, exist := os.LookupEnv("DB_TYPE")
	DB.dbType = dbType
	DB.loadDBConfig()

	logger.Debug("Use default database: ", !exist)

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
		logger.Debug("Using default database")
		DB.dbType = DB.Types().Default()
	}

	logger.Debug("Database connection established")
	logger.Info("Database type: ", DB.Type())

	return DB.CreateDatabase()
}

// ConnectPostgresEnv connects to the Postgres database using environment variables
// It checks for the testing environment and uses the appropriate database credentials
func (db) ConnectPostgresEnv() {
	useTesting, exist := os.LookupEnv("DB_TESTING")
	v := logger.Lava(1, "Using individual database environment variables")
	v.LavaStart()
	if exist && useTesting != "TRUE" {
		DB.ConnectPostgres(
			os.Getenv("DB_HOST"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"),
			os.Getenv("DB_PORT"),
		)
	} else {
		logger.Info("Using testing database")
		DB.ConnectPostgres(
			os.Getenv("DB_HOST_TEST"),
			os.Getenv("DB_USER_TEST"),
			os.Getenv("DB_PASSWORD_TEST"),
			os.Getenv("DB_NAME_TEST"),
			os.Getenv("DB_PORT_TEST"),
		)
	}
	v.LavaEnd()
}

// ConnectPostgres connects to the Postgres database using the provided credentials
// It uses the gorm library to establish the connection
func (db) ConnectPostgres(host string, user string, password string, dbname string, port string) {
	var err error
	v := logger.Lava(1, "Creating Postgres connection string, use the environment variable CONNECTION_STRING instead")
	v.LavaStart()
	dsn := "host=" + host +
		" user=" + user +
		" password=" + password +
		" dbname=" + dbname +
		" port=" + port +
		" sslmode=disable"

	v.LavaEnd()
	logger.Info("Connecting to database: ", dsn)
	DB.Gorm().db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		logger.Fatal(err)
	}
}

// ConnectSQLiteEnv connects to the SQLite database using environment variables
// It checks for the database name in the environment variables and uses a default if not found
func (db) ConnectSQLiteEnv() {
	dbname := os.Getenv("DB_NAME")
	DB.ConnectSQLite(cmp.Or(dbname, "atvsqlite.db"))
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
	logger.Lava(1, "Refactor getting MongoDB connection string to be here")
	DB.ConnectMongoDB(DB.connectionString)
}

// ConnectMongoDB connects to the MongoDB database using the provided port
// It uses the mongo driver to establish the connection
func (db) ConnectMongoDB(conectionString string) {
	// Create a Client to a MongoDB server and use Ping to verify that the
	// server is running.

	logger.Debug("Connecting to MongoDB: ", conectionString)
	clientOpts := options.Client().ApplyURI(conectionString)
	client, err := mongo.Connect(clientOpts)

	if err != nil {
		logger.Fatal(err)
	}

	logger.Lava(1, "Clean up this code")
	// Call Ping to verify that the deployment is up and the Client was
	// configured successfully. As mentioned in the Ping documentation, this
	// reduces application resiliency as the server may be temporarily
	// unavailable when Ping is called.
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		logger.Fatal(err)
	}

	DB.Mongo().client = client
	DB.Mongo().db = client.Database(DB.name)

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
	if err := DB.Gorm().DB().Raw(checkQuery, DB.name).Scan(&exists).Error; err != nil {
		logger.Error("Error checking database existence: ", err)
		return err
	}

	if exists {
		logger.Debug("Database", DB.name, "already exists")
		return nil
	}

	// Create the new database
	createQuery := fmt.Sprintf("CREATE DATABASE \"%s\"", DB.name)
	if err := DB.Gorm().DB().Exec(createQuery).Error; err != nil {
		logger.Error("failed to create database", DB.name, ": ", err)
		return fmt.Errorf("failed to create database '%s': %w", DB.name, err)
	}

	logger.Debug("Database", DB.name, "created successfully")
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
			err := DB.Gorm().DB().AutoMigrate(model)

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
func (db) ClosePostgres() {
	sqlDB, err := DB.Gorm().DB().DB()
	if err != nil {
		logger.Error("Error getting SQL DB: ", err)
		return
	}
	sqlDB.Close()
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
		logger.Info("No need to close SQLite connection")

	default:
		logger.Warning("Unknown database type, no specific close method available")
	}
}
