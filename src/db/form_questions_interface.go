package db

import (
	"dainxor/atv/configs"
	"dainxor/atv/logger"
	"dainxor/atv/models"
	"dainxor/atv/types"
)

type formQuestionsNS struct{}

var FormQuestions formQuestionsNS

func (formQuestionsNS) Create(t models.FormQuestionCreate) types.Result[models.FormQuestionDB] {
	resultQuestionDB := t.ToInsert()
	if resultQuestionDB.IsErr() {
		logger.Warning("Error converting form question to DB model:", resultQuestionDB.Error())
		return types.ResultErr[models.FormQuestionDB](resultQuestionDB.Error())
	}

	questionDB := resultQuestionDB.Value()
	resultCreate := configs.DB.InsertOne(questionDB)

	if resultCreate.IsErr() {
		logger.Warning("Failed to create form question in MongoDB:", resultCreate.Error())
		return types.ResultErr[models.FormQuestionDB](resultCreate.Error())
	}

	resultGet := configs.DB.FindOne(models.Filter.ID(resultCreate.Value()), &questionDB)

	if resultGet.IsErr() {
		logger.Warning("Failed to get created form question in MongoDB:", resultGet.Error())
		return types.ResultErr[models.FormQuestionDB](resultGet.Error())
	}

	questionDB.ID = resultCreate.Value()
	return types.ResultOk(questionDB)
}
func (formQuestionsNS) GetByID(id string, filter models.FilterObject) types.Result[models.FormQuestionDB] {
	oid, err := models.ID.ToBson(id)

	if err != nil {
		logger.Warning("Failed to convert ID to ObjectID: ", err)
		httpErr := types.Error(
			types.Http.C400().UnprocessableEntity(),
			"Invalid value",
			"Invalid ID format: "+err.Error(),
			"Form Question Type ID: "+id,
		)
		return types.ResultErr[models.FormQuestionDB](&httpErr)
	}

	filter = models.Filter.AddPart(filter, models.Filter.ID(oid))
	filter = models.Filter.AddPart(filter, models.Filter.NotDeleted())
	var questionType models.FormQuestionDB

	resultGet := configs.DB.FindOne(filter, questionType)
	if resultGet.IsErr() {
		logger.Warning("Failed to get form question type by ID: ", resultGet.Error())

		return types.ResultErr[models.FormQuestionDB](resultGet.Error())
	}
	questionType = resultGet.Value().(models.FormQuestionDB)
	return types.ResultOk(questionType)
}
