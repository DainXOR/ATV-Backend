package configs

import (
	"context"
	"dainxor/atv/logger"
	"dainxor/atv/models"
	"dainxor/atv/types"
	"dainxor/atv/utils"
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

var _ InterfaceDBAccessor = (*mongoType)(nil)

func NewMongoAccessor() *mongoType {
	return &mongoType{}
}

func (m *mongoType) Connect(dbName, connectionString string) error {
	clientOpts := options.Client().ApplyURI(connectionString)
	client, err := mongo.Connect(clientOpts)
	if err != nil {
		logger.Error("Failed to connect to MongoDB:", err)
		logger.Debug("Connection string used:", connectionString)
		logger.Debug("Database name used:", dbName)
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
func (m *mongoType) Disconnect() error {
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
	return bson.D{}
}

func (m mongoType) in(name string) *mongo.Collection {
	return m.db.Collection(name)
}
func (m mongoType) from(v models.DBModelInterface) *mongo.Collection {
	// Use reflection to get the collection name from the struct
	collectionName := v.TableName()
	return m.in(collectionName)
}

func (m mongoType) InsertOne(element models.DBModelInterface) types.Result[models.DBID] {
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
func (m mongoType) InsertMany(elements ...models.DBModelInterface) types.Result[[]models.DBID] {
	if len(elements) == 0 {
		return types.ResultErr[[]models.DBID](DBErr.InvalidInput())
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

func (m mongoType) FindOne(filter any, model models.DBModelInterface) types.Result[models.DBModelInterface] {
	ctx, cancel := m.Context()
	defer cancel()

	concreteModel := utils.InstanceOfUnderlying(model)
	err := m.from(model).FindOne(ctx, filter).Decode(concreteModel)
	if errors.Is(err, mongo.ErrNoDocuments) {
		logger.Debug("No document found for filter:", filter)
		return types.ResultErr[models.DBModelInterface](DBErr.NotFound())
	}

	model = utils.AsDeref(concreteModel).(models.DBModelInterface)
	return types.ResultOf(model, err, err != nil)
}
func (m mongoType) FindMany(filter any, model models.DBModelInterface) types.Result[[]models.DBModelInterface] {
	ctx, cancel := m.Context()
	defer cancel()

	cursor, err := m.from(model).Find(ctx, filter)
	if err != nil {
		logger.Error("Failed to find documents:", err)
		return types.ResultErr[[]models.DBModelInterface](err)
	}
	defer cursor.Close(ctx)

	concreteSlice := utils.SliceOfUnderlying(model)
	cursorErr := cursor.All(ctx, &concreteSlice)
	temp := utils.AsSliceOfNoPtr[models.DBModelInterface](concreteSlice)

	return types.ResultOf(temp, cursorErr, cursorErr != nil)
}

func (m mongoType) UpdateOne(filter any, model models.DBModelInterface) error {
	ctx, cancel := m.Context()
	defer cancel()

	update := bson.D{{Key: "$set", Value: model}}
	updateResult, err := m.from(model).UpdateOne(ctx, filter, update)
	if err != nil {
		logger.Error("Failed to update document:", err)
		return err
	}
	if updateResult.MatchedCount == 0 {
		logger.Warning("No documents matched the filter for update:", filter)
		return DBErr.NotFound()
	}
	if updateResult.ModifiedCount == 0 {
		logger.Warning("No documents were modified by the update:", update)
		return DBErr.NotModified()
	}

	return nil
}
func (m mongoType) UpdateMany(filter any, model models.DBModelInterface) error {
	ctx, cancel := m.Context()
	defer cancel()

	update := bson.D{{Key: "$set", Value: model}}
	res, err := m.from(model).UpdateMany(ctx, filter, update)
	if err != nil {
		logger.Warning("Failed to update documents:", err)
		return err
	}
	if res.MatchedCount == 0 {
		logger.Warning("No documents matched the filter for update:", filter)
		return DBErr.NotFound()
	}
	if res.ModifiedCount == 0 {
		logger.Warning("No documents were modified by the update:", update)
		return DBErr.NotModified()
	}

	return nil
}

func (m mongoType) PatchOne(filter any, model models.DBModelInterface) error {
	ctx, cancel := m.Context()
	defer cancel()

	update := bson.D{{Key: "$set", Value: model}}
	updateResult, err := m.from(model).UpdateOne(ctx, filter, update)
	if err != nil {
		logger.Error("Failed to patch document:", err)
		return err
	}
	if updateResult.MatchedCount == 0 {
		logger.Warning("No documents matched the filter for update:", filter)
		return DBErr.NotFound()
	}
	if updateResult.ModifiedCount == 0 {
		logger.Warning("No documents were modified by the update:", update)
		return DBErr.NotModified()
	}

	return nil
}
func (m mongoType) PatchMany(filter any, model models.DBModelInterface) error {
	ctx, cancel := m.Context()
	defer cancel()

	update := bson.D{{Key: "$set", Value: model}}
	updateResult, err := m.from(model).UpdateMany(ctx, filter, update)
	if err != nil {
		logger.Error("Failed to update documents:", err)
		return err
	}
	if updateResult.MatchedCount == 0 {
		logger.Warning("No documents matched the filter for update:", filter)
		return DBErr.NotFound()
	}
	if updateResult.ModifiedCount == 0 {
		logger.Warning("No documents were modified by the update:", update)
		return DBErr.NotModified()
	}

	return nil
}

func (m mongoType) SoftDeleteOne(filter any, model models.DBModelInterface) error {
	ctx, cancel := m.Context()
	defer cancel()

	update := bson.M{"$set": bson.M{"deleted_at": time.Now()}}
	updateResult, err := m.from(model).UpdateOne(ctx, filter, update)
	if err != nil {
		logger.Error("Failed to soft delete document:", err)
		return err
	}
	if updateResult.MatchedCount == 0 {
		logger.Warning("No documents matched the filter for soft delete:", filter)
		return DBErr.NotFound()
	}
	if updateResult.ModifiedCount == 0 {
		logger.Warning("No documents were modified by the soft delete:", update)
		return DBErr.NotModified()
	}

	return nil
}
func (m mongoType) SoftDeleteMany(filter any, model models.DBModelInterface) error {
	ctx, cancel := m.Context()
	defer cancel()

	update := bson.M{"$set": bson.M{"deleted_at": time.Now()}}
	updateResult, err := m.from(model).UpdateMany(ctx, filter, update)
	if err != nil {
		logger.Warning("Failed to update documents:", err)
		return err
	}
	if updateResult.MatchedCount == 0 {
		logger.Warning("No documents matched the filter for update:", filter)
		return DBErr.NotFound()
	}
	if updateResult.ModifiedCount == 0 {
		logger.Warning("No documents were modified by the update:", model)
		return DBErr.NotModified()
	}

	return nil
}

func (m mongoType) PermanentDeleteOne(filter any, model models.DBModelInterface) error {
	ctx, cancel := m.Context()
	defer cancel()

	deleteResult, err := m.from(model).DeleteOne(ctx, filter)
	if err != nil {
		logger.Error("Failed to delete document:", err)
		return err
	}
	if deleteResult.DeletedCount == 0 {
		logger.Warning("No documents matched the filter for delete:", filter)
		return DBErr.NotFound()
	}

	return nil
}
func (m mongoType) PermanentDeleteMany(filter any, model models.DBModelInterface) error {
	ctx, cancel := m.Context()
	defer cancel()

	deleteResult, err := m.from(model).DeleteMany(ctx, filter)
	if err != nil {
		logger.Error("Failed to delete documents:", err)
		return err
	}
	if deleteResult.DeletedCount == 0 {
		logger.Warning("No documents matched the filter for delete:", filter)
		return DBErr.NotFound()
	}

	return nil
}
