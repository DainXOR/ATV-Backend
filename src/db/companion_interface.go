package db

import (
	"dainxor/atv/configs"
	"dainxor/atv/logger"
	"dainxor/atv/models"
	"dainxor/atv/types"
	"dainxor/atv/utils"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type companionType struct{}

var Companion companionType

func (companionType) Create(companion models.CompanionCreate) types.Result[models.CompanionDBMongo] {
	companionDB := companion.ToInsert()
	result, err := configs.DB.InsertOne(companionDB)

	if err != nil {
		logger.Error("Failed to create companion in MongoDB: ", err)
		return types.ResultErr[models.CompanionDBMongo](err)
	}

	companionDB.ID, err = models.DBIDFrom(result.InsertedID)

	if err != nil {
		logger.Error("Failed to convert inserted ID to ObjectID: ", err)
		httpErr := types.ErrorInternal(
			"Failed to create companion",
			"Failed to convert inserted ID to ObjectID",
			"Error: "+err.Error(),
		)

		return types.ResultErr[models.CompanionDBMongo](&httpErr)
	}

	return types.ResultOk(companionDB)
}

func (companionType) GetByID(id string) types.Result[models.CompanionDBMongo] {
	oid, err := bson.ObjectIDFromHex(id)

	if err != nil {
		logger.Error("Failed to convert ID to ObjectID: ", err)
		httpErr := types.Error(
			types.Http.UnprocessableEntity(),
			"Invalid value",
			"Invalid ID format: "+err.Error(),
			"Companion ID: "+id,
		)
		return types.ResultErr[models.CompanionDBMongo](&httpErr)
	}

	filter := bson.D{{Key: "_id", Value: oid}}
	var companionF models.CompanionDBMongoReceiver

	err = configs.DB.FindOne(filter, &companionF)
	if err != nil {
		var httpErr types.HttpError

		if err == mongo.ErrNoDocuments {
			logger.Error("Failed to get companion by ID: ", err)
			httpErr = types.ErrorNotFound(
				"Companion not found",
				"Companion with ID "+id+" not found",
			)
		} else {
			logger.Error("Failed to get companion by ID: ", err)
			httpErr = types.ErrorInternal(
				"Failed to retrieve companion",
				"Decoding error",
				err.Error(),
				"Companion ID: "+id,
			)
		}

		return types.ResultErr[models.CompanionDBMongo](&httpErr)
	}

	return types.ResultOk(companionF.ToDB())
}
func (companionType) GetByNumberID(idNumber string) types.Result[models.CompanionDBMongo] {
	filter := bson.D{{Key: "id_number", Value: idNumber}}
	var companion models.CompanionDBMongoReceiver

	err := configs.DB.FindOne(filter, &companion)
	if err != nil {
		var httpErr types.HttpError
		if err == mongo.ErrNoDocuments {
			logger.Error("Failed to get companion by ID number: ", err)
			httpErr = types.ErrorNotFound(
				"Companion not found",
				"Companion with ID number "+idNumber+" not found",
			)
		} else {
			logger.Error("Failed to get companion by ID number: ", err)
			httpErr = types.ErrorInternal(
				"Failed to retrieve companion",
				"Decoding error",
				err.Error(),
				"Companion ID number: "+idNumber,
			)
		}

		return types.ResultErr[models.CompanionDBMongo](&httpErr)
	}

	return types.ResultOk(companion.ToDB())
}
func (companionType) GetByEmail(email string) types.Result[models.CompanionDBMongo] {
	filter := bson.D{{Key: "email", Value: email}}
	var companionF models.CompanionDBMongoReceiver

	err := configs.DB.FindOne(filter, &companionF)
	if err != nil {
		var httpErr types.HttpError

		if err == mongo.ErrNoDocuments {
			logger.Error("Failed to get companion by email: ", err)
			httpErr = types.ErrorNotFound(
				"Companion not found",
				"Companion with email "+email+" not found",
			)
		} else {
			logger.Error("Failed to get companion by email: ", err)
			httpErr = types.ErrorInternal(
				"Failed to retrieve companion",
				"Decoding error",
				err.Error(),
				"Companion email: "+email,
			)
		}

		return types.ResultErr[models.CompanionDBMongo](&httpErr)
	}

	return types.ResultOk(companionF.ToDB())
}
func (companionType) GetAll() types.Result[[]models.CompanionDBMongo] {
	filter := bson.D{{Key: "deleted_at", Value: models.TimeZero()}} // Filter to exclude deleted companions
	companionsR := models.CompanionDBMongo{}.ReceiverList()

	err := configs.DB.FindAll(filter, &companionsR)
	if err != nil {
		logger.Error("Failed to get all companions from MongoDB:", err)
		httpErr := types.ErrorInternal(
			"Failed to retrieve companions",
			err.Error(),
		)

		return types.ResultErr[[]models.CompanionDBMongo](&httpErr)
	}

	companions := utils.Map(companionsR, models.CompanionDBMongoReceiver.ToDB)
	logger.Debug("Retrieved", len(companions), "companions from MongoDB database")
	return types.ResultOk(companions)
}

func (companionType) UpdateByID(id string, companion models.CompanionCreate) types.Result[models.CompanionDBMongo] {
	oid, err := models.BsonIDFrom(id)
	if err != nil {
		logger.Error("Failed to convert ID to ObjectID: ", err)
		httpErr := types.Error(
			types.Http.UnprocessableEntity(),
			"Invalid value",
			"Invalid ID format: "+err.Error(),
			"Companion ID: "+id,
		)
		return types.ResultErr[models.CompanionDBMongo](&httpErr)
	}

	filter := bson.D{{Key: "_id", Value: oid}}
	update := bson.D{{Key: "$set", Value: companion.ToUpdate()}}
	companionDB := companion.ToUpdate().Receiver()

	result := configs.DB.PatchOne(filter, update, companionDB)
	// .From(models.CompanionDBMongo{}).UpdateOne(ctx, filter, update)

	if result.IsErr() {
		err := result.Error()
		logger.Error("Failed to update companion in MongoDB: ", err)
		httpErr := types.ErrorInternal(
			"Failed to update companion",
			err.Error(),
			"Companion ID: "+id,
		)
		return types.ResultErr[models.CompanionDBMongo](&httpErr)
	}

	if result.Value().MatchedCount == 0 {
		httpErr := types.ErrorNotFound(
			"Companion not found",
			"Companion with ID "+id+" not found",
		)
		return types.ResultErr[models.CompanionDBMongo](&httpErr)
	}

	if result.Value().ModifiedCount == 0 {
		logger.Info("No changes made to companion with ID: ", id)
		logger.Lava(2, "Send a more proper code for no changes made")
		httpErr := types.Error(
			types.Http.C200().Accepted(),
			"No changes made",
			"Companion with ID "+id+" was not modified",
		)
		return types.ResultErr[models.CompanionDBMongo](&httpErr)
	}

	return Companion.GetByID(id)
}

func (companionType) PatchByID(id string, companion models.CompanionCreate) types.Result[models.CompanionDBMongo] {
	oid, err := models.BsonIDFrom(id)
	if err != nil {
		logger.Error("Failed to convert ID to ObjectID: ", err)
		httpErr := types.Error(
			types.Http.UnprocessableEntity(),
			"Invalid value",
			"Invalid ID format: "+err.Error(),
			"Companion ID: "+id,
		)
		return types.ResultErr[models.CompanionDBMongo](&httpErr)
	}

	companionDB := companion.ToUpdate()
	if companionDB == (models.CompanionDBMongo{}) {
		logger.Error("Error converting companion to DB model")
		httpErr := types.Error(
			types.Http.C400().UnprocessableEntity(),
			"Invalid value",
			"Invalid companion data",
			"Companion ID: "+id,
		)
		return types.ResultErr[models.CompanionDBMongo](&httpErr)
	}

	filter := bson.D{{Key: "_id", Value: oid}}
	update := bson.D{{Key: "$set", Value: companionDB}}
	receiver := companionDB.Receiver()

	result := configs.DB.PatchOne(filter, update, &receiver)

	if result.IsErr() {
		err := result.Error()
		logger.Error("Failed to update companion in MongoDB: ", err)
		httpErr := types.ErrorInternal(
			"Failed to update companion",
			err.Error(),
			"Companion ID: "+id,
		)
		return types.ResultErr[models.CompanionDBMongo](&httpErr)
	}

	if result.Value().MatchedCount == 0 {
		httpErr := types.ErrorNotFound(
			"Companion not found",
			"Companion with ID "+id+" not found",
		)
		return types.ResultErr[models.CompanionDBMongo](&httpErr)
	}

	if result.Value().ModifiedCount == 0 {
		logger.Info("No changes made to companion with ID: ", id)
		logger.Lava(2, "Send a more proper code for no changes made")
		httpErr := types.Error(
			types.Http.C200().Accepted(),
			"No changes made",
			"Companion with ID "+id+" was not modified",
		)
		return types.ResultErr[models.CompanionDBMongo](&httpErr)
	}

	return types.ResultOk(receiver.ToDB())
}

func (companionType) DeleteByID(id string) types.Result[models.CompanionDBMongo] {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		logger.Error("Failed to convert ID to ObjectID: ", err)
		httpErr := types.Error(
			types.Http.UnprocessableEntity(),
			"Invalid value",
			"Invalid ID format: "+err.Error(),
			"Companion ID: "+id,
		)
		return types.ResultErr[models.CompanionDBMongo](&httpErr)
	}

	filter := bson.D{{Key: "_id", Value: oid}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "deleted_at", Value: time.Now()}}}}
	ctx, cancel := configs.DB.Context()
	defer cancel()

	var deletedCompanion models.CompanionDBMongoReceiver
	//result, err := configs.DB.UpdateOne(filter, update, deletedCompanion)
	v := logger.Lava(1, "Use the code above to update the companion with deleted_at field")
	v.LavaStart()
	result, err := configs.DB.From(models.CompanionDBMongo{}).UpdateOne(ctx, filter, update)
	if err != nil {
		logger.Error("Failed to delete companion in MongoDB: ", err)
		httpErr := types.ErrorInternal(
			"Failed to delete companion",
			err.Error(),
			"Companion ID: "+id,
		)
		return types.ResultErr[models.CompanionDBMongo](&httpErr)
	}

	if result.MatchedCount == 0 {
		httpErr := types.ErrorNotFound(
			"Companion not found",
			"Companion with ID "+id+" not found",
		)
		return types.ResultErr[models.CompanionDBMongo](&httpErr)
	}

	err = configs.DB.FindOne(filter, &deletedCompanion)
	if err != nil {
		logger.Error("Failed to retrieve deleted companion: ", err)
		httpErr := types.ErrorInternal(
			"Failed to retrieve deleted companion",
			err.Error(),
			"Companion ID: "+id,
		)
		return types.ResultErr[models.CompanionDBMongo](&httpErr)
	}

	v.LavaEnd()

	return types.ResultOk(deletedCompanion.ToDB())
}
func (companionType) DeletePermanentByID(id string) types.Result[models.CompanionDBMongo] {
	logger.Warning("Permanently deleting companion by ID: ", id)
	oid, err := models.BsonIDFrom(id)
	if err != nil {
		logger.Error("Failed to convert ID to ObjectID: ", err)
		httpErr := types.Error(
			types.Http.UnprocessableEntity(),
			"Invalid value",
			"Invalid ID format: "+err.Error(),
			"Companion ID: "+id,
		)
		return types.ResultErr[models.CompanionDBMongo](&httpErr)
	}

	filter := bson.D{{Key: "_id", Value: oid}, {Key: "deleted_at", Value: bson.M{"$ne": time.Time{}}}} // Ensure the companion is marked as deleted
	ctx, cancel := configs.DB.Context()
	defer cancel()

	var companion models.CompanionDBMongoReceiver
	err = configs.DB.FindOne(filter, &companion)
	if err != nil {
		logger.Debug("Failed to find companion for permanent deletion: ", err)

		if err == mongo.ErrNoDocuments {
			httpErr := types.ErrorNotFound(
				"Companion not found",
				"Companion with ID "+id+" not found or not marked as deleted",
			)
			return types.ResultErr[models.CompanionDBMongo](&httpErr)
		}

		httpErr := types.ErrorInternal(
			"Failed to find companion for permanent deletion",
			err.Error(),
			"Companion ID: "+id,
		)
		return types.ResultErr[models.CompanionDBMongo](&httpErr)
	}

	result, err := configs.DB.From(models.CompanionDBMongo{}).DeleteOne(ctx, filter)
	if err != nil {
		logger.Debug("Failed to permanently delete companion in MongoDB: ", err)
		httpErr := types.ErrorInternal(
			"Failed to permanently delete companion",
			err.Error(),
			"Companion ID: "+id,
		)
		return types.ResultErr[models.CompanionDBMongo](&httpErr)
	}

	if result.DeletedCount == 0 {
		httpErr := types.ErrorNotFound(
			"Companion not found.",
			"Companion with ID "+id+" not found.",
			"Ensure the companion is marked as deleted before permanent deletion.",
		)
		return types.ResultErr[models.CompanionDBMongo](&httpErr)
	}

	return types.ResultOk(companion.ToDB())
}
func (companionType) DeletePermanentAll() types.Result[[]models.CompanionDBMongo] {
	filter := bson.D{{Key: "deleted_at", Value: bson.M{"$ne": nil}}}
	ctx, cancel := configs.DB.Context()
	defer cancel()

	result, err := configs.DB.From(models.CompanionDBMongo{}).DeleteMany(ctx, filter)
	if err != nil {
		logger.Error("Failed to permanently delete all companions in MongoDB: ", err)
		httpErr := types.ErrorInternal(
			"Failed to permanently delete all companions",
			err.Error(),
		)
		return types.ResultErr[[]models.CompanionDBMongo](&httpErr)
	}

	if result.DeletedCount == 0 {
		httpErr := types.ErrorNotFound(
			"No deleted companions found",
			"No companions marked as deleted found",
		)
		return types.ResultErr[[]models.CompanionDBMongo](&httpErr)
	}

	return types.ResultOk([]models.CompanionDBMongo{})
}
