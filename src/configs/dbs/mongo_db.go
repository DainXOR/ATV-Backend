package db

import (
	"context"
	"dainxor/atv/logger"
	"dainxor/atv/models"
	"dainxor/atv/types"
	"dainxor/atv/utils"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type mongoType struct {
	client        *mongo.Client
	db            *mongo.Database
	disconectFunc func()
}

var mongoT mongoType

//var _ InterfaceDB = (*mongoType)(nil)

func (mongoType) Disconnect() {
	if mongoT.disconectFunc != nil {
		mongoT.disconectFunc()
	} else {
		logger.Warning("MongoDB disconect function is nil, nothing to do")
	}
}
func (mongoType) Context() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}

func (mongoType) CreateOne(document models.DBModelInterface) (*mongo.InsertOneResult, error) {
	ctx, cancel := mongoT.Context()
	defer cancel()

	return mongoT.db.Collection(document.TableName()).InsertOne(ctx, document)
}
func (mongoType) CreateMany(documents any) types.Result[any] {
	return types.ResultErr[any](errors.ErrUnsupported)
}

func (mongoType) GetOne(filter any, result models.DBModelInterface) types.Result[models.DBModelInterface] {
	ctx, cancel := mongoT.Context()
	defer cancel()

	err := mongoT.db.Collection(result.TableName()).FindOne(ctx, filter).Decode(result)
	if errors.Is(err, mongo.ErrNoDocuments) {
		logger.Warning("No document found for filter:", filter)
		return types.ResultErr[models.DBModelInterface](ErrNotFound)
	}

	return types.ResultOf(result, err, err != nil)
}
func (mongoType) GetAll(filter any, result any) types.Result[any] {
	logger.Lava("0.1.1", "This mf should be refactored to use []models.DBModelInterface instead of any for the result")

	eType, err := utils.SliceType(result)
	if err != nil {
		logger.Error("This function ONLY works with slices or pointers to slices")
		return types.ResultErr[any](ErrInvalidInput)
	}

	iType, ok := eType.(interface{ TableName() string })
	if !ok {
		logger.Error("Result type does NOT IMPLEMENT TableName method")
		return types.ResultErr[any](ErrInvalidInput)
	}

	ctx, cancel := mongoT.Context()
	defer cancel()

	cursor, err := mongoT.db.Collection(iType.TableName()).Find(ctx, filter)
	if err != nil {
		logger.Error("Failed to find documents:", err)
		return types.ResultErr[any](err)
	}
	defer cursor.Close(ctx)

	cursorErr := cursor.All(ctx, result)
	return types.ResultOf(result, cursorErr, cursorErr != nil)
}

func (mongoType) UpdateOne(filter any, update any, result models.DBModelInterface) types.Result[models.DBModelInterface] {
	ctx, cancel := mongoT.Context()
	defer cancel()

	updateResult, err := mongoT.db.Collection(result.TableName()).UpdateOne(ctx, filter, update)
	if err != nil {
		logger.Error("Failed to update document:", err)
		return types.ResultErr[models.DBModelInterface](err)
	}
	if updateResult.MatchedCount == 0 {
		logger.Warning("No documents matched the filter for update:", filter)
		return types.ResultErr[models.DBModelInterface](ErrNotFound)
	}
	if updateResult.ModifiedCount == 0 {
		logger.Warning("No documents were modified by the update:", update)
		return types.ResultErr[models.DBModelInterface](ErrNotModified)
	}

	res := mongoT.GetOne(filter, result)
	if res.IsErr() {
		logger.Error("Failed to find updated document:", res.Error())
		return res
	}
	return res
}

func (mongoType) PatchOne(filter any, update any, result models.DBModelInterface) types.Result[mongo.UpdateResult] {
	ctx, cancel := mongoT.Context()
	defer cancel()

	updateResult, err := mongoT.db.Collection(result.TableName()).UpdateOne(ctx, filter, update)
	if err != nil {
		logger.Error("Failed to update document:", err)
		return types.ResultErr[mongo.UpdateResult](err)
	}

	res := mongoT.GetOne(filter, result)

	return types.ResultOf(*updateResult, res.Error(), res.IsErr())
}

// ConnectMongoDB connects to the MongoDB database
func (mongoType) ConnectMongoDB(dbName, conectionString string) {
	logger.Debug("Connecting to MongoDB: ", conectionString)
	clientOpts := options.Client().ApplyURI(conectionString)
	client, err := mongo.Connect(clientOpts)
	if err != nil {
		logger.Fatal(err)
	}

	ctx, cancel := mongoT.Context()
	defer cancel()
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		logger.Fatal(err)
	}

	mongoT.client = client
	mongoT.db = client.Database(dbName)

	mongoT.disconectFunc = func() {
		if err = mongoT.client.Disconnect(context.Background()); err != nil {
			logger.Fatal("Error disconnecting from MongoDB: ", err)
		}
	}
}
