package db

import (
	"dainxor/atv/configs"
	"dainxor/atv/logger"
	"dainxor/atv/models"
	"dainxor/atv/types"
	"dainxor/atv/utils"

	"go.mongodb.org/mongo-driver/v2/bson"
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

func (specialityType) GetByID(id string) types.Result[models.SpecialityDB] {
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

	filter := bson.D{{Key: "_id", Value: oid}}

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
func (specialityType) GetAll() types.Result[[]models.SpecialityDB] {
	filter := bson.D{models.Filter.NotDeleted()} // Filter to exclude deleted specialities

	resultSpecialities := configs.DB.FindAll(filter, models.SpecialityDB{})
	if resultSpecialities.IsErr() {
		logger.Warning("Failed to get all specialities from database:", resultSpecialities.Error())
		return types.ResultErr[[]models.SpecialityDB](resultSpecialities.Error())
	}

	specialities := utils.Map(resultSpecialities.Value(), models.InterfaceTo[models.SpecialityDB])
	logger.Debug("Retrieved", len(specialities), "specialities from database")
	return types.ResultOk(specialities)
}
