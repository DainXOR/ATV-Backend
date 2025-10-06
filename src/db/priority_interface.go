package db

import (
	"dainxor/atv/configs"
	"dainxor/atv/logger"
	"dainxor/atv/models"
	"dainxor/atv/types"
	"dainxor/atv/utils"
)

type priorityType struct{}

var Priority priorityType

func (priorityType) Create(u models.PriorityCreate) types.Result[models.PriorityDB] {
	priorityDB := u.ToInsert()
	resultID := configs.DB.InsertOne(priorityDB)

	if resultID.IsErr() {
		logger.Error("Error inserting priority: ", resultID.Error())
		return types.ResultErr[models.PriorityDB](resultID.Error())
	}

	priorityDB.ID = resultID.Value()
	return types.ResultOk(priorityDB)
}

func (priorityType) GetByID(id string, filter models.FilterObject) types.Result[models.PriorityDB] {
	oid, err := models.ID.ToDB(id)

	if err != nil {
		logger.Error("Failed to convert ID to ObjectID: ", err)
		httpErr := types.Error(
			types.Http.C400().UnprocessableEntity(),
			"Invalid value",
			"Invalid ID format: "+err.Error(),
			"Priority ID: "+id,
		)
		return types.ResultErr[models.PriorityDB](&httpErr)
	}

	filter = models.Filter.AddPart(filter, models.Filter.ID(oid))
	filter = models.Filter.AddPart(filter, models.Filter.NotDeleted())

	resultPriority := configs.DB.FindOne(filter, models.PriorityDB{})
	if resultPriority.IsErr() {
		logger.Warning("Failed to get priority by ID: ", resultPriority.Error())
		var httpErr types.HttpError

		switch resultPriority.Error() {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"Priority not found",
				"Priority with ID "+id+" not found",
			)

		default:
			logger.Error("Failed to get priority by ID: ", resultPriority.Error())
			httpErr = types.ErrorInternal(
				"Failed to retrieve priority",
				resultPriority.Error().Error(),
				"Priority ID: "+id,
			)
		}

		return types.ResultErr[models.PriorityDB](&httpErr)
	}

	return types.ResultOk(resultPriority.Value().(models.PriorityDB))
}
func (priorityType) GetAll(filter models.FilterObject) types.Result[[]models.PriorityDB] {
	filter = models.Filter.AddPart(filter, models.Filter.NotDeleted())

	resultSpecialities := configs.DB.FindAll(filter, models.PriorityDB{})
	if resultSpecialities.IsErr() {
		logger.Warning("Failed to get all specialities from database:", resultSpecialities.Error())
		var httpErr types.HttpError

		switch resultSpecialities.Error() {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"Specialities not found",
				"No specialities found",
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to retrieve specialities",
				resultSpecialities.Error().Error(),
			)
		}
		return types.ResultErr[[]models.PriorityDB](&httpErr)
	}

	specialities := utils.Map(resultSpecialities.Value(), models.InterfaceTo[models.PriorityDB])
	logger.Debug("Retrieved", len(specialities), "specialities from database")
	return types.ResultOk(specialities)
}

func (priorityType) UpdateByID(id string, priority models.PriorityCreate, filter models.FilterObject) types.Result[models.PriorityDB] {
	oid, err := models.ID.ToDB(id)
	if err != nil {
		logger.Warning("Failed to convert ID to ObjectID: ", err)
		httpErr := types.Error(
			types.Http.C400().UnprocessableEntity(),
			"Invalid value",
			"Invalid ID format: "+err.Error(),
			"Priority ID: "+id,
		)
		return types.ResultErr[models.PriorityDB](&httpErr)
	}

	priorityUpdate := priority.ToUpdate()

	filter = models.Filter.AddPart(filter, models.Filter.ID(oid))
	filter = models.Filter.AddPart(filter, models.Filter.NotDeleted())
	err = configs.DB.UpdateOne(filter, priorityUpdate)
	if err != nil {
		logger.Warning("Failed to update priority in database: ", err)
		var httpErr types.HttpError

		switch err {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"Priority not found",
				"Priority with ID "+id+" not found",
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to update priority",
				err.Error(),
				"Priority ID: "+id,
			)
		}
		return types.ResultErr[models.PriorityDB](&httpErr)
	}

	resultPriorityDB := configs.DB.FindOne(filter, models.PriorityDB{})
	if resultPriorityDB.IsErr() {
		logger.Warning("Failed to retrieve updated priority: ", resultPriorityDB.Error())
		var httpErr types.HttpError

		switch resultPriorityDB.Error() {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"Priority not found",
				"Priority with ID "+id+" not found",
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to retrieve updated priority",
				resultPriorityDB.Error().Error(),
				"Priority ID: "+id,
			)
		}
		return types.ResultErr[models.PriorityDB](&httpErr)
	}

	return types.ResultOk(resultPriorityDB.Value().(models.PriorityDB))
}

func (priorityType) PatchByID(id string, priority models.PriorityCreate, filter models.FilterObject) types.Result[models.PriorityDB] {
	oid, err := models.ID.ToDB(id)
	if err != nil {
		logger.Warning("Failed to convert ID to ObjectID: ", err)
		httpErr := types.Error(
			types.Http.C400().UnprocessableEntity(),
			"Invalid value",
			"Invalid ID format: "+err.Error(),
			"Priority ID: "+id,
		)
		return types.ResultErr[models.PriorityDB](&httpErr)
	}

	priorityPatch := priority.ToUpdate()

	filter = models.Filter.AddPart(filter, models.Filter.ID(oid))
	filter = models.Filter.AddPart(filter, models.Filter.NotDeleted())
	err = configs.DB.PatchOne(filter, priorityPatch)
	if err != nil {
		logger.Warning("Failed to patch priority in database: ", err)
		var httpErr types.HttpError

		switch err {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"Priority not found",
				"Priority with ID "+id+" not found",
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to patch priority",
				err.Error(),
				"Priority ID: "+id,
			)
		}
		return types.ResultErr[models.PriorityDB](&httpErr)
	}

	resultPriorityDB := configs.DB.FindOne(filter, models.PriorityDB{})
	if resultPriorityDB.IsErr() {
		logger.Warning("Failed to retrieve updated priority: ", resultPriorityDB.Error())
		var httpErr types.HttpError

		switch resultPriorityDB.Error() {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"Priority not found",
				"Priority with ID "+id+" not found",
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to retrieve updated priority",
				resultPriorityDB.Error().Error(),
				"Priority ID: "+id,
			)
		}
		return types.ResultErr[models.PriorityDB](&httpErr)
	}

	return types.ResultOk(resultPriorityDB.Value().(models.PriorityDB))
}

func (priorityType) DeleteByID(id string, filter models.FilterObject) types.Result[models.PriorityDB] {
	oid, err := models.ID.ToDB(id)
	if err != nil {
		logger.Warning("Failed to convert ID to ObjectID: ", err)
		httpErr := types.Error(
			types.Http.C400().UnprocessableEntity(),
			"Invalid value",
			"Invalid ID format: "+err.Error(),
			"Priority ID: "+id,
		)
		return types.ResultErr[models.PriorityDB](&httpErr)
	}

	filter = models.Filter.AddPart(filter, models.Filter.ID(oid))
	filter = models.Filter.AddPart(filter, models.Filter.NotDeleted())
	resultPriority := configs.DB.FindOne(filter, models.PriorityDB{})
	if resultPriority.IsErr() {
		logger.Warning("Failed to retrieve priority: ", resultPriority.Error())
		var httpErr types.HttpError

		switch resultPriority.Error() {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"Priority not found",
				"Priority with ID "+id+" not found",
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to retrieve priority",
				resultPriority.Error().Error(),
				"Priority ID: "+id,
			)
		}

		return types.ResultErr[models.PriorityDB](&httpErr)
	}

	deletedPriority := resultPriority.Value().(models.PriorityDB)
	err = configs.DB.SoftDeleteOne(filter, deletedPriority)
	if err != nil {
		logger.Warning("Failed to delete priority in database: ", err)
		var httpErr types.HttpError

		switch err {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"Priority not found",
				"Priority with ID "+id+" not found",
			)

		case configs.DBErr.NotModified():
			httpErr = types.Error(
				types.Http.C400().PreconditionFailed(),
				"Priority was already marked as deleted",
				"Priority ID: "+id,
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to delete priority",
				err.Error(),
				"Priority ID: "+id,
			)
		}

		return types.ResultErr[models.PriorityDB](&httpErr)
	}

	return types.ResultOk(deletedPriority)
}
func (priorityType) DeletePermanentByID(id string, filter models.FilterObject) types.Result[models.PriorityDB] {
	logger.Warning("Permanently deleting priority by ID: ", id)
	oid, err := models.ID.ToDB(id)
	if err != nil {
		logger.Warning("Failed to convert ID to ObjectID: ", err)
		httpErr := types.Error(
			types.Http.C400().UnprocessableEntity(),
			"Invalid value",
			"Invalid ID format: "+err.Error(),
			"Priority ID: "+id,
		)
		return types.ResultErr[models.PriorityDB](&httpErr)
	}

	filter = models.Filter.AddPart(filter, models.Filter.ID(oid))
	filter = models.Filter.AddPart(filter, models.Filter.Deleted())
	resultPriority := configs.DB.FindOne(filter, models.PriorityDB{})
	if resultPriority.IsErr() {
		logger.Warning("Failed to find priority for permanent deletion: ", resultPriority.Error())
		var httpErr types.HttpError

		switch resultPriority.Error() {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"Priority not found",
				"Priority with ID "+id+" not found or not marked as deleted",
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to find priority for permanent deletion",
				resultPriority.Error().Error(),
				"Priority ID: "+id,
			)
		}

		return types.ResultErr[models.PriorityDB](&httpErr)
	}

	err = configs.DB.PermanentDeleteOne(filter, models.PriorityDB{})
	if err != nil {
		logger.Warning("Failed to permanently delete priority in database: ", err)
		var httpErr types.HttpError

		switch err {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"Priority not found",
				"Priority with ID "+id+" not found or not marked as deleted",
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to permanently delete priority",
				err.Error(),
				"Priority ID: "+id,
			)
		}

		return types.ResultErr[models.PriorityDB](&httpErr)
	}

	return types.ResultOk(resultPriority.Value().(models.PriorityDB))
}
func (priorityType) DeletePermanentAll(filter models.FilterObject) types.Result[[]models.PriorityDB] {
	filter = models.Filter.AddPart(filter, models.Filter.Deleted())

	resultPrioritys := configs.DB.FindAll(filter, models.PriorityDB{})
	if resultPrioritys.IsErr() {
		logger.Warning("Failed to find priority for permanent deletion: ", resultPrioritys.Error())
		var httpErr types.HttpError

		switch resultPrioritys.Error() {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"Prioritys not found",
				"No prioritys marked as deleted found",
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to find priority for permanent deletion",
				resultPrioritys.Error().Error(),
			)
		}

		return types.ResultErr[[]models.PriorityDB](&httpErr)
	}

	err := configs.DB.PermanentDeleteMany(filter, models.PriorityDB{})
	if err != nil {
		logger.Warning("Failed to permanently delete all prioritys in database: ", err)
		var httpErr types.HttpError

		switch err {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"Prioritys not found",
				"No prioritys marked as deleted found",
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to permanently delete all prioritys",
				err.Error(),
			)
		}
		return types.ResultErr[[]models.PriorityDB](&httpErr)
	}

	prioritys := utils.Map(resultPrioritys.Value(), models.InterfaceTo[models.PriorityDB])
	return types.ResultOk(prioritys)
}
