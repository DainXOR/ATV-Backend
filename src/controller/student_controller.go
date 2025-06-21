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

type studentType struct{}

var Student studentType

func (studentType) GetByIDGorm(c *gin.Context) {
	logger.Debug("Using GORM")
	id := c.Param("id")
	logger.Debug("Getting user by ID: ", id)

	result := db.Student.GetByIDGorm(id)

	if result.IsErr() {
		err := result.Error()
		cerror := err.(*types.HttpError)
		c.JSON(cerror.Code, err)
		return
	}

	user := result.Value()
	c.JSON(types.Http.C200().Ok(),
		models.Response(
			user.ToResponse(),
			logger.DeprecateMsg(1, 2, "Use /api/v1/student/:id instead"),
		),
	)
}
func (studentType) GetByIDMongo(c *gin.Context) {
	logger.Debug("Using MongoDB")
	id := c.Param("id")
	logger.Debug("Getting student by ID: ", id)

	result := db.Student.GetByIDMongo(id)

	if result.IsErr() {
		err := result.Error()
		cerror := err.(*types.HttpError)
		c.JSON(cerror.Code, err)
		return
	}

	user := result.Value()
	c.JSON(http.StatusOK,
		models.Response(
			user.ToResponse(),
			"",
		),
	)
}

func (studentType) GetAllGorm(c *gin.Context) {
	logger.Debug("Using GORM")

	result := db.Student.GetAllGorm()

	if result.IsErr() {
		err := result.Error()
		cerror := err.(*types.HttpError)
		c.JSON(cerror.Code, err)
		return
	}

	users := utils.Map(result.Value(), models.UserDBGorm.ToResponse)
	if len(users) == 0 {
		logger.Warning("No users found in GORM database")
		c.JSON(http.StatusNotFound, types.Error(
			types.Http.NotFound(),
			"No users found",
			"No users found in the GORM database",
		))
		return
	}
	c.JSON(http.StatusOK,
		models.Response(
			users,
			logger.DeprecateMsg(1, 2, "Use /api/v1/student/all instead"),
		),
	)
}
func (studentType) GetAllMongo(c *gin.Context) {
	logger.Debug("Using MongoDB")

	result := db.Student.GetAllMongo()

	if result.IsErr() {
		err := result.Error()
		cerror := err.(*types.HttpError)
		c.JSON(cerror.Code, err)
		return
	}

	students := utils.Map(result.Value(), models.StudentDBMongo.ToResponse)
	if len(students) == 0 {
		logger.Warning("No students found in MongoDB database")
		c.JSON(http.StatusNotFound, types.Error(
			types.Http.NotFound(),
			"No students found",
			"No students found in the MongoDB database",
		))
		return
	}
	c.JSON(http.StatusOK,
		models.Response(
			students,
			"",
		),
	)
}

func (studentType) CreateGorm(c *gin.Context) {
	logger.Debug("Using GORM")
	var body models.StudentCreate

	if err := c.ShouldBindJSON(&body); err != nil {
		expected := utils.StructToString(body)
		logger.Error(err.Error())
		logger.Error("Failed to create user: JSON request body is invalid")
		logger.Error("Expected body: ", expected)

		c.JSON(http.StatusBadRequest,
			types.Error(
				types.Http.BadRequest(),
				"Invalid request body",
				"Expected body: "+expected,
			),
		)
		return
	}

	logger.Debug("Creating user: ", body)
	logger.Debug("Username: ", body.FirstName, body.LastName)

	result := db.Student.CreateGorm(body)

	if result.IsErr() {
		err := result.Error()
		cerror := err.(*types.HttpError)
		c.JSON(cerror.Code, err)
		return
	}

	user := result.Value()
	c.JSON(http.StatusCreated,
		models.Response(
			user.ToResponse(),
			logger.DeprecateMsg(1, 2, "Use /api/v1/student instead"),
		),
	)
}
func (studentType) CreateMongo(c *gin.Context) {
	logger.Debug("Using MongoDB")
	var body models.StudentCreate

	if err := c.ShouldBindJSON(&body); err != nil {
		expected := utils.StructToString(body)
		logger.Error(err.Error())
		logger.Error("Failed to create student: JSON request body is invalid")
		logger.Error("Expected body: ", expected)

		c.JSON(http.StatusBadRequest,
			types.Error(
				types.Http.BadRequest(),
				"Invalid request body",
				"Expected body: "+expected,
			),
		)
		return
	}

	logger.Debug("Creating student in MongoDB: ", body)

	result := db.Student.CreateMongo(body)

	if result.IsErr() {
		logger.Error("Failed to create student in MongoDB: ", result.Error())
		err := result.Error()
		httpErr := err.(*types.HttpError)
		c.JSON(httpErr.Code, err)
		return
	}

	user := result.Value()
	c.JSON(http.StatusCreated,
		models.Response(
			user.ToResponse(),
			"",
		),
	)
}

// Update updates an existing user in the database
// This will override zeroed fields
func (studentType) Update(c *gin.Context) {
	var body models.StudentCreate

	if err := c.ShouldBindJSON(&body); err != nil {
		expected := utils.StructToString(body)
		logger.Error(err.Error())
		logger.Error("Failed to update user: JSON request body is invalid")
		logger.Error("Expected body: ", expected)

		c.JSON(http.StatusBadRequest,
			types.Error(
				types.Http.BadRequest(),
				"Invalid request body",
				"Expected body: "+expected,
			),
		)
		return
	}

	id := c.Param("id")
	logger.Debug("Updating user by ID: ", id)

	result := db.Student.UpdateGorm(id, body)

	if result.IsErr() {
		err := result.Error()
		cerror := err.(*types.HttpError)
		c.JSON(cerror.Code, err)
		return
	}

	user := result.Value()
	c.JSON(http.StatusOK,
		models.Response(
			user.ToResponse(),
			logger.DeprecateMsg(1, 2, "Use /api/v1/student/:id instead"),
		),
	)
}

// Patch updates an existing user in the database
// This will keep previous value in zeroed fields
func (studentType) Patch(c *gin.Context) {
	var body models.StudentCreate

	if err := c.ShouldBindJSON(&body); err != nil {
		expected := utils.StructToString(body)
		logger.Error(err.Error())
		logger.Error("Failed to update user: JSON request body is invalid")
		logger.Error("Expected body: ", expected)

		c.JSON(http.StatusBadRequest,
			types.Error(
				types.Http.BadRequest(),
				"Invalid request body",
				"Expected body: "+expected,
			),
		)
		return
	}

	id := c.Param("id")
	logger.Debug("Patching user by ID: ", id)

	result := db.Student.PatchGorm(id, body)

	if result.IsErr() {
		err := result.Error()
		cerror := err.(*types.HttpError)
		c.JSON(cerror.Code, err)
		return
	}

	user := result.Value()
	c.JSON(http.StatusOK,
		models.Response(
			user.ToResponse(),
			logger.DeprecateMsg(1, 2, "Use /api/v1/student/:id instead"),
		),
	)
}
