package service

import (
	"dainxor/atv/dao"
	"dainxor/atv/logger"
	"dainxor/atv/models"
	"dainxor/atv/types"
	"dainxor/atv/utils"

	"github.com/gin-gonic/gin"
)

type contactReasonType struct{}

var ContactReason contactReasonType

func (contactReasonType) Create(c *gin.Context) {
	var body models.ContactReasonCreate

	if err := c.ShouldBindJSON(&body); err != nil {
		expected := utils.StructToTagString(body, "json")
		logger.Error(err.Error())
		logger.Error("Failed to create contact reason: JSON request body is invalid")
		logger.Error("Expected body: ", expected)

		c.JSON(types.Http.C400().BadRequest(),
			types.EmptyResponse(
				"Invalid request body",
				"Expected body: "+expected,
			),
		)
		return
	}

	logger.Debug("Creating contact reason in MongoDB: ", body)

	result := dao.ContactReason.Create(body)

	if result.IsErr() {
		logger.Error("Failed to create contact reason in MongoDB: ", result.Error())
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

	contactReason := result.Value()
	c.JSON(types.Http.C200().Created(),
		types.Response(
			contactReason.ToResponse(),
			"",
		),
	)
}

func (contactReasonType) GetByID(c *gin.Context) {
	id := c.Param("id")
	logger.Debug("Getting contact reason by ID: ", id)
	filter := models.Filter.Create(c.Request.URL.Query())

	result := dao.ContactReason.GetByID(id, filter)

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

	contactReason := result.Value()
	c.JSON(types.Http.C200().Ok(),
		types.Response(
			contactReason.ToResponse(),
			"",
		),
	)
}
func (contactReasonType) GetAll(c *gin.Context) {
	filter := models.Filter.Create(c.Request.URL.Query())
	result := dao.ContactReason.GetAll(filter)

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

	contactReasons := utils.Map(result.Value(), models.ContactReasonDB.ToResponse)
	if len(contactReasons) == 0 {
		logger.Warning("No contact reasons found in MongoDB database")
		c.JSON(types.Http.C400().NotFound(),
			types.EmptyResponse(
				"No contact reasons found",
			))
		return
	}
	c.JSON(types.Http.C200().Ok(),
		types.Response(
			contactReasons,
			"",
		),
	)
}

func (contactReasonType) UpdateByID(c *gin.Context) {
	var body models.ContactReasonCreate

	if err := c.ShouldBindJSON(&body); err != nil {
		expected := utils.StructToTagString(body, "json")
		logger.Error(err.Error())
		logger.Error("Failed to update contact reason: JSON request body is invalid")
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
	logger.Debug("Updating contact reason by ID: ", id)
	filter := models.Filter.Create(c.Request.URL.Query())

	result := dao.ContactReason.UpdateByID(id, body, filter)
	if result.IsErr() {
		err := result.Error()
		cerror := err.(*types.HttpError)
		c.JSON(cerror.Code, err)
		return
	}

	contactReason := result.Value()
	c.JSON(types.Http.C200().Ok(),
		types.Response(
			contactReason.ToResponse(),
			"",
		),
	)
}

func (contactReasonType) PatchByID(c *gin.Context) {
	var body models.ContactReasonCreate

	if err := c.ShouldBindJSON(&body); err != nil {
		expected := utils.StructToTagString(body, "json")
		logger.Error(err.Error())
		logger.Error("Failed to patch contact reason: JSON request body is invalid")
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

	result := dao.ContactReason.PatchByID(id, body, filter)

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

	contactReason := result.Value()
	c.JSON(types.Http.C200().Ok(),
		types.Response(
			contactReason.ToResponse(),
			"",
		),
	)
}

func (contactReasonType) DeleteByID(c *gin.Context) {
	id := c.Param("id")
	logger.Debug("Deleting contact reason by ID: ", id)
	filter := models.Filter.Create(c.Request.URL.Query())

	result := dao.ContactReason.DeleteByID(id, filter)

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
			"Contact reason marked for deletion",
		),
	)
}
