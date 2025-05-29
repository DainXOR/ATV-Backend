package configs

import (
	"context"
	"dainxor/atv/logger"
	"dainxor/atv/types"
	"dainxor/atv/utils"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type mongoType struct {
	client *mongo.Client
	db     *mongo.Database
}
type gormType struct {
	db *gorm.DB
}
type db struct {
	dbType string

	conectionString string

	user string
	pass string
	name string
	host string
	port string
}

var DB db
var DataBase *gorm.DB
var mongoT mongoType
var gormT gormType

func (db) Gorm() *gormType {
	return &gormT
}
func (db) Mongo() *mongoType {
	return &mongoT
}

func (gormType) DB() *gorm.DB {
	return gormT.db
}

func (mongoType) DB() *mongo.Database {
	return mongoT.db
}
func (mongoType) Collection(name string) *mongo.Collection {
	return mongoT.db.Collection(name)
}
func (mongoType) Context() context.Context {
	return context.TODO()
}

func (db) GetFirst(dest any, id string) types.HttpError {
	if DB.Type() == "MONGODB" {
		collectionName := dest.(interface{ TableName() string }).TableName()
		err := DB.Mongo().DB().Collection(collectionName).FindOne(DB.Mongo().Context(), map[string]any{"_id": id}).Decode(&dest)
		ret := types.HttpError{}
		ret.Err = err
		return ret
	}
	if DB.Type() == "POSTGRES" || DB.Type() == "SQLITE" {
		DB.Gorm().DB().First(&dest, id)
		return types.HttpError{}
	}

	return types.HttpError{}
}

// LoadDBConfig loads the database configuration from environment variables
// and sets the default values if not found. It also sets the database type.
func (db) loadDBConfig() {
	useTesting, exist := os.LookupEnv("DB_TESTING")
	envUser := "DB_USER"
	envPass := "DB_PASSWORD"
	envName := "DB_NAME"
	envHost := "DB_HOST"
	envPort := "DB_PORT"

	if exist || useTesting == "TRUE" {
		logger.Info("Using testing database")
		envUser += "_TEST"
		envPass += "_TEST"
		envName += "_TEST"
		envHost += "_TEST"
		envPort += "_TEST"
	}

	dbType, exist := os.LookupEnv("DB_TYPE")
	if exist {
		DB.dbType = dbType
	} else {
		logger.Warning("DB_TYPE not found, using default")
		DB.dbType = "POSTGRES"
	}

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
}

// EnvInit initializes the database connection based on the environment variables
// It checks for the database type and connects to the appropriate database
func (db) EnvInit() error {
	dbType, exist := os.LookupEnv("DB_TYPE")
	DB.dbType = dbType
	DB.loadDBConfig()

	logger.Info("Use default database: ", !exist)

	switch DB.dbType {
	case "POSTGRES":
		logger.Info("Using Postgres database")
		DB.ConnectPostgresEnv()

	case "MONGODB":
		logger.Info("Using MongoDB database")
		DB.ConnectMongoDBEnv()

	case "SQLITE":
		fallthrough
	default:
		logger.Info("Using SQLite database")
		DB.dbType = "SQLITE"
		DB.ConnectSQLiteEnv()
	}

	logger.Info("Database connection established")
	logger.Info("Database type: ", DB.dbType)

	return DB.CreateDatabase()
}

// ConnectPostgresEnv connects to the Postgres database using environment variables
// It checks for the testing environment and uses the appropriate database credentials
func (db) ConnectPostgresEnv() {
	useTesting, exist := os.LookupEnv("DB_TESTING")
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
}

// ConnectPostgres connects to the Postgres database using the provided credentials
// It uses the gorm library to establish the connection
func (db) ConnectPostgres(host string, user string, password string, dbname string, port string) {
	var err error
	dsn := "host=" + host +
		" user=" + user +
		" password=" + password +
		" dbname=" + dbname +
		" port=" + port +
		" sslmode=disable"
	logger.Info("Connecting to database: ", dsn)
	gormT.db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		logger.Fatal(err)
	}
}

// ConnectSQLiteEnv connects to the SQLite database using environment variables
// It checks for the database name in the environment variables and uses a default if not found
func (db) ConnectSQLiteEnv() {
	dbname, exist := os.LookupEnv("DB_NAME")
	if exist {
		DB.ConnectSQLite(dbname)
	} else {
		DB.ConnectSQLite("atvsqlite.db")
	}
}

// ConnectSQLite connects to the SQLite database using the provided database name
// It uses the gorm library to establish the connection
func (db) ConnectSQLite(dbname string) {
	var err error
	gormT.db, err = gorm.Open(sqlite.Open("atvsqlite.db"), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}
}

// ConnectMongoDBEnv connects to the MongoDB database using environment variables
// It checks for the testing environment and uses the appropriate database credentials
func (db) ConnectMongoDBEnv() {
	useTesting, exist := os.LookupEnv("DB_TESTING")
	if exist && useTesting != "TRUE" {
		DB.ConnectMongoDB(os.Getenv("CONECTION_STRING"))
	} else {
		logger.Info("Using testing database")
		DB.ConnectMongoDB(os.Getenv("CONECTION_STRING_TEST"))
	}
}

// ConnectMongoDB connects to the MongoDB database using the provided port
// It uses the mongo driver to establish the connection
func (db) ConnectMongoDB(conectionString string) {
	// Create a Client to a MongoDB server and use Ping to verify that the
	// server is running.

	clientOpts := options.Client().ApplyURI(conectionString)

	client, err := mongo.Connect(clientOpts)
	if err != nil {
		logger.Fatal(err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			logger.Fatal(err)
		}
	}()

	// Call Ping to verify that the deployment is up and the Client was
	// configured successfully. As mentioned in the Ping documentation, this
	// reduces application resiliency as the server may be temporarily
	// unavailable when Ping is called.
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		logger.Fatal(err)
	}
}

// CreateDatabase creates the database based on the database type
// It checks if the database already exists and creates it if not
func (db) CreateDatabase() error {
	logger.Info("Creating database")

	switch DB.dbType {
	case "POSTGRES":
		logger.Info("Creating Postgres database")
		err := DB.CreatePostgresDatabase()

		if err != nil {
			logger.Error("Error creating Postgres database: ", err)
			return err
		}

	case "MONGODB":
		logger.Info("Creating MongoDB database")
		DB.CreateMongoDBDatabase()

	case "SQLITE":
		fallthrough
	default:
		logger.Info("Creating SQLite database")
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
		log.Printf("Database '%s' already exists", DB.name)
		logger.Info("Database", DB.name, "already exists")
		return nil
	}

	// Create the new database
	createQuery := fmt.Sprintf("CREATE DATABASE \"%s\"", DB.name)
	if err := DB.Gorm().DB().Exec(createQuery).Error; err != nil {
		return fmt.Errorf("failed to create database '%s': %w", DB.name, err)
	}

	log.Printf("Database '%s' created successfully", DB.name)
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

	if DB.dbType == "POSTGRES" {
		for _, model := range models {
			err := gormT.DB().AutoMigrate(model)

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
func (db) Close() {
	if DB.dbType == "POSTGRES" {
		sqlDB, err := gormT.DB().DB()
		if err != nil {
			logger.Error("Error getting SQL DB: ", err)
			return
		}
		sqlDB.Close()
	} else if DB.dbType == "MONGODB" {
		if err := mongoT.DB().Client().Disconnect(context.Background()); err != nil {
			logger.Error("Error disconnecting from MongoDB: ", err)
		}
	} else {
		logger.Info("No need to close SQLite connection")
	}
}

// GetDBType returns the database type in use
func (db) Type() string {
	return DB.dbType
}
