package service

import (
	"dainxor/atv/dao"
	"dainxor/atv/logger"
	"dainxor/atv/models"
	"dainxor/atv/types"
	"dainxor/atv/utils"

	"github.com/gin-gonic/gin"
)

type formsNS struct{}

var Forms formsNS

func (formsNS) Create(c *gin.Context) {
	var body models.FormCreate

	if err := c.ShouldBindJSON(&body); err != nil {
		expected := utils.StructToString(body)
		logger.Error(err.Error())
		logger.Error("Failed to create form: JSON request body is invalid")
		logger.Error("Expected body: ", expected)

		c.JSON(types.Http.C400().BadRequest(),
			types.EmptyResponse(
				"Invalid request body",
				"Expected body: "+expected,
			),
		)
		return
	}

	logger.Debug("Creating form in MongoDB: ", body)

	result := dao.Forms.Create(body)

	if result.IsErr() {
		logger.Error("Failed to create form in MongoDB: ", result.Error())
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

func (formsNS) GetByID(c *gin.Context) {
	id := c.Param("id")
	logger.Debug("Getting form by ID: ", id)
	filter := models.Filter.Create(c.Request.URL.Query())

	result := dao.Forms.GetByID(id, filter)

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

	form := result.Value()
	c.JSON(types.Http.C200().Ok(),
		types.Response(
			form.ToResponse(),
			"",
		),
	)
}
func (formsNS) GetAll(c *gin.Context) {

}

func (formsNS) UpdateByID(c *gin.Context) {}

func (formsNS) PatchByID(c *gin.Context) {}

func (formsNS) DeleteByID(c *gin.Context) {}
