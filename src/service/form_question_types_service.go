package service

import (
	"dainxor/atv/dao"
	"dainxor/atv/logger"
	"dainxor/atv/models"
	"dainxor/atv/types"
	"dainxor/atv/utils"

	"github.com/gin-gonic/gin"
)

type formQuestionTypesNS struct{}

var FormQuestionTypes formQuestionTypesNS

func (formQuestionTypesNS) Create(c *gin.Context) {
	var body models.FormQuestionTypeCreate

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

	logger.Debug("Creating form question in MongoDB: ", body)

	result := dao.FormQuestionTypes.Create(body)

	if result.IsErr() {
		logger.Error("Failed to create form question in MongoDB: ", result.Error())
		err := result.Error()
		httpErr := err.(*types.HttpError)
		c.JSON(httpErr.Code,
			types.EmptyResponse(
				httpErr.Msg(),
				httpErr.Details(),
			),
		)
		return
	}

	companion := result.Value()
	c.JSON(types.Http.C200().Created(),
		types.Response(
			companion.ToResponse(),
			"",
		),
	)
}

func (formQuestionTypesNS) GetByID(c *gin.Context) {
	id := c.Param("id")
	logger.Debug("Getting question by ID: ", id)
	filter := models.Filter.Create(c.Request.URL.Query())

	result := dao.FormQuestionTypes.GetByID(id, filter)

	if result.IsErr() {
		err := result.Error()
		cerror := err.(*types.HttpError)
		c.JSON(cerror.Code,
			types.EmptyResponse(
				cerror.Msg(),
				cerror.Details(),
			),
		)
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
func (formQuestionTypesNS) GetAll(c *gin.Context) {
	filter := models.Filter.Create(c.Request.URL.Query())
	result := dao.FormQuestionTypes.GetAll(filter)

	if result.IsErr() {
		handleErrorAnswer(c, result.Error())
		return
	}

	objects := utils.Map(result.Value(), models.FormQuestionTypeDB.ToResponse)
	if len(objects) == 0 {
		logger.Warning("No question types found")
		c.JSON(types.Http.C400().NotFound(),
			types.EmptyResponse(
				"No question types found",
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

func (formQuestionTypesNS) UpdateByID(c *gin.Context) {}

func (formQuestionTypesNS) PatchByID(c *gin.Context) {}

func (formQuestionTypesNS) DeleteByID(c *gin.Context) {}
