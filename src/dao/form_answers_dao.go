package dao

import (
	"dainxor/atv/configs"
	"dainxor/atv/logger"
	"dainxor/atv/models"
	"dainxor/atv/types"
	"dainxor/atv/utils"
)

type formAnswersNS struct{}

var FormAnswers formAnswersNS

func (formAnswersNS) Create(t models.FormAnswerCreate) types.Result[models.FormAnswerDB] {
	resultDB := t.ToInsert()
	if resultDB.IsErr() {
		logger.Warning("Error converting form answer to DB model:", resultDB.Error())
		return types.ResultErr[models.FormAnswerDB](resultDB.Error())
	}

	objectDB := resultDB.Value()
	resultCreate := configs.DB.InsertOne(objectDB)

	if resultCreate.IsErr() {
		logger.Warning("Failed to create form answer in MongoDB:", resultCreate.Error())
		return types.ResultErr[models.FormAnswerDB](resultCreate.Error())
	}

	objectDB.ID = resultCreate.Value()
	return types.ResultOk(objectDB)
}

func (formAnswersNS) GetByID(id string, filter models.FilterObject) types.Result[models.FormAnswerDB] {
	oid, err := models.ID.ToDB(id)

	if err != nil {
		logger.Warning("Failed to convert ID to ObjectID: ", err)
		httpErr := types.Error(
			types.Http.C400().UnprocessableEntity(),
			"Invalid value",
			"Invalid ID format: "+err.Error(),
			"Form Answer ID: "+id,
		)
		return types.ResultErr[models.FormAnswerDB](&httpErr)
	}

	filter = models.Filter.AddPart(filter, models.Filter.ID(oid))
	filter = models.Filter.AddPart(filter, models.Filter.NotDeleted())
	var answerType models.FormAnswerDB

	resultGet := configs.DB.FindOne(filter, answerType)
	if resultGet.IsErr() {
		logger.Warning("Failed to get form answer by ID: ", resultGet.Error())

		return types.ResultErr[models.FormAnswerDB](resultGet.Error())
	}
	answerType = resultGet.Value().(models.FormAnswerDB)
	return types.ResultOk(answerType)
}
func (formAnswersNS) GetAll(filter models.FilterObject) types.Result[[]models.FormAnswerDB] {
	filter = models.Filter.AddPart(filter, models.Filter.NotDeleted())

	resultObjects := configs.DB.FindAll(filter, models.FormAnswerDB{})
	if resultObjects.IsErr() {
		logger.Warning("Failed to get all objects: ", resultObjects.Error())
		return types.ResultErr[[]models.FormAnswerDB](resultObjects.Error())
	}

	objectsDB := utils.Map(resultObjects.Value(), models.InterfaceTo[models.FormAnswerDB])
	logger.Debug("Retrieved", len(objectsDB), "objects from db")

	return types.ResultOk(objectsDB)
}
