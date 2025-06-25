package controller

import (
	"dainxor/atv/db"
	"dainxor/atv/logger"
	"dainxor/atv/models"
	"dainxor/atv/types"
	"dainxor/atv/utils"

	"github.com/gin-gonic/gin"
)

type universityType struct{}

var University universityType

func (universityType) Create(c *gin.Context) {
	var body models.UniversityCreate

	if err := c.ShouldBindJSON(&body); err != nil {
		expected := utils.StructToString(body)
		logger.Error(err.Error())
		logger.Error("Failed to create university: JSON request body is invalid")
		logger.Error("Expected body: ", expected)

		c.JSON(types.Http.C400().BadRequest(),
			types.EmptyResponse(
				"Invalid request body",
				"Expected body: "+expected,
			),
		)
		return
	}

	logger.Debug("Creating university in MongoDB: ", body)

	result := db.University.Create(body)

	if result.IsErr() {
		logger.Error("Failed to create university in MongoDB: ", result.Error())
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

	university := result.Value()
	c.JSON(types.Http.C200().Created(),
		types.Response(
			university.ToResponse(),
			"",
		),
	)
}

func (universityType) GetByID(c *gin.Context) {
	id := c.Param("id")
	logger.Debug("Getting university by ID: ", id)

	result := db.University.GetByID(id)

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

	university := result.Value()
	c.JSON(types.Http.C200().Ok(),
		types.Response(
			university.ToResponse(),
			"",
		),
	)
}
func (universityType) GetAll(c *gin.Context) {
	result := db.University.GetAll()

	if result.IsErr() {
		err := result.Error().(*types.HttpError)
		c.JSON(err.Code,
			types.EmptyResponse(
				err.Msg(),
				err.Details(),
			),
		)
		return
	}

	students := utils.Map(result.Value(), models.UniversityDBMongo.ToResponse)
	if len(students) == 0 {
		logger.Warning("No universities found in MongoDB database")
		c.JSON(types.Http.C400().NotFound(),
			types.EmptyResponse(
				"No universities found",
			))
		return
	}
	c.JSON(types.Http.C200().Ok(),
		types.Response(
			students,
			"",
		),
	)
}
