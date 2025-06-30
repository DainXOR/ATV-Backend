package controller

import (
	"dainxor/atv/db"
	"dainxor/atv/logger"
	"dainxor/atv/models"
	"dainxor/atv/types"
	"dainxor/atv/utils"

	"github.com/gin-gonic/gin"
)

type specialityType struct{}

var Speciality specialityType

func (specialityType) Create(c *gin.Context) {
	var body models.SpecialityCreate

	if err := c.ShouldBindJSON(&body); err != nil {
		expected := utils.StructToString(body)
		logger.Error(err.Error())
		logger.Error("Failed to create speciality: JSON request body is invalid")
		logger.Error("Expected body: ", expected)

		c.JSON(types.Http.C400().BadRequest(),
			types.EmptyResponse(
				"Invalid request body",
				"Expected body: "+expected,
			),
		)
		return
	}

	logger.Debug("Creating speciality in MongoDB: ", body)

	result := db.Speciality.Create(body)

	if result.IsErr() {
		logger.Error("Failed to create speciality in MongoDB: ", result.Error())
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

	speciality := result.Value()
	c.JSON(types.Http.C200().Created(),
		types.Response(
			speciality.ToResponse(),
			"",
		),
	)
}

func (specialityType) GetByID(c *gin.Context) {
	id := c.Param("id")
	logger.Debug("Getting speciality by ID: ", id)

	result := db.Speciality.GetByID(id)

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

	speciality := result.Value()
	c.JSON(types.Http.C200().Ok(),
		types.Response(
			speciality.ToResponse(),
			"",
		),
	)
}
func (specialityType) GetAll(c *gin.Context) {
	result := db.Speciality.GetAll()

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

	students := utils.Map(result.Value(), models.SpecialityDBMongo.ToResponse)
	if len(students) == 0 {
		logger.Warning("No specialities found in MongoDB database")
		c.JSON(types.Http.C400().NotFound(),
			types.EmptyResponse(
				"No specialities found",
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
