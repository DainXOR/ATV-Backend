package db

import (
	"dainxor/atv/configs"
	"dainxor/atv/logger"
	"dainxor/atv/models"
	"dainxor/atv/types"
	"dainxor/atv/utils"
)

type universityType struct{}

var University universityType

func (universityType) Create(u models.UniversityCreate) types.Result[models.UniversityDB] {
	universityDB := u.ToInsert()
	resultID := configs.DB.InsertOne(universityDB)

	if resultID.IsErr() {
		logger.Warning("Error inserting university: ", resultID.Error())
		return types.ResultErr[models.UniversityDB](resultID.Error())
	}

	universityDB.ID = resultID.Value()
	return types.ResultOk(universityDB)
}

func (universityType) GetByID(id string, filter models.FilterObject) types.Result[models.UniversityDB] {
	oid, err := models.ID.ToDB(id)

	if err != nil {
		logger.Warning("Failed to convert ID to ObjectID: ", err)
		httpErr := types.Error(
			types.Http.C400().UnprocessableEntity(),
			"Invalid value",
			"Invalid ID format: "+err.Error(),
			"University ID: "+id,
		)
		return types.ResultErr[models.UniversityDB](&httpErr)
	}

	filter = models.Filter.AddPart(filter, models.Filter.ID(oid))
	filter = models.Filter.AddPart(filter, models.Filter.NotDeleted())

	resultUniversity := configs.DB.FindOne(filter, models.UniversityDB{})
	if resultUniversity.IsErr() {
		logger.Warning("Failed to get university by ID: ", resultUniversity.Error())
		var httpErr types.HttpError

		switch resultUniversity.Error() {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"University not found",
				"University with ID "+id+" not found",
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to retrieve university",
				resultUniversity.Error().Error(),
				"University ID: "+id,
			)
		}

		return types.ResultErr[models.UniversityDB](&httpErr)
	}

	return types.ResultOk(resultUniversity.Value().(models.UniversityDB))
}
func (universityType) GetAll(filter models.FilterObject) types.Result[[]models.UniversityDB] {
	filter = models.Filter.AddPart(filter, models.Filter.NotDeleted())

	resultUniversities := configs.DB.FindAll(filter, models.UniversityDB{})
	if resultUniversities.IsErr() {
		logger.Warning("Failed to get all universities from database:", resultUniversities.Error())
		var httpErr types.HttpError

		switch resultUniversities.Error() {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"Universities not found",
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to retrieve universities",
				resultUniversities.Error().Error(),
			)

		}

		return types.ResultErr[[]models.UniversityDB](&httpErr)
	}

	universities := utils.Map(resultUniversities.Value(), models.InterfaceTo[models.UniversityDB])
	logger.Debug("Retrieved", len(universities), "universities from database")
	return types.ResultOk(universities)
}

func (universityType) UpdateByID(id string, university models.UniversityCreate, filter models.FilterObject) types.Result[models.UniversityDB] {
	oid, err := models.ID.ToDB(id)
	if err != nil {
		logger.Warning("Failed to convert ID to ObjectID: ", err)
		httpErr := types.Error(
			types.Http.C400().UnprocessableEntity(),
			"Invalid value",
			"Invalid ID format: "+err.Error(),
			"University ID: "+id,
		)
		return types.ResultErr[models.UniversityDB](&httpErr)
	}

	universityUpdate := university.ToUpdate()

	filter = models.Filter.AddPart(filter, models.Filter.ID(oid))
	filter = models.Filter.AddPart(filter, models.Filter.NotDeleted())
	err = configs.DB.UpdateOne(filter, universityUpdate)
	if err != nil {
		logger.Warning("Failed to update university in database: ", err)
		var httpErr types.HttpError

		switch err {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"University not found",
				"University with ID "+id+" not found",
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to update university",
				err.Error(),
				"University ID: "+id,
			)
		}
		return types.ResultErr[models.UniversityDB](&httpErr)
	}

	resultUniversityDB := configs.DB.FindOne(filter, models.UniversityDB{})
	if resultUniversityDB.IsErr() {
		logger.Warning("Failed to retrieve updated university: ", resultUniversityDB.Error())
		var httpErr types.HttpError

		switch resultUniversityDB.Error() {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"University not found",
				"University with ID "+id+" not found",
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to retrieve updated university",
				resultUniversityDB.Error().Error(),
				"University ID: "+id,
			)
		}
		return types.ResultErr[models.UniversityDB](&httpErr)
	}

	return types.ResultOk(resultUniversityDB.Value().(models.UniversityDB))
}

func (universityType) PatchByID(id string, university models.UniversityCreate, filter models.FilterObject) types.Result[models.UniversityDB] {
	oid, err := models.ID.ToDB(id)
	if err != nil {
		logger.Warning("Failed to convert ID to ObjectID: ", err)
		httpErr := types.Error(
			types.Http.C400().UnprocessableEntity(),
			"Invalid value",
			"Invalid ID format: "+err.Error(),
			"University ID: "+id,
		)
		return types.ResultErr[models.UniversityDB](&httpErr)
	}

	universityPatch := university.ToUpdate()

	filter = models.Filter.AddPart(filter, models.Filter.ID(oid))
	filter = models.Filter.AddPart(filter, models.Filter.NotDeleted())
	err = configs.DB.PatchOne(filter, universityPatch)
	if err != nil {
		logger.Warning("Failed to patch university in database: ", err)
		var httpErr types.HttpError

		switch err {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"University not found",
				"University with ID "+id+" not found",
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to patch university",
				err.Error(),
				"University ID: "+id,
			)
		}
		return types.ResultErr[models.UniversityDB](&httpErr)
	}

	resultUniversityDB := configs.DB.FindOne(filter, models.UniversityDB{})
	if resultUniversityDB.IsErr() {
		logger.Warning("Failed to retrieve updated university: ", resultUniversityDB.Error())
		var httpErr types.HttpError

		switch resultUniversityDB.Error() {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"University not found",
				"University with ID "+id+" not found",
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to retrieve updated university",
				resultUniversityDB.Error().Error(),
				"University ID: "+id,
			)
		}
		return types.ResultErr[models.UniversityDB](&httpErr)
	}

	return types.ResultOk(resultUniversityDB.Value().(models.UniversityDB))
}

func (universityType) DeleteByID(id string, filter models.FilterObject) types.Result[models.UniversityDB] {
	oid, err := models.ID.ToDB(id)
	if err != nil {
		logger.Warning("Failed to convert ID to ObjectID: ", err)
		httpErr := types.Error(
			types.Http.C400().UnprocessableEntity(),
			"Invalid value",
			"Invalid ID format: "+err.Error(),
			"University ID: "+id,
		)
		return types.ResultErr[models.UniversityDB](&httpErr)
	}

	filter = models.Filter.AddPart(filter, models.Filter.ID(oid))
	filter = models.Filter.AddPart(filter, models.Filter.NotDeleted())
	resultUniversity := configs.DB.FindOne(filter, models.UniversityDB{})
	if resultUniversity.IsErr() {
		logger.Warning("Failed to retrieve university: ", resultUniversity.Error())
		var httpErr types.HttpError

		switch resultUniversity.Error() {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"University not found",
				"University with ID "+id+" not found",
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to retrieve university",
				resultUniversity.Error().Error(),
				"University ID: "+id,
			)
		}

		return types.ResultErr[models.UniversityDB](&httpErr)
	}

	deletedUniversity := resultUniversity.Value().(models.UniversityDB)
	err = configs.DB.SoftDeleteOne(filter, deletedUniversity)
	if err != nil {
		logger.Warning("Failed to delete university in database: ", err)
		var httpErr types.HttpError

		switch err {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"University not found",
				"University with ID "+id+" not found",
			)

		case configs.DBErr.NotModified():
			httpErr = types.Error(
				types.Http.C400().PreconditionFailed(),
				"University was already marked as deleted",
				"University ID: "+id,
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to delete university",
				err.Error(),
				"University ID: "+id,
			)
		}

		return types.ResultErr[models.UniversityDB](&httpErr)
	}

	return types.ResultOk(deletedUniversity)
}
func (universityType) DeletePermanentByID(id string, filter models.FilterObject) types.Result[models.UniversityDB] {
	logger.Warning("Permanently deleting university by ID: ", id)
	oid, err := models.ID.ToDB(id)
	if err != nil {
		logger.Warning("Failed to convert ID to ObjectID: ", err)
		httpErr := types.Error(
			types.Http.C400().UnprocessableEntity(),
			"Invalid value",
			"Invalid ID format: "+err.Error(),
			"University ID: "+id,
		)
		return types.ResultErr[models.UniversityDB](&httpErr)
	}

	filter = models.Filter.AddPart(filter, models.Filter.ID(oid))
	filter = models.Filter.AddPart(filter, models.Filter.Deleted())
	resultUniversity := configs.DB.FindOne(filter, models.UniversityDB{})
	if resultUniversity.IsErr() {
		logger.Warning("Failed to find university for permanent deletion: ", resultUniversity.Error())
		var httpErr types.HttpError

		switch resultUniversity.Error() {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"University not found",
				"University with ID "+id+" not found or not marked as deleted",
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to find university for permanent deletion",
				resultUniversity.Error().Error(),
				"University ID: "+id,
			)
		}

		return types.ResultErr[models.UniversityDB](&httpErr)
	}

	err = configs.DB.PermanentDeleteOne(filter, models.UniversityDB{})
	if err != nil {
		logger.Warning("Failed to permanently delete university in database: ", err)
		var httpErr types.HttpError

		switch err {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"University not found",
				"University with ID "+id+" not found or not marked as deleted",
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to permanently delete university",
				err.Error(),
				"University ID: "+id,
			)
		}

		return types.ResultErr[models.UniversityDB](&httpErr)
	}

	return types.ResultOk(resultUniversity.Value().(models.UniversityDB))
}
func (universityType) DeletePermanentAll(filter models.FilterObject) types.Result[[]models.UniversityDB] {
	filter = models.Filter.AddPart(filter, models.Filter.Deleted())

	resultUniversitys := configs.DB.FindAll(filter, models.UniversityDB{})
	if resultUniversitys.IsErr() {
		logger.Warning("Failed to find university for permanent deletion: ", resultUniversitys.Error())
		var httpErr types.HttpError

		switch resultUniversitys.Error() {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"Universitys not found",
				"No universitys marked as deleted found",
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to find university for permanent deletion",
				resultUniversitys.Error().Error(),
			)
		}

		return types.ResultErr[[]models.UniversityDB](&httpErr)
	}

	err := configs.DB.PermanentDeleteMany(filter, models.UniversityDB{})
	if err != nil {
		logger.Warning("Failed to permanently delete all universitys in database: ", err)
		var httpErr types.HttpError

		switch err {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"Universitys not found",
				"No universitys marked as deleted found",
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to permanently delete all universitys",
				err.Error(),
			)
		}
		return types.ResultErr[[]models.UniversityDB](&httpErr)
	}

	universitys := utils.Map(resultUniversitys.Value(), models.InterfaceTo[models.UniversityDB])
	return types.ResultOk(universitys)
}
