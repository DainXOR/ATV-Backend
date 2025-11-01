package db

import (
	"dainxor/atv/configs"
	"dainxor/atv/logger"
	"dainxor/atv/models"
	"dainxor/atv/types"

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
