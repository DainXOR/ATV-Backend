package dao

import (
	"dainxor/atv/configs"
	"dainxor/atv/logger"
	"dainxor/atv/models"
	"dainxor/atv/types"
	"dainxor/atv/utils"
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

func (sessionTypeType) GetByID(id string, filter models.FilterObject) types.Result[models.SessionTypeDB] {
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

	filter = models.Filter.AddPart(filter, models.Filter.ID(oid))
	filter = models.Filter.AddPart(filter, models.Filter.NotDeleted())
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
func (sessionTypeType) GetAll(filter models.FilterObject) types.Result[[]models.SessionTypeDB] {
	filter = models.Filter.AddPart(filter, models.Filter.NotDeleted())

	resultSessionTypes := configs.DB.FindAll(filter, models.SessionTypeDB{})
	if resultSessionTypes.IsErr() {
		logger.Warning("Failed to get all session types from MongoDB:", resultSessionTypes.Error())
		return types.ResultErr[[]models.SessionTypeDB](resultSessionTypes.Error())
	}

	sessionTypes := utils.Map(resultSessionTypes.Value(), models.InterfaceTo[models.SessionTypeDB])
	logger.Debug("Retrieved", len(sessionTypes), "session types from MongoDB database")
	return types.ResultOk(sessionTypes)
}

func (sessionTypeType) UpdateByID(id string, sessionType models.SessionTypeCreate, filter models.FilterObject) types.Result[models.SessionTypeDB] {
	oid, err := models.ID.ToDB(id)
	if err != nil {
		logger.Warning("Failed to convert ID to ObjectID: ", err)
		httpErr := types.Error(
			types.Http.C400().UnprocessableEntity(),
			"Invalid value",
			"Invalid ID format: "+err.Error(),
			"SessionType ID: "+id,
		)
		return types.ResultErr[models.SessionTypeDB](&httpErr)
	}

	sessionTypeUpdate := sessionType.ToUpdate()

	filter = models.Filter.AddPart(filter, models.Filter.ID(oid))
	filter = models.Filter.AddPart(filter, models.Filter.NotDeleted())
	err = configs.DB.UpdateOne(filter, sessionTypeUpdate)
	if err != nil {
		logger.Warning("Failed to update sessionType in database: ", err)
		var httpErr types.HttpError

		switch err {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"SessionType not found",
				"SessionType with ID "+id+" not found",
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to update sessionType",
				err.Error(),
				"SessionType ID: "+id,
			)
		}
		return types.ResultErr[models.SessionTypeDB](&httpErr)
	}

	resultSessionTypeDB := configs.DB.FindOne(filter, models.SessionTypeDB{})
	if resultSessionTypeDB.IsErr() {
		logger.Warning("Failed to retrieve updated sessionType: ", resultSessionTypeDB.Error())
		var httpErr types.HttpError

		switch resultSessionTypeDB.Error() {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"SessionType not found",
				"SessionType with ID "+id+" not found",
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to retrieve updated sessionType",
				resultSessionTypeDB.Error().Error(),
				"SessionType ID: "+id,
			)
		}
		return types.ResultErr[models.SessionTypeDB](&httpErr)
	}

	return types.ResultOk(resultSessionTypeDB.Value().(models.SessionTypeDB))
}

func (sessionTypeType) PatchByID(id string, sessionType models.SessionTypeCreate, filter models.FilterObject) types.Result[models.SessionTypeDB] {
	oid, err := models.ID.ToDB(id)
	if err != nil {
		logger.Warning("Failed to convert ID to ObjectID: ", err)
		httpErr := types.Error(
			types.Http.C400().UnprocessableEntity(),
			"Invalid value",
			"Invalid ID format: "+err.Error(),
			"SessionType ID: "+id,
		)
		return types.ResultErr[models.SessionTypeDB](&httpErr)
	}

	sessionTypePatch := sessionType.ToUpdate()

	filter = models.Filter.AddPart(filter, models.Filter.ID(oid))
	filter = models.Filter.AddPart(filter, models.Filter.NotDeleted())
	err = configs.DB.PatchOne(filter, sessionTypePatch)
	if err != nil {
		logger.Warning("Failed to patch sessionType in database: ", err)
		var httpErr types.HttpError

		switch err {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"SessionType not found",
				"SessionType with ID "+id+" not found",
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to patch sessionType",
				err.Error(),
				"SessionType ID: "+id,
			)
		}
		return types.ResultErr[models.SessionTypeDB](&httpErr)
	}

	resultSessionTypeDB := configs.DB.FindOne(filter, models.SessionTypeDB{})
	if resultSessionTypeDB.IsErr() {
		logger.Warning("Failed to retrieve updated sessionType: ", resultSessionTypeDB.Error())
		var httpErr types.HttpError

		switch resultSessionTypeDB.Error() {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"SessionType not found",
				"SessionType with ID "+id+" not found",
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to retrieve updated sessionType",
				resultSessionTypeDB.Error().Error(),
				"SessionType ID: "+id,
			)
		}
		return types.ResultErr[models.SessionTypeDB](&httpErr)
	}

	return types.ResultOk(resultSessionTypeDB.Value().(models.SessionTypeDB))
}

func (sessionTypeType) DeleteByID(id string, filter models.FilterObject) types.Result[models.SessionTypeDB] {
	oid, err := models.ID.ToDB(id)
	if err != nil {
		logger.Warning("Failed to convert ID to ObjectID: ", err)
		httpErr := types.Error(
			types.Http.C400().UnprocessableEntity(),
			"Invalid value",
			"Invalid ID format: "+err.Error(),
			"SessionType ID: "+id,
		)
		return types.ResultErr[models.SessionTypeDB](&httpErr)
	}

	filter = models.Filter.AddPart(filter, models.Filter.ID(oid))
	filter = models.Filter.AddPart(filter, models.Filter.NotDeleted())
	resultSessionType := configs.DB.FindOne(filter, models.SessionTypeDB{})
	if resultSessionType.IsErr() {
		logger.Warning("Failed to retrieve sessionType: ", resultSessionType.Error())
		var httpErr types.HttpError

		switch resultSessionType.Error() {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"SessionType not found",
				"SessionType with ID "+id+" not found",
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to retrieve sessionType",
				resultSessionType.Error().Error(),
				"SessionType ID: "+id,
			)
		}

		return types.ResultErr[models.SessionTypeDB](&httpErr)
	}

	deletedSessionType := resultSessionType.Value().(models.SessionTypeDB)
	err = configs.DB.SoftDeleteOne(filter, deletedSessionType)
	if err != nil {
		logger.Warning("Failed to delete sessionType in database: ", err)
		var httpErr types.HttpError

		switch err {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"SessionType not found",
				"SessionType with ID "+id+" not found",
			)

		case configs.DBErr.NotModified():
			httpErr = types.Error(
				types.Http.C400().PreconditionFailed(),
				"SessionType was already marked as deleted",
				"SessionType ID: "+id,
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to delete sessionType",
				err.Error(),
				"SessionType ID: "+id,
			)
		}

		return types.ResultErr[models.SessionTypeDB](&httpErr)
	}

	return types.ResultOk(deletedSessionType)
}
func (sessionTypeType) DeletePermanentByID(id string, filter models.FilterObject) types.Result[models.SessionTypeDB] {
	logger.Warning("Permanently deleting sessionType by ID: ", id)
	oid, err := models.ID.ToDB(id)
	if err != nil {
		logger.Warning("Failed to convert ID to ObjectID: ", err)
		httpErr := types.Error(
			types.Http.C400().UnprocessableEntity(),
			"Invalid value",
			"Invalid ID format: "+err.Error(),
			"SessionType ID: "+id,
		)
		return types.ResultErr[models.SessionTypeDB](&httpErr)
	}

	filter = models.Filter.AddPart(filter, models.Filter.ID(oid))
	filter = models.Filter.AddPart(filter, models.Filter.Deleted())
	resultSessionType := configs.DB.FindOne(filter, models.SessionTypeDB{})
	if resultSessionType.IsErr() {
		logger.Warning("Failed to find sessionType for permanent deletion: ", resultSessionType.Error())
		var httpErr types.HttpError

		switch resultSessionType.Error() {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"SessionType not found",
				"SessionType with ID "+id+" not found or not marked as deleted",
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to find sessionType for permanent deletion",
				resultSessionType.Error().Error(),
				"SessionType ID: "+id,
			)
		}

		return types.ResultErr[models.SessionTypeDB](&httpErr)
	}

	err = configs.DB.PermanentDeleteOne(filter, models.SessionTypeDB{})
	if err != nil {
		logger.Warning("Failed to permanently delete sessionType in database: ", err)
		var httpErr types.HttpError

		switch err {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"SessionType not found",
				"SessionType with ID "+id+" not found or not marked as deleted",
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to permanently delete sessionType",
				err.Error(),
				"SessionType ID: "+id,
			)
		}

		return types.ResultErr[models.SessionTypeDB](&httpErr)
	}

	return types.ResultOk(resultSessionType.Value().(models.SessionTypeDB))
}
func (sessionTypeType) DeletePermanentAll(filter models.FilterObject) types.Result[[]models.SessionTypeDB] {
	filter = models.Filter.AddPart(filter, models.Filter.Deleted())

	resultSessionTypes := configs.DB.FindAll(filter, models.SessionTypeDB{})
	if resultSessionTypes.IsErr() {
		logger.Warning("Failed to find sessionType for permanent deletion: ", resultSessionTypes.Error())
		var httpErr types.HttpError

		switch resultSessionTypes.Error() {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"SessionTypes not found",
				"No sessionTypes marked as deleted found",
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to find sessionType for permanent deletion",
				resultSessionTypes.Error().Error(),
			)
		}

		return types.ResultErr[[]models.SessionTypeDB](&httpErr)
	}

	err := configs.DB.PermanentDeleteMany(filter, models.SessionTypeDB{})
	if err != nil {
		logger.Warning("Failed to permanently delete all sessionTypes in database: ", err)
		var httpErr types.HttpError

		switch err {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"SessionTypes not found",
				"No sessionTypes marked as deleted found",
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to permanently delete all sessionTypes",
				err.Error(),
			)
		}
		return types.ResultErr[[]models.SessionTypeDB](&httpErr)
	}

	sessionTypes := utils.Map(resultSessionTypes.Value(), models.InterfaceTo[models.SessionTypeDB])
	return types.ResultOk(sessionTypes)
}
