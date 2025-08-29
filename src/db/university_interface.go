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
