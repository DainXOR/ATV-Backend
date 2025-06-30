package db

import (
	"dainxor/atv/configs"
	"dainxor/atv/logger"
	"dainxor/atv/models"
	"dainxor/atv/types"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type specialityType struct{}

var Speciality specialityType

func (specialityType) Create(u models.SpecialityCreate) types.Result[models.SpecialityDBMongo] {
	specialityDB := u.ToInsert()
	result, err := configs.DB.InsertOne(specialityDB)

	if err != nil {
		logger.Error("Error inserting speciality: ", err)
		return types.ResultErr[models.SpecialityDBMongo](err)
	}

	specialityDB.ID, err = models.DBIDFrom(result.InsertedID)

	if err != nil {
		logger.Error("Error converting inserted ID to PrimitiveID: ", err)
		httpErr := types.ErrorInternal(
			"Failed to create speciality",
			"Failed to convert inserted ID to PrimitiveID",
			"Error: "+err.Error(),
		)
		return types.ResultErr[models.SpecialityDBMongo](&httpErr)
	}

	return types.ResultOk(specialityDB)
}

func (specialityType) GetByID(id string) types.Result[models.SpecialityDBMongo] {
	oid, err := models.BsonIDFrom(id)

	if err != nil {
		logger.Error("Failed to convert ID to ObjectID: ", err)
		httpErr := types.Error(
			types.Http.UnprocessableEntity(),
			"Invalid value",
			"Invalid ID format: "+err.Error(),
			"Speciality ID: "+id,
		)
		return types.ResultErr[models.SpecialityDBMongo](&httpErr)
	}

	filter := bson.D{{Key: "_id", Value: oid}}
	var speciality models.SpecialityDBMongo

	err = configs.DB.FindOne(filter, &speciality)
	if err != nil {
		var httpErr types.HttpError

		if err == mongo.ErrNoDocuments {
			logger.Error("Failed to get speciality by ID: ", err)
			httpErr = types.ErrorNotFound(
				"Speciality not found",
				"Speciality with ID "+id+" not found",
			)
		} else {
			logger.Error("Failed to get speciality by ID: ", err)
			httpErr = types.ErrorInternal(
				"Failed to retrieve speciality",
				"Decoding error",
				err.Error(),
				"Speciality ID: "+id,
			)
		}

		return types.ResultErr[models.SpecialityDBMongo](&httpErr)
	}

	return types.ResultOk(speciality)
}
func (specialityType) GetAll() types.Result[[]models.SpecialityDBMongo] {
	filter := bson.D{{Key: "deleted_at", Value: nil}} // Filter to exclude deleted specialities
	specialities := []models.SpecialityDBMongo{}

	err := configs.DB.FindAll(filter, &specialities)
	if err != nil {
		logger.Error("Failed to get all specialities from MongoDB:", err)
		httpErr := types.ErrorInternal(
			"Failed to retrieve specialities",
			err.Error(),
		)

		return types.ResultErr[[]models.SpecialityDBMongo](&httpErr)
	}

	logger.Debug("Retrieved", len(specialities), "specialities from MongoDB database")
	return types.ResultOk(specialities)
}
