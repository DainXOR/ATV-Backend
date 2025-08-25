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

func (studentType) GetByIDMongo(c *gin.Context) {
	id := c.Param("id")
	logger.Debug("Getting student by ID: ", id)

	result := db.Student.GetByID(id)

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

	student := result.Value()
	c.JSON(http.StatusOK,
		types.Response(
			student.ToResponse(),
			"",
		),
	)
}
func (studentType) GetAllMongo(c *gin.Context) {
	result := db.Student.GetAll()

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

	students := utils.Map(result.Value(), models.StudentDB.ToResponse)
	if len(students) == 0 {
		logger.Warning("No students found in MongoDB database")
		c.JSON(types.Http.C400().NotFound(),
			types.EmptyResponse(
				"No students found",
			))
		return
	}
	c.JSON(types.Http.C200().Ok(),
		types.Response(
			students,
			"",
		),
	)
}

func (studentType) CreateMongo(c *gin.Context) {
	var body models.StudentCreate

	if err := c.ShouldBindJSON(&body); err != nil {
		expected := utils.StructToString(body)
		logger.Error(err.Error())
		logger.Error("Failed to create student: JSON request body is invalid")
		logger.Error("Expected body: ", expected)

		c.JSON(types.Http.C400().BadRequest(),
			types.EmptyResponse(
				"Invalid request body",
				"Expected body: "+expected,
			),
		)
		return
	}

	logger.Debug("Creating student in MongoDB: ", body)

	result := db.Student.Create(body)

	if result.IsErr() {
		logger.Error("Failed to create student in MongoDB: ", result.Error())
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

	student := result.Value()
	c.JSON(types.Http.C200().Created(),
		types.Response(
			student.ToResponse(),
			"",
		),
	)
}

func (studentType) UpdateMongo(c *gin.Context) {
	var body models.StudentCreate

	if err := c.ShouldBindJSON(&body); err != nil {
		expected := utils.StructToString(body)
		logger.Error(err.Error())
		logger.Error("Failed to update student: JSON request body is invalid")
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
	logger.Debug("Updating student by ID: ", id)

	result := db.Student.UpdateByID(id, body)
	if result.IsErr() {
		err := result.Error()
		cerror := err.(*types.HttpError)
		c.JSON(cerror.Code, err)
		return
	}

	student := result.Value()
	c.JSON(http.StatusOK,
		types.Response(
			student.ToResponse(),
			"",
		),
	)
}

func (studentType) PatchMongo(c *gin.Context) {
	var body models.StudentCreate

	if err := c.ShouldBindJSON(&body); err != nil {
		expected := utils.StructToString(body)
		logger.Error(err.Error())
		logger.Error("Failed to patch student: JSON request body is invalid")
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

	result := db.Student.PatchByID(id, body)

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

	student := result.Value()
	c.JSON(types.Http.C200().Ok(),
		types.Response(
			student.ToResponse(),
			"",
		),
	)
}

// DeleteByID deletes a student by ID
func (studentType) DeleteByID(c *gin.Context) {
	id := c.Param("id")
	logger.Debug("Deleting student by ID: ", id)

	result := db.Student.DeleteByID(id)

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
			"Student marked for deletion",
		),
	)
}
func (studentType) ForceDeleteByID(c *gin.Context) {
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
	logger.Info("Force deleting student by ID: ", id)

	result := db.Student.DeletePermanentByID(id)

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
			"Student deleted permanently",
		),
	)
}

func (studentType) GetByIDGorm(c *gin.Context) {
	c.Header("Location", "/api/v1/student/"+c.Param("id"))
	c.JSON(types.Http.C300().MovedPermanently(),
		types.EmptyResponse(
			logger.DeprecateMsg(types.V("0.0.3"), types.V("0.1.1"), "Use /api/v1/student/:id instead"),
		),
	)
}
func (studentType) GetAllGorm(c *gin.Context) {
	c.Header("Location", "/api/v1/student/all")
	c.JSON(http.StatusOK,
		types.EmptyResponse(
			logger.DeprecateMsg(types.V("0.0.3"), types.V("0.1.1"), "Use /api/v1/student/all instead"),
		),
	)
}

func (studentType) CreateGorm(c *gin.Context) {
	c.Header("Location", "/api/v1/student")

	c.JSON(types.Http.C300().MovedPermanently(),
		types.EmptyResponse(
			logger.DeprecateMsg(types.V("0.0.3"), types.V("0.1.1"), "Use /api/v1/student instead"),
		),
	)
}

// UpdateGorm updates an existing student in the database
// This will override zeroed fields
func (studentType) UpdateGorm(c *gin.Context) {
	c.Header("Location", "/api/v1/student/"+c.Param("id"))

	c.JSON(types.Http.C300().MovedPermanently(),
		types.EmptyResponse(
			logger.DeprecateMsg(types.V("0.0.3"), types.V("0.1.1"), "Use /api/v1/student/:id instead"),
		),
	)
}

// PatchGorm updates an existing student in the database
// This will keep previous value in zeroed fields
func (studentType) PatchGorm(c *gin.Context) {
	c.Header("Location", "/api/v1/student/"+c.Param("id"))

	c.JSON(types.Http.C300().MovedPermanently(),
		types.EmptyResponse(
			logger.DeprecateMsg(types.V("0.0.3"), types.V("0.1.1"), "Use /api/v1/student/:id instead"),
		),
	)
}
