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

func (companionType) Create(companion models.CompanionCreate) types.Result[models.CompanionDB] {
	companionDB := companion.ToInsert()
	resultCreate := configs.DB.InsertOne(companionDB)
	if resultCreate.IsErr() {
		logger.Warning("Failed to create companion in MongoDB:", resultCreate.Error())
		return types.ResultErr[models.CompanionDB](resultCreate.Error())
	}

	resultGet := configs.DB.FindOne(bson.D{{Key: "_id", Value: resultCreate.Value()}}, &companionDB)

	if resultGet.IsErr() {
		logger.Warning("Failed to get created companion in MongoDB:", resultGet.Error())
		return types.ResultErr[models.CompanionDB](resultGet.Error())
	}

	var err error
	companionDB.ID, err = models.ID.ToDB(resultCreate.Value())

	if err != nil {
		logger.Warning("Failed to convert inserted ID to ObjectID:", err)
		httpErr := types.ErrorInternal(
			"Failed to create companion",
			"Failed to convert inserted ID to ObjectID",
			"Error: "+err.Error(),
		)

		return types.ResultErr[models.CompanionDB](&httpErr)
	}

	return types.ResultOk(companionDB)
}

func (companionType) GetByID(id string) types.Result[models.CompanionDB] {
	oid, err := models.ID.ToBson(id)

	if err != nil {
		logger.Warning("Failed to convert ID to ObjectID: ", err)
		httpErr := types.Error(
			types.Http.C400().UnprocessableEntity(),
			"Invalid value",
			"Invalid ID format: "+err.Error(),
			"Companion ID: "+id,
		)
		return types.ResultErr[models.CompanionDB](&httpErr)
	}

	filter := bson.D{{Key: "_id", Value: oid}}
	var companion models.CompanionDB

	resultGet := configs.DB.FindOne(filter, companion)
	if resultGet.IsErr() {
		logger.Warning("Failed to get companion by ID: ", resultGet.Error())

		return types.ResultErr[models.CompanionDB](resultGet.Error())
	}
	companion = resultGet.Value().(models.CompanionDB)
	return types.ResultOk(companion)
}
func (companionType) GetByNumberID(idNumber string) types.Result[models.CompanionDB] {
	filter := bson.D{{Key: "number_id", Value: idNumber}}
	var companion models.CompanionDB

	resultGet := configs.DB.FindOne(filter, &companion)
	if resultGet.IsErr() {
		logger.Warning("Failed to get companion by number id:", resultGet.Error())

		return types.ResultErr[models.CompanionDB](resultGet.Error())
	}

	companion = resultGet.Value().(models.CompanionDB)
	return types.ResultOk(companion)
}
func (companionType) GetByEmail(email string) types.Result[models.CompanionDB] {
	filter := bson.D{{Key: "email", Value: email}}
	var companion models.CompanionDB

	resultGet := configs.DB.FindOne(filter, &companion)
	if resultGet.IsErr() {
		logger.Warning("Failed to get companion by email:", resultGet.Error())

		return types.ResultErr[models.CompanionDB](resultGet.Error())

	}

	companion = resultGet.Value().(models.CompanionDB)
	return types.ResultOk(companion)
}
func (companionType) GetAll() types.Result[[]models.CompanionDB] {
	filter := bson.D{models.Filter.NotDeleted()} // Filter to exclude deleted companions
	companions := []models.CompanionDB{{}}

	resultGet := configs.DB.FindAll(filter, companions[0])
	if resultGet.IsErr() {
		logger.Error("Failed to get all companions from MongoDB:", resultGet.Error())

		return types.ResultErr[[]models.CompanionDB](resultGet.Error())
	}

	castMapper := func(m models.DBModelInterface) models.CompanionDB { return m.(models.CompanionDB) }
	companions = utils.Map(resultGet.Value(), castMapper)
	logger.Debug("Retrieved", len(companions), "companions from MongoDB database")
	return types.ResultOk(companions)
}

func (companionType) UpdateByID(id string, companion models.CompanionCreate) types.Result[models.CompanionDB] {
	oid, err := models.ID.ToDB(id)
	if err != nil {
		logger.Error("Failed to convert ID to ObjectID: ", err)
		httpErr := types.Error(
			types.Http.C400().UnprocessableEntity(),
			"Invalid value",
			"Invalid ID format: "+err.Error(),
			"Companion ID: "+id,
		)
		return types.ResultErr[models.CompanionDB](&httpErr)
	}

	companionDB := companion.ToUpdate()
	if companionDB == (models.CompanionDB{}) {
		logger.Error("Error converting companion to DB model")
		httpErr := types.Error(
			types.Http.C400().UnprocessableEntity(),
			"Invalid value",
			"Invalid companion data",
			"Companion ID: "+id,
		)
		return types.ResultErr[models.CompanionDB](&httpErr)
	}

	filter := bson.D{{Key: "_id", Value: oid}}
	err = configs.DB.UpdateOne(filter, companionDB)

	if err != nil {
		logger.Error("Failed to update companion in MongoDB: ", err)
		httpErr := types.ErrorInternal(
			"Failed to update companion",
			err.Error(),
			"Companion ID: "+id,
		)
		return types.ResultErr[models.CompanionDB](&httpErr)
	}

	return Companion.GetByID(id)
}

func (companionType) PatchByID(id string, companion models.CompanionCreate) types.Result[models.CompanionDB] {
	oid, err := models.ID.ToDB(id)
	if err != nil {
		logger.Error("Failed to convert ID to ObjectID: ", err)
		httpErr := types.Error(
			types.Http.C400().UnprocessableEntity(),
			"Invalid value",
			"Invalid ID format: "+err.Error(),
			"Companion ID: "+id,
		)
		return types.ResultErr[models.CompanionDB](&httpErr)
	}

	companionDB := companion.ToUpdate()
	if companionDB == (models.CompanionDB{}) {
		logger.Error("Error converting companion to DB model")
		httpErr := types.Error(
			types.Http.C400().UnprocessableEntity(),
			"Invalid value",
			"Invalid companion data",
			"Companion ID: "+id,
		)
		return types.ResultErr[models.CompanionDB](&httpErr)
	}

	filter := bson.D{{Key: "_id", Value: oid}}
	err = configs.DB.PatchOne(filter, companionDB)

	if err != nil {
		logger.Error("Failed to update companion in MongoDB: ", err)
		httpErr := types.ErrorInternal(
			"Failed to update companion",
			err.Error(),
			"Companion ID: "+id,
		)
		return types.ResultErr[models.CompanionDB](&httpErr)
	}

	return Companion.GetByID(id)
}

func (companionType) DeleteByID(id string) types.Result[models.CompanionDB] {
	oid, err := models.ID.ToBson(id)
	if err != nil {
		logger.Error("Failed to convert ID to ObjectID: ", err)
		httpErr := types.Error(
			types.Http.C400().UnprocessableEntity(),
			"Invalid value",
			"Invalid ID format: "+err.Error(),
			"Companion ID: "+id,
		)
		return types.ResultErr[models.CompanionDB](&httpErr)
	}

	filter := bson.D{{Key: "_id", Value: oid}}
	var companion models.CompanionDB
	companionResult := configs.DB.FindOne(filter, companion)

	if companionResult.IsErr() {
		logger.Error("Failed to find companion for soft delete: ", companionResult.Error())
		httpErr := types.ErrorInternal(
			"Failed to find companion for soft delete",
			companionResult.Error().Error(),
			"Companion ID: "+id,
		)
		return types.ResultErr[models.CompanionDB](&httpErr)
	}

	companion = companionResult.Value().(models.CompanionDB)
	err = configs.DB.UpdateOne(filter, companion)

	if err != nil {
		logger.Error("Failed to delete companion in MongoDB: ", err)
		httpErr := types.ErrorInternal(
			"Failed to delete companion",
			err.Error(),
			"Companion ID: "+id,
		)
		return types.ResultErr[models.CompanionDB](&httpErr)
	}

	return types.ResultOk(companion)
}
func (companionType) DeletePermanentByID(id string) types.Result[models.CompanionDB] {
	logger.Warning("Permanently deleting companion by ID: ", id)
	oid, err := models.ID.ToBson(id)
	if err != nil {
		logger.Error("Failed to convert ID to ObjectID: ", err)
		httpErr := types.Error(
			types.Http.C400().UnprocessableEntity(),
			"Invalid value",
			"Invalid ID format: "+err.Error(),
			"Companion ID: "+id,
		)
		return types.ResultErr[models.CompanionDB](&httpErr)
	}

	filter := bson.D{{Key: "_id", Value: oid}, {Key: "deleted_at", Value: bson.M{"$ne": time.Time{}}}} // Ensure the companion is marked as deleted
	var companion models.CompanionDB
	companionResult := configs.DB.FindOne(filter, companion)

	if companionResult.IsErr() {
		logger.Debug("Failed to find companion for permanent deletion: ", companionResult.Error())

		if companionResult.Error() == mongo.ErrNoDocuments {
			httpErr := types.ErrorNotFound(
				"Companion not found",
				"Companion with ID "+id+" not found or not marked as deleted",
			)
			return types.ResultErr[models.CompanionDB](&httpErr)
		}

		httpErr := types.ErrorInternal(
			"Failed to find companion for permanent deletion",
			err.Error(),
			"Companion ID: "+id,
		)
		return types.ResultErr[models.CompanionDB](&httpErr)
	}

	companion = companionResult.Value().(models.CompanionDB)
	err = configs.DB.PermanentDeleteOne(filter, companion)
	if err != nil {
		logger.Debug("Failed to permanently delete companion in MongoDB: ", err)
		httpErr := types.ErrorInternal(
			"Failed to permanently delete companion",
			err.Error(),
			"Companion ID: "+id,
		)
		return types.ResultErr[models.CompanionDB](&httpErr)
	}

	return types.ResultOk(companion)
}
func (companionType) DeletePermanentAll() types.Result[[]models.CompanionDB] {
	filter := bson.D{{Key: "deleted_at", Value: bson.M{"$ne": nil}}}

	var companion models.CompanionDB
	companionResult := configs.DB.FindAll(filter, companion)
	if companionResult.IsErr() {
		logger.Error("Failed to find companions for permanent deletion: ", companionResult.Error())
		httpErr := types.ErrorInternal(
			"Failed to find companions for permanent deletion",
			companionResult.Error().Error(),
			"Companion filter: "+filter.String(),
		)
		return types.ResultErr[[]models.CompanionDB](&httpErr)
	}

	err := configs.DB.PermanentDeleteMany(filter, companion)
	if err != nil {
		logger.Error("Failed to permanently delete all companions in MongoDB: ", err)
		httpErr := types.ErrorInternal(
			"Failed to permanently delete all companions",
			err.Error(),
		)
		return types.ResultErr[[]models.CompanionDB](&httpErr)
	}

	companions := utils.Map(companionResult.Value(), func(m models.DBModelInterface) models.CompanionDB {
		return m.(models.CompanionDB)
	})
	return types.ResultOk(companions)
}
