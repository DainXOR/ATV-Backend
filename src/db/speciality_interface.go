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

func (specialityType) Create(u models.SpecialityCreate) types.Result[models.SpecialityDB] {
	specialityDB := u.ToInsert()
	result, err := configs.DB.InsertOne(&specialityDB)

	if err != nil {
		logger.Error("Error inserting speciality: ", err)
		return types.ResultErr[models.SpecialityDB](err)
	}

	specialityDB.ID, err = models.DBIDFrom(result.InsertedID)

	if err != nil {
		logger.Error("Error converting inserted ID to PrimitiveID: ", err)
		httpErr := types.ErrorInternal(
			"Failed to create speciality",
			"Failed to convert inserted ID to PrimitiveID",
			"Error: "+err.Error(),
		)
		return types.ResultErr[models.SpecialityDB](&httpErr)
	}

	return types.ResultOk(specialityDB)
}

func (specialityType) GetByID(id string) types.Result[models.SpecialityDB] {
	oid, err := models.ID.ToBson(id)

	if err != nil {
		logger.Error("Failed to convert ID to ObjectID: ", err)
		httpErr := types.Error(
			types.Http.C400().UnprocessableEntity(),
			"Invalid value",
			"Invalid ID format: "+err.Error(),
			"Speciality ID: "+id,
		)
		return types.ResultErr[models.SpecialityDB](&httpErr)
	}

	filter := bson.D{{Key: "_id", Value: oid}}
	var speciality models.SpecialityDB

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

		return types.ResultErr[models.SpecialityDB](&httpErr)
	}

	return types.ResultOk(speciality)
}
func (specialityType) GetAll() types.Result[[]models.SpecialityDB] {
	filter := bson.D{{Key: "deleted_at", Value: nil}} // Filter to exclude deleted specialities
	specialities := []models.SpecialityDB{}

	err := configs.DB.FindAll(filter, &specialities)
	if err != nil {
		logger.Error("Failed to get all specialities from MongoDB:", err)
		httpErr := types.ErrorInternal(
			"Failed to retrieve specialities",
			err.Error(),
		)

		return types.ResultErr[[]models.SpecialityDB](&httpErr)
	}

	logger.Debug("Retrieved", len(specialities), "specialities from MongoDB database")
	return types.ResultOk(specialities)
}
