package service

import (
	"dainxor/atv/dao"
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

	result := dao.University.Create(body)

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
	filter := models.Filter.Create(c.Request.URL.Query())

	result := dao.University.GetByID(id, filter)

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
	filter := models.Filter.Create(c.Request.URL.Query())
	result := dao.University.GetAll(filter)

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

	universities := utils.Map(result.Value(), models.UniversityDB.ToResponse)
	if len(universities) == 0 {
		logger.Warning("No universities found in MongoDB database")
		c.JSON(types.Http.C400().NotFound(),
			types.EmptyResponse(
				"No universities found",
			))
		return
	}
	c.JSON(types.Http.C200().Ok(),
		types.Response(
			universities,
			"",
		),
	)
}

func (universityType) UpdateByID(c *gin.Context) {
	var body models.UniversityCreate

	if err := c.ShouldBindJSON(&body); err != nil {
		expected := utils.StructToString(body)
		logger.Error(err.Error())
		logger.Error("Failed to update university: JSON request body is invalid")
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
	logger.Debug("Updating university by ID: ", id)
	filter := models.Filter.Create(c.Request.URL.Query())

	result := dao.University.UpdateByID(id, body, filter)
	if result.IsErr() {
		err := result.Error()
		cerror := err.(*types.HttpError)
		c.JSON(cerror.Code, err)
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

func (universityType) PatchByID(c *gin.Context) {
	var body models.UniversityCreate

	if err := c.ShouldBindJSON(&body); err != nil {
		expected := utils.StructToString(body)
		logger.Error(err.Error())
		logger.Error("Failed to patch university: JSON request body is invalid")
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

	result := dao.University.PatchByID(id, body, filter)

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

func (universityType) DeleteByID(c *gin.Context) {
	id := c.Param("id")
	logger.Debug("Deleting university by ID: ", id)
	filter := models.Filter.Create(c.Request.URL.Query())

	result := dao.University.DeleteByID(id, filter)

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
			"University marked for deletion",
		),
	)
}
