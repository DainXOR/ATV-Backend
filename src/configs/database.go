package configs

import (
	"dainxor/atv/logger"
	"dainxor/atv/models"
	"dainxor/atv/types"
	"errors"

	"os"
)

type dbNS struct {
	dbType           string
	accessor         InterfaceDBAccessor
	name             string
	connectionString string
}

var DB dbNS
var accessors = make(map[string]InterfaceDBAccessor)

func init() {
	DB.DefineAccessor("MONGO", NewMongoAccessor())

	DB.LoadEnv()
}

// DefineAccessor registers a new database accessor for a given database type
func (dbNS) DefineAccessor(dbType string, accessor InterfaceDBAccessor) error {
	if dbType == "" {
		logger.Error("Database type is empty")
		return errors.New("Database type is empty")
	}
	if accessor == nil {
		logger.Error("Database accessor is nil")
		return errors.New("Database accessor is nil")
	}

	logger.Debug("Defining database accessor for type:", dbType)
	accessors[dbType] = accessor
	return nil
}

// LoadDBConfig loads the database configuration from environment variables
// and sets the default values if not found. It also sets the database type.
func (dbNS) LoadEnv() error {
	var present bool
	var err error

	DB.dbType, present = os.LookupEnv("DB_TYPE")
	if !present {
		logger.Error("Database type is not set in environment variables")
		err = errors.New("Database type is not set")
	}
	logger.Debug("Using database type:", DB.dbType)

	DB.connectionString, present = os.LookupEnv(DB.dbType + "_STRING")
	if !present {
		logger.Error("Database connection string is not set in environment variables")
		err = errors.New("Database connection string is not set")
	}
	logger.Debug("Using database connection string:", DB.connectionString)

	DB.name, present = os.LookupEnv("DB_NAME")
	if !present {
		logger.Error("Database name is not set in environment variables")
		err = errors.New("Database name is not set")

	} else {
		dbEnvName := App.Environment()

		// In case of debug mode, use development database
		if dbEnvName == App.Mode().Debug() {
			dbEnvName = App.Mode().Development()
		}

		DB.name = DB.name + "-" + dbEnvName
		logger.Debug("Using database name:", DB.name)
	}

	return err
}

// ReloadConnection reloads the database connection using the current environment variables
func (dbNS) ReloadConnection() error {
	DB.Close()
	DB.LoadEnv()
	if err := DB.Start(); err != nil {
		logger.Error("Failed to reload DB connection:", err)
		return err
	}
	return nil
}

/*
	Database accessor management
*/

func (dbNS) Use(db InterfaceDBAccessor) *dbNS {
	DB.accessor = db
	return &DB
}

/*
	Database connection management
*/

func (dbNS) Start() error {
	if DB.accessor == nil {
		if len(accessors) == 0 || accessors[DB.dbType] == nil {
			logger.Error("Database accessor is not set")
			return errors.New("Database accessor is not set")
		}

		DB.accessor = accessors[DB.dbType]
	}

	if DB.connectionString == "" || DB.name == "" {
		logger.Error("Database connection string or name is not set")
		return errors.New("Database connection string or name is not set")
	}

	logger.Debug("Connecting to database:", DB.name)
	logger.Debug("Using connection string:", DB.connectionString)
	logger.Debug("Using accessor type:", DB.dbType)
	return DB.accessor.Connect(DB.name, DB.connectionString)
}
func (dbNS) ConnectTo(dbName, connectionString string) error {
	if DB.accessor == nil {
		logger.Error("Database accessor is not set")
		return errors.New("Database accessor is not set")
	}

	DB.connectionString = connectionString
	DB.name = dbName

	return DB.Start()
}
func (dbNS) ConnectOnce(dbName, connectionString string) error {
	if DB.accessor == nil {
		logger.Error("Database accessor is not set")
		return errors.New("Database accessor is not set")
	}

	return DB.accessor.Connect(dbName, connectionString)
}

/*
	Database operations
*/

func (dbNS) InsertOne(document models.DBModelInterface) types.Result[models.DBID] {
	return DB.accessor.InsertOne(document)
}
func (dbNS) InsertMany(documents ...models.DBModelInterface) types.Result[[]models.DBID] {
	return DB.accessor.InsertMany(documents...)
}

func (dbNS) FindOne(filter any, result models.DBModelInterface) types.Result[models.DBModelInterface] {
	return DB.accessor.FindOne(filter, result)
}
func (dbNS) FindAll(filter any, result models.DBModelInterface) types.Result[[]models.DBModelInterface] {
	return DB.accessor.FindMany(filter, result)
}

func (dbNS) UpdateOne(filter any, update models.DBModelInterface) error {
	return DB.accessor.UpdateOne(filter, update)
}
func (dbNS) UpdateMany(filter any, update models.DBModelInterface) error {
	return DB.accessor.UpdateMany(filter, update)
}

func (dbNS) PatchOne(filter any, update models.DBModelInterface) error {
	return DB.accessor.PatchOne(filter, update)
}
func (dbNS) PatchMany(filter any, update models.DBModelInterface) error {
	return DB.accessor.PatchMany(filter, update)
}

func (dbNS) SoftDeleteOne(filter any, model models.DBModelInterface) error {
	return DB.accessor.SoftDeleteOne(filter, model)
}
func (dbNS) SoftDeleteMany(filter any, model models.DBModelInterface) error {
	return DB.accessor.SoftDeleteMany(filter, model)
}

func (dbNS) PermanentDeleteOne(filter any, model models.DBModelInterface) error {
	return DB.accessor.PermanentDeleteOne(filter, model)
}
func (dbNS) PermanentDeleteMany(filter any, model models.DBModelInterface) error {
	return DB.accessor.PermanentDeleteMany(filter, model)
}

// Migrate performs database migrations for the provided models
// It uses the gorm library to automatically migrate the models to the database
func (dbNS) Migrate(models ...models.DBModelInterface) {
	logger.Info("Starting migrations")

	if err := DB.accessor.Migrate(models...); err != nil {
		logger.Error("Migration failed:", err)
		return
	}

	logger.Info("Migrations completed")
}

func (dbNS) Close() error {
	if DB.accessor == nil {
		logger.Error("Database accessor is not set")
		return errors.New("Database accessor is not set")
	}

	return DB.accessor.Disconnect()
}
