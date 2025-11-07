package service

import (
	"dainxor/atv/dao"
	"dainxor/atv/logger"
	"dainxor/atv/models"
	"dainxor/atv/types"
	"dainxor/atv/utils"

	"github.com/gin-gonic/gin"
)

type alertsType struct{}

var Alert alertsType

func (alertsType) Create(c *gin.Context) {
	var body models.AlertCreate

	if err := c.ShouldBindJSON(&body); err != nil {
		expected := utils.StructToString(body)
		logger.Error(err.Error())
		logger.Error("Failed to create alerts: JSON request body is invalid")
		logger.Error("Expected body: ", expected)

		c.JSON(types.Http.C400().BadRequest(),
			types.EmptyResponse(
				"Invalid request body",
				"Expected body: "+expected,
			),
		)
		return
	}

	logger.Debug("Creating alerts in MongoDB: ", body)

	result := dao.Alert.Create(body)

	if result.IsErr() {
		logger.Error("Failed to create alerts in MongoDB: ", result.Error())
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

	alerts := result.Value()
	c.JSON(types.Http.C200().Created(),
		types.Response(
			alerts.ToResponse(),
			"",
		),
	)
}

func (alertsType) GetByID(c *gin.Context) {
	id := c.Param("id")
	logger.Debug("Getting alerts by ID: ", id)
	filter := models.Filter.Create(c.Request.URL.Query())

	result := dao.Alert.GetByID(id, filter)

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

	alerts := result.Value()
	c.JSON(types.Http.C200().Ok(),
		types.Response(
			alerts.ToResponse(),
			"",
		),
	)
}
func (alertsType) GetAll(c *gin.Context) {
	filter := models.Filter.Create(c.Request.URL.Query())
	result := dao.Alert.GetAll(filter)

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

	specialities := utils.Map(result.Value(), models.AlertDB.ToResponse)
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
func (alertsType) UpdateByID(c *gin.Context) {
	var body models.AlertCreate

	if err := c.ShouldBindJSON(&body); err != nil {
		expected := utils.StructToString(body)
		logger.Error(err.Error())
		logger.Error("Failed to update alerts: JSON request body is invalid")
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
	logger.Debug("Updating alerts by ID: ", id)
	filter := models.Filter.Create(c.Request.URL.Query())

	result := dao.Alert.UpdateByID(id, body, filter)
	if result.IsErr() {
		err := result.Error()
		cerror := err.(*types.HttpError)
		c.JSON(cerror.Code, err)
		return
	}

	alerts := result.Value()
	c.JSON(types.Http.C200().Ok(),
		types.Response(
			alerts.ToResponse(),
			"",
		),
	)
}

func (alertsType) PatchByID(c *gin.Context) {
	var body models.AlertCreate

	if err := c.ShouldBindJSON(&body); err != nil {
		expected := utils.StructToString(body)
		logger.Error(err.Error())
		logger.Error("Failed to patch alerts: JSON request body is invalid")
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
	result := dao.Alert.PatchByID(id, body, filter)

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

	alerts := result.Value()
	c.JSON(types.Http.C200().Ok(),
		types.Response(
			alerts.ToResponse(),
			"",
		),
	)
}

func (alertsType) DeleteByID(c *gin.Context) {
	id := c.Param("id")
	logger.Debug("Deleting alerts by ID: ", id)
	filter := models.Filter.Create(c.Request.URL.Query())

	result := dao.Alert.DeleteByID(id, filter)

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
			"Alerts marked for deletion",
		),
	)
}
