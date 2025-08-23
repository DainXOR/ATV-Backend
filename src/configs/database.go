package configs

import (
	"dainxor/atv/logger"
	"dainxor/atv/models"
	"dainxor/atv/types"
	"errors"

	"os"
)

type dbType struct {
	accessor         InterfaceDBAccessor
	name             string
	connectionString string
}

var DB dbType

func init() {
	DB.LoadEnv()
}

// LoadDBConfig loads the database configuration from environment variables
// and sets the default values if not found. It also sets the database type.
func (dbType) LoadEnv() {
	useTesting, exist := os.LookupEnv("DB_TESTING")

	if exist || useTesting == "TRUE" {
		logger.Debug("Using testing database")
		DB.connectionString = os.Getenv("CONECTION_STRING_TEST")
		DB.name = os.Getenv("DB_NAME_TEST")
	} else {
		logger.Debug("Using production database")
		DB.connectionString = os.Getenv("CONECTION_STRING")
		DB.name = os.Getenv("DB_NAME")
	}

	if DB.connectionString == "" {
		logger.Error("Database connection string is not set")
	}
	if DB.name == "" {
		logger.Error("Database name is not set")
	}

}
func (dbType) ReloadConnection() {
	DB.Close()
	DB.LoadEnv()
	DB.Start()
}

func (dbType) Use(db InterfaceDBAccessor) *dbType {
	DB.accessor = db
	return &DB
}
func (dbType) Start() error {
	if DB.accessor == nil {
		logger.Error("Database accessor is not set")
		return errors.New("Database accessor is not set")
	}

	if DB.connectionString == "" || DB.name == "" {
		logger.Error("Database connection string or name is not set")
		return errors.New("Database connection string or name is not set")
	}

	return DB.accessor.Connect(DB.name, DB.connectionString)
}
func (dbType) ConnectTo(dbName, connectionString string) error {
	if DB.accessor == nil {
		logger.Error("Database accessor is not set")
		return errors.New("Database accessor is not set")
	}

	return DB.accessor.Connect(dbName, connectionString)
}

func (dbType) InsertOne(document models.DBModelInterface) types.Result[models.DBID] {
	return DB.accessor.InsertOne(document)
}
func (dbType) InsertMany(documents ...models.DBModelInterface) types.Result[[]models.DBID] {
	return DB.accessor.InsertMany(documents...)
}

func (dbType) FindOne(filter any, result models.DBModelInterface) types.Result[models.DBModelInterface] {
	return DB.accessor.FindOne(filter, result)
}
func (dbType) FindAll(filter any, result models.DBModelInterface) types.Result[[]models.DBModelInterface] {
	return DB.accessor.FindMany(filter, result)
}

func (dbType) UpdateOne(filter any, update models.DBModelInterface) error {
	return DB.accessor.UpdateOne(filter, update)
}
func (dbType) UpdateMany(filter any, update models.DBModelInterface) error {
	return DB.accessor.UpdateMany(filter, update)
}

func (dbType) PatchOne(filter any, update models.DBModelInterface) error {
	return DB.accessor.PatchOne(filter, update)
}
func (dbType) PatchMany(filter any, update models.DBModelInterface) error {
	return DB.accessor.PatchMany(filter, update)
}

func (dbType) SoftDeleteOne(filter any, model models.DBModelInterface) error {
	return DB.accessor.SoftDeleteOne(filter, model)
}
func (dbType) SoftDeleteMany(filter any, model models.DBModelInterface) error {
	return DB.accessor.SoftDeleteMany(filter, model)
}

func (dbType) PermanentDeleteOne(filter any, model models.DBModelInterface) error {
	return DB.accessor.PermanentDeleteOne(filter, model)
}
func (dbType) PermanentDeleteMany(filter any, model models.DBModelInterface) error {
	return DB.accessor.PermanentDeleteMany(filter, model)
}

// Migrate performs database migrations for the provided models
// It uses the gorm library to automatically migrate the models to the database
func (dbType) Migrate(models ...models.DBModelInterface) {
	logger.Info("Starting migrations")

	if err := DB.accessor.Migrate(models...); err != nil {
		logger.Error("Migration failed:", err)
		return
	}

	logger.Info("Migrations completed")
}

func (dbType) Close() {
	DB.accessor.Disconnect()
}
