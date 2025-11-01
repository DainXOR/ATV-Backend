package service

import (
	"dainxor/atv/db"
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

	logger.Debug("Creating form question in MongoDB: ", body)

	result := db.FormQuestions.Create(body)

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

func (formQuestionsNS) GetByID() {}
func (formQuestionsNS) GetAll()  {}

func (formQuestionsNS) UpdateByID() {}

func (formQuestionsNS) PatchByID() {}

func (formQuestionsNS) DeleteByID() {}
