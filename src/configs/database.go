package configs

import (
	"dainxor/atv/logger"
	"dainxor/atv/models"
	"dainxor/atv/types"
	"dainxor/atv/utils"
	"errors"

	"go.mongodb.org/mongo-driver/v2/mongo"

	"os"
)

type dbType struct {
	dbAccessor       InterfaceDB
	dbName           string
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
		DB.dbName = os.Getenv("DB_NAME_TEST")
	} else {
		logger.Debug("Using production database")
		DB.connectionString = os.Getenv("CONECTION_STRING")
		DB.dbName = os.Getenv("DB_NAME")
	}

	if DB.connectionString == "" {
		logger.Error("Database connection string is not set")
	}
	if DB.dbName == "" {
		logger.Error("Database name is not set")
	}

}
func (dbType) ReloadConnection() {
	DB.Close()
	DB.LoadEnv()
	DB.Start()
}

func (dbType) Use(db InterfaceDB) *dbType {
	DB.dbAccessor = db
	return &DB
}
func (dbType) Start() error {
	if DB.dbAccessor == nil {
		logger.Error("Database accessor is not set")
		return errors.New("Database accessor is not set")
	}

	if DB.connectionString == "" || DB.dbName == "" {
		logger.Error("Database connection string or name is not set")
		return errors.New("Database connection string or name is not set")
	}

	return DB.dbAccessor.Connect(DB.dbName, DB.connectionString)
}
func (dbType) ConnectTo(dbName, connectionString string) error {
	if DB.dbAccessor == nil {
		logger.Error("Database accessor is not set")
		return errors.New("Database accessor is not set")
	}

	return DB.dbAccessor.Connect(dbName, connectionString)
}

func (dbType) FindOne(filter any, result models.DBModelInterface) error {
	return DB.dbAccessor.GetOne(filter, result).Error()
}
func (dbType) FindAll(filter any, result any) error {
	logger.Lava(types.V("0.1.1"), "This mf should be refactored to use []models.DBModelInterface instead of any for the result")
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

func (dbType) InsertOne(document models.DBModelInterface) (*mongo.InsertOneResult, error) {
	ctx, cancel := DB.Context()
	defer cancel()

	return DB.From(document).InsertOne(ctx, document)
}
func (dbType) UpdateOne(filter any, update any, result models.DBModelInterface) types.Result[mongo.UpdateResult] {
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
func (dbType) PatchOne(filter any, update any, result models.DBModelInterface) types.Result[mongo.UpdateResult] {
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

// Migrate performs database migrations for the provided models
// It uses the gorm library to automatically migrate the models to the database
func (dbType) Migrate(models ...models.DBModelInterface) {
	logger.Info("Starting migrations")

	if err := DB.dbAccessor.Migrate(models...); err != nil {
		logger.Error("Migration failed:", err)
		return
	}

	logger.Info("Migrations completed")
}

func (dbType) Close() {
	DB.dbAccessor.Disconnect()
}
