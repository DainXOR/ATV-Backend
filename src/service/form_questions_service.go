package service

import (
	"dainxor/atv/dao"
	"dainxor/atv/logger"
	"dainxor/atv/models"
	"dainxor/atv/types"
	"dainxor/atv/utils"

	"github.com/gin-gonic/gin"
)

type formQuestionsNS struct{}

var FormQuestions formQuestionsNS

func (formQuestionsNS) Create(c *gin.Context) {
	var body models.FormQuestionCreate

	if err := c.ShouldBindJSON(&body); err != nil {
		expected := utils.StructToString(body)
		logger.Error(err.Error())
		logger.Error("Failed to create form question: JSON request body is invalid")
		logger.Error("Expected body: ", expected)

		c.JSON(types.Http.C400().BadRequest(),
			types.EmptyResponse(
				"Invalid request body",
				"Expected body: "+expected,
			),
		)
		return
	}

	logger.Debug("Creating form question in db: ", body)

	result := dao.FormQuestions.Create(body)

	if result.IsErr() {
		logger.Error("Failed to create form question in db: ", result.Error())
		handleErrorAnswer(c, result.Error())
		return
	}

	object := result.Value()
	c.JSON(types.Http.C200().Created(),
		types.Response(
			object.ToResponse(),
			"",
		),
	)
}

func (formQuestionsNS) GetByID(c *gin.Context) {
	id := c.Param("id")
	logger.Debug("Getting question by ID: ", id)
	filter := models.Filter.Create(c.Request.URL.Query())

	result := dao.FormQuestions.GetByID(id, filter)

	if result.IsErr() {
		handleErrorAnswer(c, result.Error())
		return
	}

	question := result.Value()
	c.JSON(types.Http.C200().Ok(),
		types.Response(
			question.ToResponse(),
			"",
		),
	)
}
func (formQuestionsNS) GetAll(c *gin.Context) {
	filter := models.Filter.Create(c.Request.URL.Query())
	result := dao.FormQuestions.GetAll(filter)

	if result.IsErr() {
		handleErrorAnswer(c, result.Error())
		return
	}

	objects := utils.Map(result.Value(), models.FormQuestionDB.ToResponse)
	if len(objects) == 0 {
		logger.Warning("No question found")
		c.JSON(types.Http.C400().NotFound(),
			types.EmptyResponse(
				"No question found",
			))
		return
	}
	c.JSON(types.Http.C200().Ok(),
		types.Response(
			objects,
			"",
		),
	)
}

func (formQuestionsNS) UpdateByID(c *gin.Context) {}

func (formQuestionsNS) PatchByID(c *gin.Context) {}

func (formQuestionsNS) DeleteByID(c *gin.Context) {}
