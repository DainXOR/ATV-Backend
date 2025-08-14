package db

import (
	"dainxor/atv/configs"
	"dainxor/atv/logger"
	"dainxor/atv/models"
	"dainxor/atv/types"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type sessionTypeType struct{}

var SessionType sessionTypeType

func (sessionTypeType) Create(u models.SessionTypeCreate) types.Result[models.SessionTypeDB] {
	sessionTypeDB := u.ToInsert()
	result, err := configs.DB.InsertOne(&sessionTypeDB)

	if err != nil {
		logger.Error("Error inserting session type: ", err)
		return types.ResultErr[models.SessionTypeDB](err)
	}

	sessionTypeDB.ID, err = models.DBIDFrom(result.InsertedID)

	if err != nil {
		logger.Error("Error converting inserted ID to PrimitiveID: ", err)
		httpErr := types.ErrorInternal(
			"Failed to create session type",
			"Failed to convert inserted ID to PrimitiveID",
			"Error: "+err.Error(),
		)
		return types.ResultErr[models.SessionTypeDB](&httpErr)
	}

	return types.ResultOk(sessionTypeDB)
}

func (sessionTypeType) GetByID(id string) types.Result[models.SessionTypeDB] {
	oid, err := models.BsonIDFrom(id)

	if err != nil {
		logger.Error("Failed to convert ID to ObjectID: ", err)
		httpErr := types.Error(
			types.Http.C400().UnprocessableEntity(),
			"Invalid value",
			"Invalid ID format: "+err.Error(),
			"SessionType ID: "+id,
		)
		return types.ResultErr[models.SessionTypeDB](&httpErr)
	}

	filter := bson.D{{Key: "_id", Value: oid}}
	var sessionType models.SessionTypeDB

	err = configs.DB.FindOne(filter, &sessionType)
	if err != nil {
		var httpErr types.HttpError

		if err == mongo.ErrNoDocuments {
			logger.Error("Failed to get session type by ID: ", err)
			httpErr = types.ErrorNotFound(
				"SessionType not found",
				"SessionType with ID "+id+" not found",
			)
		} else {
			logger.Error("Failed to get session type by ID: ", err)
			httpErr = types.ErrorInternal(
				"Failed to retrieve session type",
				"Decoding error",
				err.Error(),
				"SessionType ID: "+id,
			)
		}

		return types.ResultErr[models.SessionTypeDB](&httpErr)
	}

	return types.ResultOk(sessionType)
}
func (sessionTypeType) GetAll() types.Result[[]models.SessionTypeDB] {
	filter := bson.D{{Key: "deleted_at", Value: models.Time.Zero()}} // Filter to exclude deleted session types
	sessionTypes := []models.SessionTypeDB{}

	err := configs.DB.FindAll(filter, &sessionTypes)
	if err != nil {
		logger.Error("Failed to get all session types from MongoDB:", err)
		httpErr := types.ErrorInternal(
			"Failed to retrieve session types",
			err.Error(),
		)

		return types.ResultErr[[]models.SessionTypeDB](&httpErr)
	}

	logger.Debug("Retrieved", len(sessionTypes), "session types from MongoDB database")
	return types.ResultOk(sessionTypes)
}
