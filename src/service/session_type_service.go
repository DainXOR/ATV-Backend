package service

import (
	"dainxor/atv/dao"
	"dainxor/atv/logger"
	"dainxor/atv/models"
	"dainxor/atv/types"
	"dainxor/atv/utils"

	"github.com/gin-gonic/gin"
)

type sessionTypeType struct{}

var SessionType sessionTypeType

func (sessionTypeType) Create(c *gin.Context) {
	var body models.SessionTypeCreate

	if err := c.ShouldBindJSON(&body); err != nil {
		expected := utils.StructToString(body)
		logger.Warning(err.Error())
		logger.Warning("Failed to create session type: JSON request body is invalid")
		logger.Warning("Expected body: ", expected)

		c.JSON(types.Http.C400().BadRequest(),
			types.EmptyResponse(
				"Invalid request body",
				"Expected body: "+expected,
			),
		)
		return
	}

	logger.Debug("Creating session type in MongoDB: ", body)
	existent := dao.SessionType.GetAll(models.Filter.Empty())
	if existent.IsOk() && len(existent.Value()) > 0 {
		match := utils.Any(existent.Value(), func(st models.SessionTypeDB) bool {
			return st.Name == body.Name
		})

		if match {
			logger.Info("Session type with the name already exists: ", body.Name)
			c.JSON(types.Http.C400().Conflict(),
				types.EmptyResponse(
					"Session type with this name already exists",
					"Name: "+body.Name,
				),
			)
			return
		}
	}

	result := dao.SessionType.Create(body)

	if result.IsErr() {
		logger.Error("Failed to create session type in MongoDB: ", result.Error())
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

	sessionType := result.Value()
	c.JSON(types.Http.C200().Created(),
		types.Response(
			sessionType.ToResponse(),
			"",
		),
	)
}

func (sessionTypeType) GetByID(c *gin.Context) {
	id := c.Param("id")
	logger.Debug("Getting session type by ID: ", id)
	filter := models.Filter.Create(c.Request.URL.Query())

	result := dao.SessionType.GetByID(id, filter)

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

	sessionType := result.Value()
	c.JSON(types.Http.C200().Ok(),
		types.Response(
			sessionType.ToResponse(),
			"",
		),
	)
}
func (sessionTypeType) GetAll(c *gin.Context) {
	filter := models.Filter.Create(c.Request.URL.Query())
	result := dao.SessionType.GetAll(filter)

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

	sessionTypes := utils.Map(result.Value(), models.SessionTypeDB.ToResponse)
	if len(sessionTypes) == 0 {
		logger.Warning("No session types found in database")
		c.JSON(types.Http.C400().NotFound(),
			types.EmptyResponse(
				"No session types found",
			))
		return
	}
	c.JSON(types.Http.C200().Ok(),
		types.Response(
			sessionTypes,
			"",
		),
	)
}

func (sessionTypeType) UpdateByID(c *gin.Context) {
	var body models.SessionTypeCreate

	if err := c.ShouldBindJSON(&body); err != nil {
		expected := utils.StructToString(body)
		logger.Error(err.Error())
		logger.Error("Failed to update sessionType: JSON request body is invalid")
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
	logger.Debug("Updating sessionType by ID: ", id)
	filter := models.Filter.Create(c.Request.URL.Query())

	result := dao.SessionType.UpdateByID(id, body, filter)
	if result.IsErr() {
		err := result.Error()
		cerror := err.(*types.HttpError)
		c.JSON(cerror.Code, err)
		return
	}

	sessionType := result.Value()
	c.JSON(types.Http.C200().Ok(),
		types.Response(
			sessionType.ToResponse(),
			"",
		),
	)
}

func (sessionTypeType) PatchByID(c *gin.Context) {
	var body models.SessionTypeCreate

	if err := c.ShouldBindJSON(&body); err != nil {
		expected := utils.StructToString(body)
		logger.Error(err.Error())
		logger.Error("Failed to patch sessionType: JSON request body is invalid")
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

	result := dao.SessionType.PatchByID(id, body, filter)

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

	sessionType := result.Value()
	c.JSON(types.Http.C200().Ok(),
		types.Response(
			sessionType.ToResponse(),
			"",
		),
	)
}

func (sessionTypeType) DeleteByID(c *gin.Context) {
	id := c.Param("id")
	logger.Debug("Deleting sessionType by ID: ", id)
	filter := models.Filter.Create(c.Request.URL.Query())

	result := dao.SessionType.DeleteByID(id, filter)

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
			"SessionType marked for deletion",
		),
	)
}
