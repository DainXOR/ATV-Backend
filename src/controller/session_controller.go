package controller

import (
	"dainxor/atv/db"
	"dainxor/atv/logger"
	"dainxor/atv/models"
	"dainxor/atv/types"
	"dainxor/atv/utils"

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

	sessions := utils.Map(result.Value(), models.SessionDBMongo.ToResponse)

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

	sessions := utils.Map(result.Value(), models.SessionDBMongo.ToResponse)
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
