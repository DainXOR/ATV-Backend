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

type sessionType struct{}

var Session sessionType

func (sessionType) Create(c *gin.Context) {
	var body models.SessionCreate

	if err := c.ShouldBindJSON(&body); err != nil {
		expected := utils.StructToString(body)
		logger.Error(err.Error())
		logger.Error("Failed to create session: JSON request body is invalid")
		logger.Error("Expected body: ", expected)

		c.JSON(types.Http.C400().BadRequest(),
			types.EmptyResponse(
				"Invalid request body",
				"Expected body: "+expected,
			),
		)
		return
	}

	logger.Debug("Creating session in MongoDB: ", body)

	result := db.Session.Create(body)

	if result.IsErr() {
		logger.Warning("Failed to create session in MongoDB: ", result.Error())
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

	session := result.Value()
	c.JSON(types.Http.C200().Created(),
		types.Response(
			session.ToResponse(),
			"",
		),
	)
}

func (sessionType) GetByID(c *gin.Context) {
	id := c.Param("id")
	logger.Debug("Getting session by ID: ", id)

	result := db.Session.GetByID(id)

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

	session := result.Value()
	c.JSON(types.Http.C200().Ok(),
		types.Response(
			session.ToResponse(),
			"",
		),
	)
}
func (sessionType) GetAllByStudentID(c *gin.Context) {
	studentID := c.Param("student_id")
	logger.Debug("Getting all sessions by student ID: ", studentID)

	result := db.Session.GetAllByStudentID(studentID)

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

	sessions := utils.Map(result.Value(), models.SessionDB.ToResponse)

	if len(sessions) == 0 {
		logger.Warning("No sessions found for student ID in MongoDB database")
		c.JSON(types.Http.C400().NotFound(),
			types.EmptyResponse(
				"No sessions found for student ID",
			))
		return
	}
	c.JSON(types.Http.C200().Ok(),
		types.Response(
			sessions,
			"",
		),
	)
}
func (sessionType) GetAll(c *gin.Context) {
	result := db.Session.GetAll()

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

	sessions := utils.Map(result.Value(), models.SessionDB.ToResponse)
	if len(sessions) == 0 {
		logger.Warning("No sessions found in MongoDB database")
		c.JSON(types.Http.C400().NotFound(),
			types.EmptyResponse(
				"No sessions found",
			))
		return
	}
	c.JSON(types.Http.C200().Ok(),
		types.Response(
			sessions,
			"",
		),
	)
}

func (sessionType) UpdateByID(c *gin.Context) {
	var body models.SessionCreate

	if err := c.ShouldBindJSON(&body); err != nil {
		expected := utils.StructToString(body)
		logger.Error(err.Error())
		logger.Error("Failed to update session: JSON request body is invalid")
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
	logger.Debug("Updating session by ID: ", id)

	result := db.Session.UpdateByID(id, body)
	if result.IsErr() {
		err := result.Error()
		cerror := err.(*types.HttpError)
		c.JSON(cerror.Code, err)
		return
	}

	session := result.Value()
	c.JSON(http.StatusOK,
		types.Response(
			session.ToResponse(),
			"",
		),
	)
}

func (sessionType) PatchByID(c *gin.Context) {
	var body models.SessionCreate

	if err := c.ShouldBindJSON(&body); err != nil {
		expected := utils.StructToString(body)
		logger.Error(err.Error())
		logger.Error("Failed to patch session: JSON request body is invalid")
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

	result := db.Session.PatchByID(id, body)

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

	session := result.Value()
	c.JSON(types.Http.C200().Ok(),
		types.Response(
			session.ToResponse(),
			"",
		),
	)
}

func (sessionType) DeleteByID(c *gin.Context) {
	id := c.Param("id")
	logger.Debug("Deleting session by ID: ", id)

	result := db.Session.DeleteByID(id)

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

	session := result.Value()
	c.JSON(types.Http.C200().Ok(),
		types.Response(
			session.ToResponse(),
			"Session marked for deletion",
		),
	)
}
