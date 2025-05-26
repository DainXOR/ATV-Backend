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

type userType struct{}

var User userType

func (userType) GetByID(c *gin.Context) {
	id := c.Param("id")
	logger.Debug("Getting user by ID: ", id)

	result := db.User.GetUserByID(id)

	if result.IsErr() {
		err := result.Error()
		cerror := err.(*types.HttpError)
		c.JSON(cerror.Code.AsInt(), err)
		return
	}

	user := result.Value()
	c.JSON(http.StatusOK, user.ToResponse())
}

func (userType) Create(c *gin.Context) {
	var body models.UserCreate

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

	result := db.User.CreateUser(body)

	if result.IsErr() {
		err := result.Error()
		cerror := err.(*types.HttpError)
		c.JSON(cerror.Code.AsInt(), err)
		return
	}

	user := result.Value()
	c.JSON(http.StatusCreated, user.ToResponse())
}

// Update updates an existing user in the database
// This will override zeroed fields
func (userType) Update(c *gin.Context) {
	var body models.UserCreate

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

	result := db.User.UpdateUser(id, body)

	if result.IsErr() {
		err := result.Error()
		cerror := err.(*types.HttpError)
		c.JSON(cerror.Code.AsInt(), err)
		return
	}

	user := result.Value()
	c.JSON(http.StatusOK, user.ToResponse())
}

// Patch updates an existing user in the database
// This will keep previous value in zeroed fields
func (userType) Patch(c *gin.Context) {
	var body models.UserCreate

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

	result := db.User.PatchUser(id, body)

	if result.IsErr() {
		err := result.Error()
		cerror := err.(*types.HttpError)
		c.JSON(cerror.Code.AsInt(), err)
		return
	}

	user := result.Value()
	c.JSON(http.StatusOK, user.ToResponse())
}
