package dao

import (
	"dainxor/atv/configs"
	"dainxor/atv/logger"
	"dainxor/atv/models"
	"dainxor/atv/types"
	"dainxor/atv/utils"
)

type specialityType struct{}

var Speciality specialityType

func (specialityType) Create(u models.SpecialityCreate) types.Result[models.SpecialityDB] {
	specialityDB := u.ToInsert()
	resultID := configs.DB.InsertOne(specialityDB)

	if resultID.IsErr() {
		logger.Error("Error inserting speciality: ", resultID.Error())
		return types.ResultErr[models.SpecialityDB](resultID.Error())
	}

	specialityDB.ID = resultID.Value()
	return types.ResultOk(specialityDB)
}

func (specialityType) GetByID(id string, filter models.FilterObject) types.Result[models.SpecialityDB] {
	oid, err := models.ID.ToDB(id)

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

	filter = models.Filter.AddPart(filter, models.Filter.ID(oid))
	filter = models.Filter.AddPart(filter, models.Filter.NotDeleted())

	resultSpeciality := configs.DB.FindOne(filter, models.SpecialityDB{})
	if resultSpeciality.IsErr() {
		logger.Warning("Failed to get speciality by ID: ", resultSpeciality.Error())
		var httpErr types.HttpError

		switch resultSpeciality.Error() {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"Speciality not found",
				"Speciality with ID "+id+" not found",
			)

		default:
			logger.Error("Failed to get speciality by ID: ", resultSpeciality.Error())
			httpErr = types.ErrorInternal(
				"Failed to retrieve speciality",
				resultSpeciality.Error().Error(),
				"Speciality ID: "+id,
			)
		}

		return types.ResultErr[models.SpecialityDB](&httpErr)
	}

	return types.ResultOk(resultSpeciality.Value().(models.SpecialityDB))
}
func (specialityType) GetAll(filter models.FilterObject) types.Result[[]models.SpecialityDB] {
	filter = models.Filter.AddPart(filter, models.Filter.NotDeleted())

	resultSpecialities := configs.DB.FindAll(filter, models.SpecialityDB{})
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
		return types.ResultErr[[]models.SpecialityDB](&httpErr)
	}

	specialities := utils.Map(resultSpecialities.Value(), models.InterfaceTo[models.SpecialityDB])
	logger.Debug("Retrieved", len(specialities), "specialities from database")
	return types.ResultOk(specialities)
}

func (specialityType) UpdateByID(id string, speciality models.SpecialityCreate, filter models.FilterObject) types.Result[models.SpecialityDB] {
	oid, err := models.ID.ToDB(id)
	if err != nil {
		logger.Warning("Failed to convert ID to ObjectID: ", err)
		httpErr := types.Error(
			types.Http.C400().UnprocessableEntity(),
			"Invalid value",
			"Invalid ID format: "+err.Error(),
			"Speciality ID: "+id,
		)
		return types.ResultErr[models.SpecialityDB](&httpErr)
	}

	specialityUpdate := speciality.ToUpdate()

	filter = models.Filter.AddPart(filter, models.Filter.ID(oid))
	filter = models.Filter.AddPart(filter, models.Filter.NotDeleted())
	err = configs.DB.UpdateOne(filter, specialityUpdate)
	if err != nil {
		logger.Warning("Failed to update speciality in database: ", err)
		var httpErr types.HttpError

		switch err {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"Speciality not found",
				"Speciality with ID "+id+" not found",
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to update speciality",
				err.Error(),
				"Speciality ID: "+id,
			)
		}
		return types.ResultErr[models.SpecialityDB](&httpErr)
	}

	resultSpecialityDB := configs.DB.FindOne(filter, models.SpecialityDB{})
	if resultSpecialityDB.IsErr() {
		logger.Warning("Failed to retrieve updated speciality: ", resultSpecialityDB.Error())
		var httpErr types.HttpError

		switch resultSpecialityDB.Error() {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"Speciality not found",
				"Speciality with ID "+id+" not found",
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to retrieve updated speciality",
				resultSpecialityDB.Error().Error(),
				"Speciality ID: "+id,
			)
		}
		return types.ResultErr[models.SpecialityDB](&httpErr)
	}

	return types.ResultOk(resultSpecialityDB.Value().(models.SpecialityDB))
}

func (specialityType) PatchByID(id string, speciality models.SpecialityCreate, filter models.FilterObject) types.Result[models.SpecialityDB] {
	oid, err := models.ID.ToDB(id)
	if err != nil {
		logger.Warning("Failed to convert ID to ObjectID: ", err)
		httpErr := types.Error(
			types.Http.C400().UnprocessableEntity(),
			"Invalid value",
			"Invalid ID format: "+err.Error(),
			"Speciality ID: "+id,
		)
		return types.ResultErr[models.SpecialityDB](&httpErr)
	}

	specialityPatch := speciality.ToUpdate()

	filter = models.Filter.AddPart(filter, models.Filter.ID(oid))
	filter = models.Filter.AddPart(filter, models.Filter.NotDeleted())
	err = configs.DB.PatchOne(filter, specialityPatch)
	if err != nil {
		logger.Warning("Failed to patch speciality in database: ", err)
		var httpErr types.HttpError

		switch err {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"Speciality not found",
				"Speciality with ID "+id+" not found",
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to patch speciality",
				err.Error(),
				"Speciality ID: "+id,
			)
		}
		return types.ResultErr[models.SpecialityDB](&httpErr)
	}

	resultSpecialityDB := configs.DB.FindOne(filter, models.SpecialityDB{})
	if resultSpecialityDB.IsErr() {
		logger.Warning("Failed to retrieve updated speciality: ", resultSpecialityDB.Error())
		var httpErr types.HttpError

		switch resultSpecialityDB.Error() {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"Speciality not found",
				"Speciality with ID "+id+" not found",
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to retrieve updated speciality",
				resultSpecialityDB.Error().Error(),
				"Speciality ID: "+id,
			)
		}
		return types.ResultErr[models.SpecialityDB](&httpErr)
	}

	return types.ResultOk(resultSpecialityDB.Value().(models.SpecialityDB))
}

func (specialityType) DeleteByID(id string, filter models.FilterObject) types.Result[models.SpecialityDB] {
	oid, err := models.ID.ToDB(id)
	if err != nil {
		logger.Warning("Failed to convert ID to ObjectID: ", err)
		httpErr := types.Error(
			types.Http.C400().UnprocessableEntity(),
			"Invalid value",
			"Invalid ID format: "+err.Error(),
			"Speciality ID: "+id,
		)
		return types.ResultErr[models.SpecialityDB](&httpErr)
	}

	filter = models.Filter.AddPart(filter, models.Filter.ID(oid))
	filter = models.Filter.AddPart(filter, models.Filter.NotDeleted())
	resultSpeciality := configs.DB.FindOne(filter, models.SpecialityDB{})
	if resultSpeciality.IsErr() {
		logger.Warning("Failed to retrieve speciality: ", resultSpeciality.Error())
		var httpErr types.HttpError

		switch resultSpeciality.Error() {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"Speciality not found",
				"Speciality with ID "+id+" not found",
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to retrieve speciality",
				resultSpeciality.Error().Error(),
				"Speciality ID: "+id,
			)
		}

		return types.ResultErr[models.SpecialityDB](&httpErr)
	}

	deletedSpeciality := resultSpeciality.Value().(models.SpecialityDB)
	err = configs.DB.SoftDeleteOne(filter, deletedSpeciality)
	if err != nil {
		logger.Warning("Failed to delete speciality in database: ", err)
		var httpErr types.HttpError

		switch err {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"Speciality not found",
				"Speciality with ID "+id+" not found",
			)

		case configs.DBErr.NotModified():
			httpErr = types.Error(
				types.Http.C400().PreconditionFailed(),
				"Speciality was already marked as deleted",
				"Speciality ID: "+id,
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to delete speciality",
				err.Error(),
				"Speciality ID: "+id,
			)
		}

		return types.ResultErr[models.SpecialityDB](&httpErr)
	}

	return types.ResultOk(deletedSpeciality)
}
func (specialityType) DeletePermanentByID(id string, filter models.FilterObject) types.Result[models.SpecialityDB] {
	logger.Warning("Permanently deleting speciality by ID: ", id)
	oid, err := models.ID.ToDB(id)
	if err != nil {
		logger.Warning("Failed to convert ID to ObjectID: ", err)
		httpErr := types.Error(
			types.Http.C400().UnprocessableEntity(),
			"Invalid value",
			"Invalid ID format: "+err.Error(),
			"Speciality ID: "+id,
		)
		return types.ResultErr[models.SpecialityDB](&httpErr)
	}

	filter = models.Filter.AddPart(filter, models.Filter.ID(oid))
	filter = models.Filter.AddPart(filter, models.Filter.Deleted())
	resultSpeciality := configs.DB.FindOne(filter, models.SpecialityDB{})
	if resultSpeciality.IsErr() {
		logger.Warning("Failed to find speciality for permanent deletion: ", resultSpeciality.Error())
		var httpErr types.HttpError

		switch resultSpeciality.Error() {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"Speciality not found",
				"Speciality with ID "+id+" not found or not marked as deleted",
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to find speciality for permanent deletion",
				resultSpeciality.Error().Error(),
				"Speciality ID: "+id,
			)
		}

		return types.ResultErr[models.SpecialityDB](&httpErr)
	}

	err = configs.DB.PermanentDeleteOne(filter, models.SpecialityDB{})
	if err != nil {
		logger.Warning("Failed to permanently delete speciality in database: ", err)
		var httpErr types.HttpError

		switch err {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"Speciality not found",
				"Speciality with ID "+id+" not found or not marked as deleted",
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to permanently delete speciality",
				err.Error(),
				"Speciality ID: "+id,
			)
		}

		return types.ResultErr[models.SpecialityDB](&httpErr)
	}

	return types.ResultOk(resultSpeciality.Value().(models.SpecialityDB))
}
func (specialityType) DeletePermanentAll(filter models.FilterObject) types.Result[[]models.SpecialityDB] {
	filter = models.Filter.AddPart(filter, models.Filter.Deleted())

	resultSpecialitys := configs.DB.FindAll(filter, models.SpecialityDB{})
	if resultSpecialitys.IsErr() {
		logger.Warning("Failed to find speciality for permanent deletion: ", resultSpecialitys.Error())
		var httpErr types.HttpError

		switch resultSpecialitys.Error() {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"Specialitys not found",
				"No specialitys marked as deleted found",
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to find speciality for permanent deletion",
				resultSpecialitys.Error().Error(),
			)
		}

		return types.ResultErr[[]models.SpecialityDB](&httpErr)
	}

	err := configs.DB.PermanentDeleteMany(filter, models.SpecialityDB{})
	if err != nil {
		logger.Warning("Failed to permanently delete all specialitys in database: ", err)
		var httpErr types.HttpError

		switch err {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"Specialitys not found",
				"No specialitys marked as deleted found",
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to permanently delete all specialitys",
				err.Error(),
			)
		}
		return types.ResultErr[[]models.SpecialityDB](&httpErr)
	}

	specialitys := utils.Map(resultSpecialitys.Value(), models.InterfaceTo[models.SpecialityDB])
	return types.ResultOk(specialitys)
}
