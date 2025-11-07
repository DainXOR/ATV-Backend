package service

import (
	"dainxor/atv/dao"
	"dainxor/atv/logger"
	"dainxor/atv/models"
	"dainxor/atv/types"
	"dainxor/atv/utils"

	"github.com/gin-gonic/gin"
)

type formAnswersNS struct{}

var FormAnswers formAnswersNS

func (formAnswersNS) Create(c *gin.Context) {
	var body models.FormAnswerCreate

	if err := c.ShouldBindJSON(&body); err != nil {
		expected := utils.StructToString(body)
		logger.Error(err.Error())
		logger.Error("Failed to create form answer: JSON request body is invalid")
		logger.Error("Expected body: ", expected)

		c.JSON(types.Http.C400().BadRequest(),
			types.EmptyResponse(
				"Invalid request body",
				"Expected body: "+expected,
			),
		)
		return
	}

	logger.Debug("Creating form answer in db: ", body)

	result := dao.FormAnswers.Create(body)

	if result.IsErr() {
		logger.Error("Failed to create form answer in db: ", result.Error())
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

func (formAnswersNS) GetByID(c *gin.Context) {
	id := c.Param("id")
	logger.Debug("Getting answer by ID: ", id)
	filter := models.Filter.Create(c.Request.URL.Query())

	result := dao.FormAnswers.GetByID(id, filter)

	if result.IsErr() {
		handleErrorAnswer(c, result.Error())
		return
	}

	answer := result.Value()
	c.JSON(types.Http.C200().Ok(),
		types.Response(
			answer.ToResponse(),
			"",
		),
	)
}
func (formAnswersNS) GetAll(c *gin.Context) {
	filter := models.Filter.Create(c.Request.URL.Query())
	result := dao.FormAnswers.GetAll(filter)

	if result.IsErr() {
		handleErrorAnswer(c, result.Error())
		return
	}

	objects := utils.Map(result.Value(), models.FormAnswerDB.ToResponse)
	if len(objects) == 0 {
		logger.Warning("No answer found")
		c.JSON(types.Http.C400().NotFound(),
			types.EmptyResponse(
				"No answer found",
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

func (formAnswersNS) UpdateByID(c *gin.Context) {}

func (formAnswersNS) PatchByID(c *gin.Context) {}

func (formAnswersNS) DeleteByID(c *gin.Context) {}
