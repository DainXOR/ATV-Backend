package controller

import (
	"dainxor/atv/db"
	"dainxor/atv/logger"
	"dainxor/atv/models"
	"dainxor/atv/types"
	"dainxor/atv/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type companionType struct{}

var Companion companionType

func (companionType) GetByIDMongo(c *gin.Context) {
	id := c.Param("id")
	//filter := Filter.Create(c.Request.URL.Query())

	logger.Debug("Getting companion by ID: ", id)

	result := db.Companion.GetByID(id)

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
func (companionType) GetAllMongo(c *gin.Context) {
	result := db.Companion.GetAll()

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

	companions := utils.Map(result.Value(), models.CompanionDB.ToResponse)
	if len(companions) == 0 {
		logger.Warning("No companions found in MongoDB database")
		c.JSON(types.Http.C400().NotFound(),
			types.EmptyResponse(
				"No companions found",
			))
		return
	}
	c.JSON(types.Http.C200().Ok(),
		types.Response(
			companions,
			"",
		),
	)
}

func (companionType) CreateMongo(c *gin.Context) {
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

	result := db.Companion.Create(body)

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

func (companionType) UpdateMongo(c *gin.Context) {
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

	result := db.Companion.UpdateByID(id, body)
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

func (companionType) PatchMongo(c *gin.Context) {
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

	result := db.Companion.PatchByID(id, body)

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

	result := db.Companion.DeleteByID(id)

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

	result := db.Companion.DeletePermanentByID(id)

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

func (companionType) GetByIDGorm(c *gin.Context) {
	c.Header("Location", "/api/v1/companion/"+c.Param("id"))
	c.JSON(types.Http.C300().MovedPermanently(),
		types.EmptyResponse(
			logger.DeprecateMsg(types.V("0.1.1"), types.V("0.1.2"), "Use /api/v1/companion/:id instead"),
		),
	)
}
func (companionType) GetAllGorm(c *gin.Context) {
	c.Header("Location", "/api/v1/companion/all")
	c.JSON(http.StatusOK,
		types.EmptyResponse(
			logger.DeprecateMsg(types.V("0.1.1"), types.V("0.1.2"), "Use /api/v1/companion/all instead"),
		),
	)
}

func (companionType) CreateGorm(c *gin.Context) {
	c.Header("Location", "/api/v1/companion")

	c.JSON(types.Http.C300().MovedPermanently(),
		types.EmptyResponse(
			logger.DeprecateMsg(types.V("0.1.1"), types.V("0.1.2"), "Use /api/v1/companion instead"),
		),
	)
}

// UpdateGorm updates an existing companion in the database
// This will override zeroed fields
func (companionType) UpdateGorm(c *gin.Context) {
	c.Header("Location", "/api/v1/companion/"+c.Param("id"))

	c.JSON(types.Http.C300().MovedPermanently(),
		types.EmptyResponse(
			logger.DeprecateMsg(types.V("0.1.1"), types.V("0.1.2"), "Use /api/v1/companion/:id instead"),
		),
	)
}

// PatchGorm updates an existing companion in the database
// This will keep previous value in zeroed fields
func (companionType) PatchGorm(c *gin.Context) {
	c.Header("Location", "/api/v1/companion/"+c.Param("id"))

	c.JSON(types.Http.C300().MovedPermanently(),
		types.EmptyResponse(
			logger.DeprecateMsg(types.V("0.1.1"), types.V("0.1.2"), "Use /api/v1/companion/:id instead"),
		),
	)
}
