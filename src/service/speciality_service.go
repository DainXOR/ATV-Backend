package service

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
	filter := models.Filter.Create(c.Request.URL.Query())

	result := db.Speciality.GetByID(id, filter)

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
	filter := models.Filter.Create(c.Request.URL.Query())
	result := db.Speciality.GetAll(filter)

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

	specialities := utils.Map(result.Value(), models.SpecialityDB.ToResponse)
	if len(specialities) == 0 {
		logger.Warning("No specialities found in MongoDB database")
		c.JSON(types.Http.C400().NotFound(),
			types.EmptyResponse(
				"No specialities found",
			))
		return
	}
	c.JSON(types.Http.C200().Ok(),
		types.Response(
			specialities,
			"",
		),
	)
}

func (specialityType) UpdateByID(c *gin.Context) {
	var body models.SpecialityCreate

	if err := c.ShouldBindJSON(&body); err != nil {
		expected := utils.StructToString(body)
		logger.Error(err.Error())
		logger.Error("Failed to update speciality: JSON request body is invalid")
		logger.Error("Expected body: ", expected)

		c.JSON(types.Http.C400().BadRequest(),
			types.EmptyResponse(
				"Invalid request body",
				"Expected body: "+expected,
			),
		)
		return
	}

	id := c.Param("id")
	logger.Debug("Updating speciality by ID: ", id)
	filter := models.Filter.Create(c.Request.URL.Query())

	result := db.Speciality.UpdateByID(id, body, filter)
	if result.IsErr() {
		err := result.Error()
		cerror := err.(*types.HttpError)
		c.JSON(cerror.Code, err)
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

func (specialityType) PatchByID(c *gin.Context) {
	var body models.SpecialityCreate

	if err := c.ShouldBindJSON(&body); err != nil {
		expected := utils.StructToString(body)
		logger.Error(err.Error())
		logger.Error("Failed to patch speciality: JSON request body is invalid")
		logger.Error("Expected body: ", expected)

		c.JSON(types.Http.C400().BadRequest(),
			types.EmptyResponse(
				"Invalid request body",
				"Expected body: "+expected,
			),
		)
		return
	}

	id := c.Param("id")
	filter := models.Filter.Create(c.Request.URL.Query())

	result := db.Speciality.PatchByID(id, body, filter)

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

func (specialityType) DeleteByID(c *gin.Context) {
	id := c.Param("id")
	logger.Debug("Deleting speciality by ID: ", id)
	filter := models.Filter.Create(c.Request.URL.Query())

	result := db.Speciality.DeleteByID(id, filter)

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

	data := result.Value().ToResponse()

	c.JSON(types.Http.C200().Accepted(),
		types.Response(
			data,
			"Speciality marked for deletion",
		),
	)
}
