package configs

import (
	"context"
	"dainxor/atv/logger"
	"dainxor/atv/models"
	"dainxor/atv/types"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type mongoType struct {
	client *mongo.Client
	db     *mongo.Database
}

var mongoT mongoType

var _ InterfaceDB = (*mongoType)(nil)

func (m mongoType) Connect(dbName, conectionString string) error {
	clientOpts := options.Client().ApplyURI(conectionString)
	client, err := mongo.Connect(clientOpts)
	if err != nil {
		logger.Error("Failed to connect to MongoDB:", err)
		return err
	}

	ctx, cancel := m.Context()
	defer cancel()
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		logger.Error("Failed to ping MongoDB:", err)
		return err
	}

	m.client = client
	m.db = client.Database(dbName)
	return nil
}
func (m mongoType) Disconnect() error {
	if m.client != nil {
		err := m.client.Disconnect(context.Background())
		if err != nil {
			return err
		}
		m.client = nil
		m.db = nil
	}
	return nil
}
func (m mongoType) Migrate(models ...models.DBModelInterface) error {
	return nil // No migrations needed for MongoDB
}
func (m mongoType) Context() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}

func (m mongoType) CreateFilter(filter ...types.SPair[string]) any {
	return bson.E{}
}

func (m mongoType) in(name string) *mongo.Collection {
	return m.db.Collection(name)
}
func (m mongoType) from(v models.DBModelInterface) *mongo.Collection {
	// Use reflection to get the collection name from the struct
	collectionName := v.TableName()
	return m.in(collectionName)
}

func (m mongoType) CreateOne(element models.DBModelInterface) types.Result[models.DBID] {
	ctx, cancel := m.Context()
	defer cancel()

	res, err := m.from(element).InsertOne(ctx, element)
	if err != nil {
		logger.Warning("Failed to insert document:", err)
		return types.ResultErr[models.DBID](err)
	}

	id, _ := models.ID.ToDB(res.InsertedID)
	return types.ResultOk(id)
}
func (m mongoType) CreateMany(elements ...models.DBModelInterface) types.Result[[]models.DBID] {
	if len(elements) == 0 {
		return types.ResultErr[[]models.DBID](dbErr.InvalidInput())
	}

	ctx, cancel := m.Context()
	defer cancel()

	var results []models.DBID
	for i, element := range elements {
		res, err := m.from(element).InsertOne(ctx, element)
		if err != nil {
			logger.Warning("Failed to insert document:", err)
			logger.Info("Inserted documents before failure:", i)

			return types.ResultOf(results, err, true)
		}
		id, _ := models.ID.ToDB(res.InsertedID)
		results = append(results, id)
	}
	return types.ResultOk(results)
}

func (m mongoType) GetOne(filter any, model models.DBModelInterface) types.Result[models.DBModelInterface] {
	ctx, cancel := m.Context()
	defer cancel()

	err := m.from(model).FindOne(ctx, filter).Decode(&model)
	if errors.Is(err, mongo.ErrNoDocuments) {
		logger.Debug("No document found for filter:", filter)
		return types.ResultErr[models.DBModelInterface](dbErr.NotFound())
	}

	return types.ResultOf(model, err, err != nil)
}
func (m mongoType) GetMany(filter any, model ...models.DBModelInterface) types.Result[[]models.DBModelInterface] {
	if len(model) == 0 {
		logger.Debug("No models provided for GetMany")
		return types.ResultErr[[]models.DBModelInterface](dbErr.InvalidInput())
	}

	ctx, cancel := m.Context()
	defer cancel()

	cursor, err := m.from(model[0]).Find(ctx, filter)
	if err != nil {
		logger.Error("Failed to find documents:", err)
		return types.ResultErr[[]models.DBModelInterface](err)
	}
	defer cursor.Close(ctx)

	cursorErr := cursor.All(ctx, &model)
	return types.ResultOf(model, cursorErr, cursorErr != nil)
}

func (m mongoType) UpdateOne(filter any, update models.DBModelInterface) error {
	ctx, cancel := m.Context()
	defer cancel()

	updateResult, err := m.from(update).UpdateOne(ctx, filter, update)
	if err != nil {
		logger.Error("Failed to update document:", err)
		return err
	}
	if updateResult.MatchedCount == 0 {
		logger.Warning("No documents matched the filter for update:", filter)
		return dbErr.NotFound()
	}
	if updateResult.ModifiedCount == 0 {
		logger.Warning("No documents were modified by the update:", update)
		return dbErr.NotModified()
	}

	return nil
}
func (m mongoType) UpdateMany(filter any, update models.DBModelInterface) error {
	ctx, cancel := m.Context()
	defer cancel()

	res, err := m.from(update).UpdateMany(ctx, filter, update)
	if err != nil {
		logger.Warning("Failed to update documents:", err)
		return err
	}
	if res.MatchedCount == 0 {
		logger.Warning("No documents matched the filter for update:", filter)
		return dbErr.NotFound()
	}
	if res.ModifiedCount == 0 {
		logger.Warning("No documents were modified by the update:", update)
		return dbErr.NotModified()
	}

	return nil
}

func (m mongoType) PatchOne(filter any, update models.DBModelInterface) error {
	ctx, cancel := m.Context()
	defer cancel()

	updateResult, err := m.from(update).UpdateOne(ctx, filter, update)
	if err != nil {
		logger.Error("Failed to patch document:", err)
		return err
	}
	if updateResult.MatchedCount == 0 {
		logger.Warning("No documents matched the filter for update:", filter)
		return dbErr.NotFound()
	}
	if updateResult.ModifiedCount == 0 {
		logger.Warning("No documents were modified by the update:", update)
		return dbErr.NotModified()
	}

	return nil
}
func (m mongoType) PatchMany(filter any, update models.DBModelInterface) error {
	ctx, cancel := m.Context()
	defer cancel()

	updateResult, err := m.from(update).UpdateMany(ctx, filter, update)
	if err != nil {
		logger.Error("Failed to update documents:", err)
		return err
	}
	if updateResult.MatchedCount == 0 {
		logger.Warning("No documents matched the filter for update:", filter)
		return dbErr.NotFound()
	}
	if updateResult.ModifiedCount == 0 {
		logger.Warning("No documents were modified by the update:", update)
		return dbErr.NotModified()
	}

	return nil
}

func (m mongoType) DeleteOne(filter any, model models.DBModelInterface) error {
	ctx, cancel := m.Context()
	defer cancel()

	deleteResult, err := m.from(model).DeleteOne(ctx, filter)
	if err != nil {
		logger.Error("Failed to delete document:", err)
		return err
	}
	if deleteResult.DeletedCount == 0 {
		logger.Warning("No documents matched the filter for delete:", filter)
		return dbErr.NotFound()
	}

	return nil
}
func (m mongoType) DeleteMany(filter any, model models.DBModelInterface) error {
	ctx, cancel := m.Context()
	defer cancel()

	deleteResult, err := m.from(model).DeleteMany(ctx, filter)
	if err != nil {
		logger.Error("Failed to delete documents:", err)
		return err
	}
	if deleteResult.DeletedCount == 0 {
		logger.Warning("No documents matched the filter for delete:", filter)
		return dbErr.NotFound()
	}

	return nil
}
