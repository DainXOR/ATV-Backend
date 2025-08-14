package db

import (
	"dainxor/atv/configs"
	"dainxor/atv/logger"
	"dainxor/atv/models"
	"dainxor/atv/types"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type universityType struct{}

var University universityType

func (universityType) Create(u models.UniversityCreate) types.Result[models.UniversityDB] {
	universityDB := u.ToInsert()
	result, err := configs.DB.InsertOne(&universityDB)

	if err != nil {
		logger.Error("Error inserting university: ", err)
		return types.ResultErr[models.UniversityDB](err)
	}

	universityDB.ID, err = models.ID.ToDB(result.InsertedID)

	if err != nil {
		logger.Error("Error converting inserted ID to PrimitiveID: ", err)
		httpErr := types.ErrorInternal(
			"Failed to create university",
			"Failed to convert inserted ID to PrimitiveID",
			"Error: "+err.Error(),
		)
		return types.ResultErr[models.UniversityDB](&httpErr)
	}

	return types.ResultOk(universityDB)
}

func (universityType) GetByID(id string) types.Result[models.UniversityDB] {
	oid, err := models.ID.ToDB(id)

	if err != nil {
		logger.Error("Failed to convert ID to ObjectID: ", err)
		httpErr := types.Error(
			types.Http.C400().UnprocessableEntity(),
			"Invalid value",
			"Invalid ID format: "+err.Error(),
			"University ID: "+id,
		)
		return types.ResultErr[models.UniversityDB](&httpErr)
	}

	filter := bson.D{{Key: "_id", Value: oid}}
	var university models.UniversityDB

	err = configs.DB.FindOne(filter, &university)
	if err != nil {
		var httpErr types.HttpError

		if err == mongo.ErrNoDocuments {
			logger.Error("Failed to get university by ID: ", err)
			httpErr = types.ErrorNotFound(
				"University not found",
				"University with ID "+id+" not found",
			)
		} else {
			logger.Error("Failed to get university by ID: ", err)
			httpErr = types.ErrorInternal(
				"Failed to retrieve university",
				"Decoding error",
				err.Error(),
				"University ID: "+id,
			)
		}

		return types.ResultErr[models.UniversityDB](&httpErr)
	}

	return types.ResultOk(university)
}
func (universityType) GetAll() types.Result[[]models.UniversityDB] {
	filter := bson.D{{Key: "deleted_at", Value: nil}} // Filter to exclude deleted universities
	universities := []models.UniversityDB{}

	err := configs.DB.FindAll(filter, &universities)
	if err != nil {
		logger.Error("Failed to get all universities from MongoDB:", err)
		httpErr := types.ErrorInternal(
			"Failed to retrieve universities",
			err.Error(),
		)

		return types.ResultErr[[]models.UniversityDB](&httpErr)
	}

	logger.Debug("Retrieved", len(universities), "universities from MongoDB database")
	return types.ResultOk(universities)
}
