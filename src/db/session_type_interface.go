package db

import (
	"dainxor/atv/configs"
	"dainxor/atv/logger"
	"dainxor/atv/models"
	"dainxor/atv/types"
	"dainxor/atv/utils"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type sessionTypeType struct{}

var SessionType sessionTypeType

func (sessionTypeType) Create(u models.SessionTypeCreate) types.Result[models.SessionTypeDB] {
	sessionTypeDB := u.ToInsert()
	resultID := configs.DB.InsertOne(sessionTypeDB)

	if resultID.IsErr() {
		logger.Warning("Error inserting session type: ", resultID.Error())
		return types.ResultErr[models.SessionTypeDB](resultID.Error())
	}

	sessionTypeDB.ID = resultID.Value()
	return types.ResultOk(sessionTypeDB)
}

func (sessionTypeType) GetByID(id string) types.Result[models.SessionTypeDB] {
	oid, err := models.ID.ToDB(id)

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
	resultSessionType := configs.DB.FindOne(filter, models.SessionTypeDB{})
	if resultSessionType.IsErr() {
		logger.Warning("Failed to get session type by ID: ", resultSessionType.Error())
		var httpErr types.HttpError

		switch resultSessionType.Error() {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"SessionType not found",
				"SessionType with ID "+id+" not found",
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to retrieve session type",
				resultSessionType.Error().Error(),
				"SessionType ID: "+id,
			)
		}

		return types.ResultErr[models.SessionTypeDB](&httpErr)
	}

	return types.ResultOk(resultSessionType.Value().(models.SessionTypeDB))
}
func (sessionTypeType) GetAll() types.Result[[]models.SessionTypeDB] {
	filter := bson.D{{Key: "deleted_at", Value: models.Time.Zero()}} // Filter to exclude deleted session types

	resultSessionTypes := configs.DB.FindAll(filter, models.SessionTypeDB{})
	if resultSessionTypes.IsErr() {
		logger.Warning("Failed to get all session types from MongoDB:", resultSessionTypes.Error())
		return types.ResultErr[[]models.SessionTypeDB](resultSessionTypes.Error())
	}

	sessionTypes := utils.Map(resultSessionTypes.Value(), models.InterfaceTo[models.SessionTypeDB])
	logger.Debug("Retrieved", len(sessionTypes), "session types from MongoDB database")
	return types.ResultOk(sessionTypes)
}
