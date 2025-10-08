package service

import (
	"dainxor/atv/db"
	"dainxor/atv/logger"
	"dainxor/atv/models"
	"dainxor/atv/types"
	"dainxor/atv/utils"

	"github.com/gin-gonic/gin"
)

type priorityType struct{}

var Priority priorityType

func (priorityType) Create(c *gin.Context) {
	var body models.PriorityCreate

	if err := c.ShouldBindJSON(&body); err != nil {
		expected := utils.StructToString(body)
		logger.Error(err.Error())
		logger.Error("Failed to create priority: JSON request body is invalid")
		logger.Error("Expected body: ", expected)

		c.JSON(types.Http.C400().BadRequest(),
			types.EmptyResponse(
				"Invalid request body",
				"Expected body: "+expected,
			),
		)
		return
	}

	logger.Debug("Creating priority in MongoDB: ", body)

	result := db.Priority.Create(body)

	if result.IsErr() {
		logger.Error("Failed to create priority in MongoDB: ", result.Error())
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

	priority := result.Value()
	c.JSON(types.Http.C200().Created(),
		types.Response(
			priority.ToResponse(),
			"",
		),
	)
}

func (priorityType) GetByID(c *gin.Context) {
	id := c.Param("id")
	logger.Debug("Getting priority by ID: ", id)
	filter := models.Filter.Create(c.Request.URL.Query())

	result := db.Priority.GetByID(id, filter)

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

	priority := result.Value()
	c.JSON(types.Http.C200().Ok(),
		types.Response(
			priority.ToResponse(),
			"",
		),
	)
}
func (priorityType) GetAll(c *gin.Context) {
	filter := models.Filter.Create(c.Request.URL.Query())
	result := db.Priority.GetAll(filter)

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

	specialities := utils.Map(result.Value(), models.PriorityDB.ToResponse)
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

func (priorityType) UpdateByID(c *gin.Context) {
	var body models.PriorityCreate

	if err := c.ShouldBindJSON(&body); err != nil {
		expected := utils.StructToString(body)
		logger.Error(err.Error())
		logger.Error("Failed to update priority: JSON request body is invalid")
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
	logger.Debug("Updating priority by ID: ", id)
	filter := models.Filter.Create(c.Request.URL.Query())

	result := db.Priority.UpdateByID(id, body, filter)
	if result.IsErr() {
		err := result.Error()
		cerror := err.(*types.HttpError)
		c.JSON(cerror.Code, err)
		return
	}

	priority := result.Value()
	c.JSON(types.Http.C200().Ok(),
		types.Response(
			priority.ToResponse(),
			"",
		),
	)
}

func (priorityType) PatchByID(c *gin.Context) {
	var body models.PriorityCreate

	if err := c.ShouldBindJSON(&body); err != nil {
		expected := utils.StructToString(body)
		logger.Error(err.Error())
		logger.Error("Failed to patch priority: JSON request body is invalid")
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

	result := db.Priority.PatchByID(id, body, filter)

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

	priority := result.Value()
	c.JSON(types.Http.C200().Ok(),
		types.Response(
			priority.ToResponse(),
			"",
		),
	)
}

func (priorityType) DeleteByID(c *gin.Context) {
	id := c.Param("id")
	logger.Debug("Deleting priority by ID: ", id)
	filter := models.Filter.Create(c.Request.URL.Query())

	result := db.Priority.DeleteByID(id, filter)

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
			"Priority marked for deletion",
		),
	)
}
