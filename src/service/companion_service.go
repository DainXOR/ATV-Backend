package service

import (
	"dainxor/atv/dao"
	"dainxor/atv/logger"
	"dainxor/atv/models"
	"dainxor/atv/types"
	"dainxor/atv/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type companionType struct{}

var Companion companionType

func (companionType) Create(c *gin.Context) {
	var body models.CompanionCreate

	if err := c.ShouldBindJSON(&body); err != nil {
		expected := utils.StructToString(body)
		logger.Error(err.Error())
		logger.Error("Failed to create companion: JSON request body is invalid")
		logger.Error("Expected body: ", expected)

		c.JSON(types.Http.C400().BadRequest(),
			types.EmptyResponse(
				"Invalid request body",
				"Expected body: "+expected,
			),
		)
		return
	}

	logger.Debug("Creating companion in MongoDB: ", body)

	result := dao.Companion.Create(body)

	if result.IsErr() {
		logger.Error("Failed to create companion in MongoDB: ", result.Error())
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

func (companionType) GetByID(c *gin.Context) {
	id := c.Param("id")
	logger.Debug("Getting companion by ID: ", id)
	filter := models.Filter.Create(c.Request.URL.Query())

	result := dao.Companion.GetByID(id, filter)

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

	companion := result.Value()
	c.JSON(http.StatusOK,
		types.Response(
			companion.ToResponse(),
			"",
		),
	)
}
func (companionType) GetAll(c *gin.Context) {
	filter := models.Filter.Create(c.Request.URL.Query())
	result := dao.Companion.GetAll(filter)

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

	if len(result.Value()) == 0 {
		logger.Warning("No companions found in database")
		c.JSON(types.Http.C400().NotFound(),
			types.EmptyResponse(
				"No companions found",
			))
		return
	}

	companionsDB := utils.Map(result.Value(), models.CompanionDB.ToResponse)
	c.JSON(types.Http.C200().Ok(),
		types.Response(
			companionsDB,
			"",
		),
	)
}

func (companionType) UpdateByID(c *gin.Context) {
	var body models.CompanionCreate

	if err := c.ShouldBindJSON(&body); err != nil {
		expected := utils.StructToString(body)
		logger.Error(err.Error())
		logger.Error("Failed to update companion: JSON request body is invalid")
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
	logger.Debug("Updating companion by ID: ", id)
	filter := models.Filter.Create(c.Request.URL.Query())

	result := dao.Companion.UpdateByID(id, body, filter)
	if result.IsErr() {
		err := result.Error()
		cerror := err.(*types.HttpError)
		c.JSON(cerror.Code, err)
		return
	}

	companion := result.Value()
	c.JSON(http.StatusOK,
		types.Response(
			companion.ToResponse(),
			"",
		),
	)
}

func (companionType) PatchByID(c *gin.Context) {
	var body models.CompanionCreate

	if err := c.ShouldBindJSON(&body); err != nil {
		expected := utils.StructToString(body)
		logger.Error(err.Error())
		logger.Error("Failed to patch companion: JSON request body is invalid")
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

	result := dao.Companion.PatchByID(id, body, filter)

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

	companion := result.Value()
	c.JSON(types.Http.C200().Ok(),
		types.Response(
			companion.ToResponse(),
			"",
		),
	)
}

// DeleteByID deletes a companion by ID
func (companionType) DeleteByID(c *gin.Context) {
	id := c.Param("id")
	logger.Debug("Deleting companion by ID: ", id)
	filter := models.Filter.Create(c.Request.URL.Query())

	result := dao.Companion.DeleteByID(id, filter)

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
			"Companion marked for deletion",
		),
	)
}
func (companionType) ForceDeleteByID(c *gin.Context) {
	confirm := c.Param("confirm")
	if confirm != "delete-permanently" {
		c.JSON(types.Http.C400().BadRequest(),
			types.EmptyResponse(
				"Invalid confirmation parameter",
				"Use 'delete-permanently' to confirm deletion",
			),
		)
		return
	}

	id := c.Param("id")
	logger.Info("Force deleting companion by ID: ", id)
	filter := models.Filter.Create(c.Request.URL.Query())

	result := dao.Companion.DeletePermanentByID(id, filter)

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

	c.JSON(types.Http.C200().Ok(),
		types.Response(
			data,
			"Companion deleted permanently",
		),
	)
}
