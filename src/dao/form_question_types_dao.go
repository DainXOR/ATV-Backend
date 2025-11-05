package dao

import (
	"dainxor/atv/configs"
	"dainxor/atv/logger"
	"dainxor/atv/models"
	"dainxor/atv/types"
	"dainxor/atv/utils"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type formQuestionTypesNS struct{}

var FormQuestionTypes formQuestionTypesNS

func (formQuestionTypesNS) Create(t models.FormQuestionTypeCreate) types.Result[models.FormQuestionTypeDB] {
	typeDB := t.ToInsert()
	resultCreate := configs.DB.InsertOne(typeDB)

	if resultCreate.IsErr() {
		logger.Warning("Failed to create companion in MongoDB:", resultCreate.Error())
		return types.ResultErr[models.FormQuestionTypeDB](resultCreate.Error())
	}

	resultGet := configs.DB.FindOne(bson.D{{Key: "_id", Value: resultCreate.Value()}}, &typeDB)

	if resultGet.IsErr() {
		logger.Warning("Failed to get created companion in MongoDB:", resultGet.Error())
		return types.ResultErr[models.FormQuestionTypeDB](resultGet.Error())
	}

	typeDB.ID = resultCreate.Value()
	return types.ResultOk(typeDB)
}

func (formQuestionTypesNS) GetByID(id string, filter models.FilterObject) types.Result[models.FormQuestionTypeDB] {
	oid, err := models.ID.ToBson(id)

	if err != nil {
		logger.Warning("Failed to convert ID to ObjectID: ", err)
		httpErr := types.Error(
			types.Http.C400().UnprocessableEntity(),
			"Invalid value",
			"Invalid ID format: "+err.Error(),
			"Form Question Type ID: "+id,
		)
		return types.ResultErr[models.FormQuestionTypeDB](&httpErr)
	}

	filter = models.Filter.AddPart(filter, models.Filter.ID(oid))
	filter = models.Filter.AddPart(filter, models.Filter.NotDeleted())
	var questionType models.FormQuestionTypeDB

	resultGet := configs.DB.FindOne(filter, questionType)
	if resultGet.IsErr() {
		logger.Warning("Failed to get form question type by ID: ", resultGet.Error())

		return types.ResultErr[models.FormQuestionTypeDB](resultGet.Error())
	}
	questionType = resultGet.Value().(models.FormQuestionTypeDB)
	return types.ResultOk(questionType)
}
func (formQuestionTypesNS) GetAll(filter models.FilterObject) types.Result[[]models.FormQuestionTypeDB] {
	filter = models.Filter.AddPart(filter, models.Filter.NotDeleted())

	resultTypes := configs.DB.FindAll(filter, models.FormQuestionTypeDB{})
	if resultTypes.IsErr() {
		logger.Warning("Failed to get all form question types: ", resultTypes.Error())
		return types.ResultErr[[]models.FormQuestionTypeDB](resultTypes.Error())
	}

	typesDB := utils.Map(resultTypes.Value(), models.InterfaceTo[models.FormQuestionTypeDB])
	logger.Debug("Retrieved", len(typesDB), "form question types from db")
	return types.ResultOk(typesDB)
}

func (formQuestionTypesNS) UpdateByID(id string, model models.FormQuestionTypeCreate, filter models.FilterObject) types.Result[models.FormQuestionTypeDB] {
	oid, err := models.ID.ToDB(id)
	if err != nil {
		logger.Warning("Failed to convert ID to ObjectID: ", err)
		httpErr := types.Error(
			types.Http.C400().UnprocessableEntity(),
			"Invalid value",
			"Invalid ID format: "+err.Error(),
			"Question type ID: "+id,
		)
		return types.ResultErr[models.FormQuestionTypeDB](&httpErr)
	}

	objectUpdate := model.ToUpdate()

	filter = models.Filter.AddPart(filter, models.Filter.ID(oid))
	filter = models.Filter.AddPart(filter, models.Filter.NotDeleted())
	err = configs.DB.UpdateOne(filter, objectUpdate)
	if err != nil {
		logger.Warning("Failed to update question type in database: ", err)
		var httpErr types.HttpError

		switch err {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"Question type not found",
				"Question type with ID "+id+" not found",
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to update question type",
				err.Error(),
				"Question type ID: "+id,
			)
		}
		return types.ResultErr[models.FormQuestionTypeDB](&httpErr)
	}

	resultObjectDB := configs.DB.FindOne(filter, models.FormQuestionTypeDB{})
	if resultObjectDB.IsErr() {
		logger.Warning("Failed to retrieve updated question type: ", resultObjectDB.Error())
		var httpErr types.HttpError

		switch resultObjectDB.Error() {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"Question type not found",
				"Question type with ID "+id+" not found",
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to retrieve updated question type",
				resultObjectDB.Error().Error(),
				"Question type ID: "+id,
			)
		}
		return types.ResultErr[models.FormQuestionTypeDB](&httpErr)
	}

	return types.ResultOk(resultObjectDB.Value().(models.FormQuestionTypeDB))
}

func (formQuestionTypesNS) PatchByID(id string, model models.FormQuestionTypeCreate, filter models.FilterObject) types.Result[models.FormQuestionTypeDB] {
	oid, err := models.ID.ToDB(id)
	if err != nil {
		logger.Warning("Failed to convert ID to ObjectID: ", err)
		httpErr := types.Error(
			types.Http.C400().UnprocessableEntity(),
			"Invalid value",
			"Invalid ID format: "+err.Error(),
			"Question type ID: "+id,
		)
		return types.ResultErr[models.FormQuestionTypeDB](&httpErr)
	}

	objectPatch := model.ToUpdate()

	filter = models.Filter.AddPart(filter, models.Filter.ID(oid))
	filter = models.Filter.AddPart(filter, models.Filter.NotDeleted())
	err = configs.DB.PatchOne(filter, objectPatch)
	if err != nil {
		logger.Warning("Failed to patch question type in database: ", err)
		var httpErr types.HttpError

		switch err {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"Question type not found",
				"Question type with ID "+id+" not found",
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to patch question type",
				err.Error(),
				"Question type ID: "+id,
			)
		}
		return types.ResultErr[models.FormQuestionTypeDB](&httpErr)
	}

	resultObjectDB := configs.DB.FindOne(filter, models.FormQuestionTypeDB{})
	if resultObjectDB.IsErr() {
		logger.Warning("Failed to retrieve updated question type: ", resultObjectDB.Error())
		var httpErr types.HttpError

		switch resultObjectDB.Error() {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"Question type not found",
				"Question type with ID "+id+" not found",
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to retrieve updated question type",
				resultObjectDB.Error().Error(),
				"Question type ID: "+id,
			)
		}
		return types.ResultErr[models.FormQuestionTypeDB](&httpErr)
	}

	return types.ResultOk(resultObjectDB.Value().(models.FormQuestionTypeDB))
}

func (formQuestionTypesNS) DeleteByID(id string, filter models.FilterObject) types.Result[models.FormQuestionTypeDB] {
	oid, err := models.ID.ToDB(id)
	if err != nil {
		logger.Warning("Failed to convert ID to ObjectID: ", err)
		httpErr := types.Error(
			types.Http.C400().UnprocessableEntity(),
			"Invalid value",
			"Invalid ID format: "+err.Error(),
			"Question type ID: "+id,
		)
		return types.ResultErr[models.FormQuestionTypeDB](&httpErr)
	}

	filter = models.Filter.AddPart(filter, models.Filter.ID(oid))
	filter = models.Filter.AddPart(filter, models.Filter.NotDeleted())
	resultObject := configs.DB.FindOne(filter, models.FormQuestionTypeDB{})
	if resultObject.IsErr() {
		logger.Warning("Failed to retrieve question type: ", resultObject.Error())
		var httpErr types.HttpError

		switch resultObject.Error() {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"Question type not found",
				"Question type with ID "+id+" not found",
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to retrieve question type",
				resultObject.Error().Error(),
				"Question type ID: "+id,
			)
		}

		return types.ResultErr[models.FormQuestionTypeDB](&httpErr)
	}

	deletedObject := resultObject.Value().(models.FormQuestionTypeDB)
	err = configs.DB.SoftDeleteOne(filter, deletedObject)
	if err != nil {
		logger.Warning("Failed to delete question type in database: ", err)
		var httpErr types.HttpError

		switch err {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"Question type not found",
				"Question type with ID "+id+" not found",
			)

		case configs.DBErr.NotModified():
			httpErr = types.Error(
				types.Http.C400().PreconditionFailed(),
				"Question type was already marked as deleted",
				"Question type ID: "+id,
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to delete question type",
				err.Error(),
				"Question type ID: "+id,
			)
		}

		return types.ResultErr[models.FormQuestionTypeDB](&httpErr)
	}

	return types.ResultOk(deletedObject)
}
func (formQuestionTypesNS) DeletePermanentByID(id string, filter models.FilterObject) types.Result[models.FormQuestionTypeDB] {
	logger.Warning("Permanently deleting question type by ID: ", id)
	oid, err := models.ID.ToDB(id)
	if err != nil {
		logger.Warning("Failed to convert ID to ObjectID: ", err)
		httpErr := types.Error(
			types.Http.C400().UnprocessableEntity(),
			"Invalid value",
			"Invalid ID format: "+err.Error(),
			"Question type ID: "+id,
		)
		return types.ResultErr[models.FormQuestionTypeDB](&httpErr)
	}

	filter = models.Filter.AddPart(filter, models.Filter.ID(oid))
	filter = models.Filter.AddPart(filter, models.Filter.Deleted())
	resultObject := configs.DB.FindOne(filter, models.FormQuestionTypeDB{})
	if resultObject.IsErr() {
		logger.Warning("Failed to find question type for permanent deletion: ", resultObject.Error())
		var httpErr types.HttpError

		switch resultObject.Error() {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"Question type not found",
				"Question type with ID "+id+" not found or not marked as deleted",
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to find question type for permanent deletion",
				resultObject.Error().Error(),
				"Question type ID: "+id,
			)
		}

		return types.ResultErr[models.FormQuestionTypeDB](&httpErr)
	}

	err = configs.DB.PermanentDeleteOne(filter, models.FormQuestionTypeDB{})
	if err != nil {
		logger.Warning("Failed to permanently delete question type in database: ", err)
		var httpErr types.HttpError

		switch err {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"Question type not found",
				"Question type with ID "+id+" not found or not marked as deleted",
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to permanently delete question type",
				err.Error(),
				"Question type ID: "+id,
			)
		}

		return types.ResultErr[models.FormQuestionTypeDB](&httpErr)
	}

	return types.ResultOk(resultObject.Value().(models.FormQuestionTypeDB))
}
func (formQuestionTypesNS) DeletePermanentAll(filter models.FilterObject) types.Result[[]models.FormQuestionTypeDB] {
	filter = models.Filter.AddPart(filter, models.Filter.Deleted())

	resultObjects := configs.DB.FindAll(filter, models.FormQuestionTypeDB{})
	if resultObjects.IsErr() {
		logger.Warning("Failed to find question type for permanent deletion: ", resultObjects.Error())
		var httpErr types.HttpError

		switch resultObjects.Error() {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"Question types not found",
				"No question types marked as deleted found",
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to find question type for permanent deletion",
				resultObjects.Error().Error(),
			)
		}

		return types.ResultErr[[]models.FormQuestionTypeDB](&httpErr)
	}

	err := configs.DB.PermanentDeleteMany(filter, models.FormQuestionTypeDB{})
	if err != nil {
		logger.Warning("Failed to permanently delete all question types in database: ", err)
		var httpErr types.HttpError

		switch err {
		case configs.DBErr.NotFound():
			httpErr = types.ErrorNotFound(
				"Question types not found",
				"No question types marked as deleted found",
			)

		default:
			httpErr = types.ErrorInternal(
				"Failed to permanently delete all question types",
				err.Error(),
			)
		}
		return types.ResultErr[[]models.FormQuestionTypeDB](&httpErr)
	}

	objects := utils.Map(resultObjects.Value(), models.InterfaceTo[models.FormQuestionTypeDB])
	return types.ResultOk(objects)
}
