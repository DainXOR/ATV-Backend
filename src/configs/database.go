package configs

import (
	"context"
	"dainxor/atv/logger"
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

type db struct {
	gormDB  *gorm.DB
	mongoDB *mongo.Client

	dbType string

	user string
	pass string
	name string
	host string
	port string
}

var DB db
var DataBase *gorm.DB

func (db) Get() *gorm.DB {
	return DataBase
}

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

func (db) EnvInit() {
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
}

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
func (db) ConnectPostgres(host string, user string, password string, dbname string, port string) {
	var err error
	dsn := "host=" + host +
		" user=" + user +
		" password=" + password +
		" dbname=" + dbname +
		" port=" + port +
		" sslmode=disable"
	logger.Info("Connecting to database: ", dsn)
	DataBase, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		logger.Fatal(err)
	}
}

func (db) ConnectSQLiteEnv() {
	dbname, exist := os.LookupEnv("DB_NAME")
	if exist {
		DB.ConnectSQLite(dbname)
	} else {
		DB.ConnectSQLite("atvsqlite.db")
	}
}
func (db) ConnectSQLite(dbname string) {
	var err error
	DataBase, err = gorm.Open(sqlite.Open("atvsqlite.db"), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}
}

func (db) ConnectMongoDBEnv() {
	useTesting, exist := os.LookupEnv("DB_TESTING")
	if exist && useTesting != "TRUE" {
		DB.ConnectMongoDB(os.Getenv("DB_PORT"))
	} else {
		logger.Info("Using testing database")
		DB.ConnectMongoDB(os.Getenv("DB_PORT_TEST"))
	}
}
func (db) ConnectMongoDB(port string) {
	// Create a Client to a MongoDB server and use Ping to verify that the
	// server is running.

	clientOpts := options.Client().ApplyURI("mongodb://localhost:27017")
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

func (db) CreateDatabase() {
	logger.Info("Creating database")

	switch DB.dbType {
	case "POSTGRES":
		logger.Info("Creating Postgres database")
		DB.CreatePostgresDatabase()

	case "MONGODB":
		logger.Info("Creating MongoDB database")
		DB.CreateMongoDBDatabase()

	case "SQLITE":
		fallthrough
	default:
		logger.Info("Creating SQLite database")
		DB.CreateSQLiteDatabase()
	}
}

func (db) CreatePostgresDatabase() error {
	var exists bool
	checkQuery := "SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = ?)"
	if err := DB.Get().Raw(checkQuery, DB.name).Scan(&exists).Error; err != nil {
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
	if err := DB.Get().Exec(createQuery).Error; err != nil {
		return fmt.Errorf("failed to create database '%s': %w", DB.name, err)
	}

	log.Printf("Database '%s' created successfully", DB.name)
	return nil
}

func (db) CreateMongoDBDatabase() {
}

func (db) CreateSQLiteDatabase() {
}

func (db) Migrate(models ...any) {
	logger.Info("Starting migrations")

	for _, model := range models {
		err := DataBase.AutoMigrate(model)

		if err != nil {
			logger.Error("Error migrating model: ", err)
			logger.Error("Model: ", utils.StructToString(model))
			logger.Fatal("Migration failed")
		}
	}

	logger.Info("Migrations completed")
}
