package dao

import (
	"dainxor/atv/configs"
	"dainxor/atv/logger"
	"dainxor/atv/models"
	"dainxor/atv/types"
	"dainxor/atv/utils"
)

type companionType struct{}

var Companion companionType

func (companionType) Create(companion models.CompanionCreate) types.Result[models.CompanionDB] {
	resultCompanionDB := companion.ToInsert()
	if resultCompanionDB.IsErr() {
		logger.Warning("Error converting companion to DB model:", resultCompanionDB.Error())
		return types.ResultErr[models.CompanionDB](resultCompanionDB.Error())
	}

	companionDB := resultCompanionDB.Value()
	resultCreate := configs.DB.InsertOne(companionDB)
	if resultCreate.IsErr() {
		logger.Warning("Failed to create companion in MongoDB:", resultCreate.Error())
		return types.ResultErr[models.CompanionDB](resultCreate.Error())
	}

	resultGet := configs.DB.FindOne(models.Filter.ID(resultCreate.Value()), &companionDB)

	if resultGet.IsErr() {
		logger.Warning("Failed to get created companion in MongoDB:", resultGet.Error())
		return types.ResultErr[models.CompanionDB](resultGet.Error())
	}

	companionDB.ID = resultCreate.Value()
	return types.ResultOk(companionDB)
}

func (companionType) GetByID(id string, filter models.FilterObject) types.Result[models.CompanionDB] {
	oid, err := models.ID.ToBson(id)

	if err != nil {
		logger.Warning("Failed to convert ID to ObjectID: ", err)
		httpErr := types.Error(
			types.Http.C400().UnprocessableEntity(),
			"Invalid value",
			"Invalid ID format: "+err.Error(),
			"Companion ID: "+id,
		)
		return types.ResultErr[models.CompanionDB](&httpErr)
	}

	filter = models.Filter.AddPart(filter, models.Filter.ID(oid))
	filter = models.Filter.AddPart(filter, models.Filter.NotDeleted())
	var companion models.CompanionDB

	resultGet := configs.DB.FindOne(filter, companion)
	if resultGet.IsErr() {
		logger.Warning("Failed to get companion by ID: ", resultGet.Error())

		return types.ResultErr[models.CompanionDB](resultGet.Error())
	}
	companion = resultGet.Value().(models.CompanionDB)
	return types.ResultOk(companion)
}
func (companionType) GetByNumberID(numberID string, filter models.FilterObject) types.Result[models.CompanionDB] {
	filter = models.Filter.Add(filter, "number_id", []string{numberID}) // This is not a DBID
	filter = models.Filter.AddPart(filter, models.Filter.NotDeleted())

	resultCompanionDB := configs.DB.FindOne(filter, models.CompanionDB{})
	if resultCompanionDB.IsErr() {
		logger.Warning("Failed to get companion by number id:", resultCompanionDB.Error())

		return types.ResultErr[models.CompanionDB](resultCompanionDB.Error())
	}

	return types.ResultOk(resultCompanionDB.Value().(models.CompanionDB))
}
func (companionType) GetByEmail(email string, filter models.FilterObject) types.Result[models.CompanionDB] {
	filter = models.Filter.Add(filter, "email", []string{email})
	filter = models.Filter.AddPart(filter, models.Filter.NotDeleted())
	var companion models.CompanionDB

	resultGet := configs.DB.FindOne(filter, &companion)
	if resultGet.IsErr() {
		logger.Warning("Failed to get companion by email:", resultGet.Error())

		return types.ResultErr[models.CompanionDB](resultGet.Error())

	}

	companion = resultGet.Value().(models.CompanionDB)
	return types.ResultOk(companion)
}
func (companionType) GetAll(filter models.FilterObject) types.Result[[]models.CompanionDB] {
	filter = models.Filter.AddPart(filter, models.Filter.NotDeleted())

	resultCompanionsDB := configs.DB.FindAll(filter, models.CompanionDB{})
	if resultCompanionsDB.IsErr() {
		logger.Error("Failed to get all companions from MongoDB:", resultCompanionsDB.Error())

		return types.ResultErr[[]models.CompanionDB](resultCompanionsDB.Error())
	}

	companionsDB := utils.Map(resultCompanionsDB.Value(), models.InterfaceTo[models.CompanionDB])
	logger.Debug("Retrieved", len(companionsDB), "companions from MongoDB database")
	return types.ResultOk(companionsDB)
}

func (companionType) UpdateByID(id string, companion models.CompanionCreate, filter models.FilterObject) types.Result[models.CompanionDB] {
	oid, err := models.ID.ToDB(id)
	if err != nil {
		logger.Error("Failed to convert ID to ObjectID: ", err)
		httpErr := types.Error(
			types.Http.C400().UnprocessableEntity(),
			"Invalid value",
			"Invalid ID format: "+err.Error(),
			"Companion ID: "+id,
		)
		return types.ResultErr[models.CompanionDB](&httpErr)
	}

	resultCompanionDB := companion.ToUpdate()
	if resultCompanionDB.IsErr() {
		logger.Error("Error converting companion to DB model")
		httpErr := types.Error(
			types.Http.C400().UnprocessableEntity(),
			"Invalid value",
			resultCompanionDB.Error().Error(),
			"Companion ID: "+id,
		)
		return types.ResultErr[models.CompanionDB](&httpErr)
	}

	filter = models.Filter.AddPart(filter, models.Filter.ID(oid))
	filter = models.Filter.AddPart(filter, models.Filter.NotDeleted())
	err = configs.DB.UpdateOne(filter, resultCompanionDB.Value())

	if err != nil {
		logger.Error("Failed to update companion in MongoDB: ", err)
		httpErr := types.ErrorInternal(
			"Failed to update companion",
			err.Error(),
			"Companion ID: "+id,
		)
		return types.ResultErr[models.CompanionDB](&httpErr)
	}

	return Companion.GetByID(id, models.Filter.Empty())
}

func (companionType) PatchByID(id string, companion models.CompanionCreate, filter models.FilterObject) types.Result[models.CompanionDB] {
	oid, err := models.ID.ToDB(id)
	if err != nil {
		logger.Error("Failed to convert ID to ObjectID: ", err)
		httpErr := types.Error(
			types.Http.C400().UnprocessableEntity(),
			"Invalid value",
			"Invalid ID format: "+err.Error(),
			"Companion ID: "+id,
		)
		return types.ResultErr[models.CompanionDB](&httpErr)
	}

	resultCompanionDB := companion.ToUpdate()
	if resultCompanionDB.IsErr() {
		logger.Error("Error converting companion to DB model")
		httpErr := types.Error(
			types.Http.C400().UnprocessableEntity(),
			"Invalid value",
			resultCompanionDB.Error().Error(),
			"Companion ID: "+id,
		)
		return types.ResultErr[models.CompanionDB](&httpErr)
	}

	filter = models.Filter.AddPart(filter, models.Filter.ID(oid))
	filter = models.Filter.AddPart(filter, models.Filter.NotDeleted())
	err = configs.DB.PatchOne(filter, resultCompanionDB.Value())

	if err != nil {
		logger.Error("Failed to update companion in MongoDB: ", err)
		httpErr := types.ErrorInternal(
			"Failed to update companion",
			err.Error(),
			"Companion ID: "+id,
		)
		return types.ResultErr[models.CompanionDB](&httpErr)
	}

	return Companion.GetByID(id, models.Filter.Empty())
}

func (companionType) DeleteByID(id string, filter models.FilterObject) types.Result[models.CompanionDB] {
	oid, err := models.ID.ToBson(id)
	if err != nil {
		logger.Error("Failed to convert ID to ObjectID: ", err)
		httpErr := types.Error(
			types.Http.C400().UnprocessableEntity(),
			"Invalid value",
			"Invalid ID format: "+err.Error(),
			"Companion ID: "+id,
		)
		return types.ResultErr[models.CompanionDB](&httpErr)
	}

	filter = models.Filter.AddPart(filter, models.Filter.ID(oid))
	filter = models.Filter.AddPart(filter, models.Filter.NotDeleted())
	var companion models.CompanionDB
	companionResult := configs.DB.FindOne(filter, companion)

	if companionResult.IsErr() {
		logger.Error("Failed to find companion for soft delete: ", companionResult.Error())
		httpErr := types.ErrorInternal(
			"Failed to find companion for soft delete",
			companionResult.Error().Error(),
			"Companion ID: "+id,
		)
		return types.ResultErr[models.CompanionDB](&httpErr)
	}

	companion = companionResult.Value().(models.CompanionDB)
	err = configs.DB.UpdateOne(filter, companion)

	if err != nil {
		logger.Error("Failed to delete companion in MongoDB: ", err)
		httpErr := types.ErrorInternal(
			"Failed to delete companion",
			err.Error(),
			"Companion ID: "+id,
		)
		return types.ResultErr[models.CompanionDB](&httpErr)
	}

	return types.ResultOk(companion)
}
func (companionType) DeletePermanentByID(id string, filter models.FilterObject) types.Result[models.CompanionDB] {
	logger.Warning("Permanently deleting companion by ID: ", id)
	oid, err := models.ID.ToBson(id)
	if err != nil {
		logger.Error("Failed to convert ID to ObjectID: ", err)
		httpErr := types.Error(
			types.Http.C400().UnprocessableEntity(),
			"Invalid value",
			"Invalid ID format: "+err.Error(),
			"Companion ID: "+id,
		)
		return types.ResultErr[models.CompanionDB](&httpErr)
	}

	filter = models.Filter.AddPart(filter, models.Filter.ID(oid))
	filter = models.Filter.AddPart(filter, models.Filter.NotDeleted())
	var companion models.CompanionDB
	companionResult := configs.DB.FindOne(filter, companion)

	if companionResult.IsErr() {
		logger.Debug("Failed to find companion for permanent deletion: ", companionResult.Error())

		if companionResult.Error() == configs.DBErr.NotFound() {
			httpErr := types.ErrorNotFound(
				"Companion not found",
				"Companion with ID "+id+" not found or not marked as deleted",
			)
			return types.ResultErr[models.CompanionDB](&httpErr)
		}

		httpErr := types.ErrorInternal(
			"Failed to find companion for permanent deletion",
			companionResult.Error().Error(),
			"Companion ID: "+id,
		)
		return types.ResultErr[models.CompanionDB](&httpErr)
	}

	companion = companionResult.Value().(models.CompanionDB)
	err = configs.DB.PermanentDeleteOne(filter, companion)
	if err != nil {
		logger.Debug("Failed to permanently delete companion in MongoDB: ", err)
		httpErr := types.ErrorInternal(
			"Failed to permanently delete companion",
			err.Error(),
			"Companion ID: "+id,
		)
		return types.ResultErr[models.CompanionDB](&httpErr)
	}

	return types.ResultOk(companion)
}
func (companionType) DeletePermanentAll(filter models.FilterObject) types.Result[[]models.CompanionDB] {
	filter = models.Filter.AddPart(filter, models.Filter.Deleted())

	var companion models.CompanionDB
	companionResult := configs.DB.FindAll(filter, companion)
	if companionResult.IsErr() {
		logger.Error("Failed to find companions for permanent deletion: ", companionResult.Error())
		httpErr := types.ErrorInternal(
			"Failed to find companions for permanent deletion",
			companionResult.Error().Error(),
			"Companion filter: "+filter.String(),
		)
		return types.ResultErr[[]models.CompanionDB](&httpErr)
	}

	err := configs.DB.PermanentDeleteMany(filter, companion)
	if err != nil {
		logger.Error("Failed to permanently delete all companions in MongoDB: ", err)
		httpErr := types.ErrorInternal(
			"Failed to permanently delete all companions",
			err.Error(),
		)
		return types.ResultErr[[]models.CompanionDB](&httpErr)
	}

	companions := utils.Map(companionResult.Value(), models.InterfaceTo[models.CompanionDB])
	return types.ResultOk(companions)
}
