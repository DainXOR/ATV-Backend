package service

import (
	"dainxor/atv/dao"
	"dainxor/atv/logger"
	"dainxor/atv/models"
	"dainxor/atv/types"
	"dainxor/atv/utils"

	"github.com/gin-gonic/gin"
)

type iGenericServicesNS[C, R any, M models.DBModelInterface[C, R], DAO dao.DAOInterface[C, R, M]] struct {
	model M
	dao   DAO
}

func GenericOf[C, R any, M models.DBModelInterface[C, R], DAO dao.DAOInterface[C, R, M]](
	model M,
	dao DAO,
) iGenericServicesNS[C, R, M, DAO] {
	return iGenericServicesNS[C, R, M, DAO]{
		model: model,
		dao:   dao,
	}
}

func (s iGenericServicesNS[C, R, M, DAO]) Create() func(*gin.Context) {
	return func(c *gin.Context) {
		var body C

		if err := c.ShouldBindJSON(&body); err != nil {
			expected := utils.StructToString(body)
			logger.Error(err.Error())
			logger.Errorf("Failed to create %s: JSON request body is invalid", s.model.TableName())
			logger.Error("Expected body: ", expected)

			c.JSON(types.Http.C400().BadRequest(),
				types.EmptyResponse(
					"Invalid request body",
					"Expected body: "+expected,
				),
			)
			return
		}

		logger.Debugf("Creating %s in db: %s", s.model.TableName(), body)

		result := s.dao.Create(body)

		if result.IsErr() {
			logger.Errorf("Failed to create %s in db: %s", s.model.TableName(), result.Error())
			handleErrorAnswer(c, result.Error())
			return
		}

		object := result.Value()
		c.JSON(types.Http.C200().Created(),
			types.Response(
				object.ToResponse(),
				"",
			),
		)
	}
}

func (s iGenericServicesNS[C, R, M, DAO]) GetByID() func(*gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")
		logger.Debug("Getting %s by ID: %s", s.model.TableName(), id)
		filter := models.Filter.Create(c.Request.URL.Query())
		result := s.dao.GetByID(id, filter)

		if result.IsErr() {
			handleErrorAnswer(c, result.Error())
			return
		}

		object := result.Value()
		c.JSON(types.Http.C200().Ok(),
			types.Response(
				object.ToResponse(),
				"",
			),
		)
	}
}
func (s iGenericServicesNS[C, R, M, DAO]) GetAll() func(*gin.Context) {
	return func(c *gin.Context) {
		filter := models.Filter.Create(c.Request.URL.Query())
		result := s.dao.GetAll(filter)

		if result.IsErr() {
			handleErrorAnswer(c, result.Error())
			return
		}

		objects := utils.Map(result.Value(), M.ToResponse)
		if len(objects) == 0 {
			logger.Warningf("No %s found in database", s.model.TableName())
			c.JSON(types.Http.C400().NotFound(),
				types.EmptyResponse(
					"No "+s.model.TableName()+" found",
				))
			return
		}
		c.JSON(types.Http.C200().Ok(),
			types.Response(
				objects,
				"",
			),
		)
	}
}

func (s iGenericServicesNS[C, R, M, DAO]) UpdateByID() func(*gin.Context) {
	return func(c *gin.Context) {
		var body C

		if err := c.ShouldBindJSON(&body); err != nil {
			expected := utils.StructToTagString(body, "json")
			logger.Error(err.Error())
			logger.Errorf("Failed to update %s: JSON request body is invalid", s.model.TableName())
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
		logger.Debugf("Updating %s by ID: ", s.model.TableName(), id)
		filter := models.Filter.Create(c.Request.URL.Query())

		result := s.dao.UpdateByID(id, body, filter)
		if result.IsErr() {
			handleErrorAnswer(c, result.Error())
			return
		}

		alerts := result.Value()
		c.JSON(types.Http.C200().Ok(),
			types.Response(
				alerts.ToResponse(),
				"",
			),
		)
	}
}
func (s iGenericServicesNS[C, R, M, DAO]) UpdateAll() func(*gin.Context) {
	return func(c *gin.Context) {
		var body C

		if err := c.ShouldBindJSON(&body); err != nil {
			expected := utils.StructToTagString(body, "json")
			logger.Error(err.Error())
			logger.Errorf("Failed to update %s: JSON request body is invalid", s.model.TableName())
			logger.Error("Expected body: ", expected)

			c.JSON(types.Http.C400().BadRequest(),
				types.EmptyResponse(
					"Invalid request body",
					"Expected body: "+expected,
				),
			)
			return
		}

		logger.Debugf("Updating %s: ", s.model.TableName())
		filter := models.Filter.Create(c.Request.URL.Query())

		result := s.dao.UpdateAll(body, filter)
		if result.IsErr() {
			handleErrorAnswer(c, result.Error())
			return
		}

		objects := utils.Map(result.Value(), M.ToResponse)
		c.JSON(types.Http.C200().Ok(),
			types.Response(
				objects,
				"",
			),
		)
	}
}

func (s iGenericServicesNS[C, R, M, DAO]) PatchByID() func(*gin.Context) {
	return func(c *gin.Context) {
		var body C

		if err := c.ShouldBindJSON(&body); err != nil {
			expected := utils.StructToTagString(body, "json")
			logger.Error(err.Error())
			logger.Errorf("Failed to patch %s: JSON request body is invalid", s.model.TableName())
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
		result := s.dao.PatchByID(id, body, filter)

		if result.IsErr() {
			handleErrorAnswer(c, result.Error())
			return
		}

		alerts := result.Value()
		c.JSON(types.Http.C200().Ok(),
			types.Response(
				alerts.ToResponse(),
				"",
			),
		)
	}
}
func (s iGenericServicesNS[C, R, M, DAO]) PatchAll() func(*gin.Context) {
	return func(c *gin.Context) {
		var body C

		if err := c.ShouldBindJSON(&body); err != nil {
			expected := utils.StructToTagString(body, "json")
			logger.Error(err.Error())
			logger.Errorf("Failed to patch %s: JSON request body is invalid", s.model.TableName())
			logger.Error("Expected body: ", expected)

			c.JSON(types.Http.C400().BadRequest(),
				types.EmptyResponse(
					"Invalid request body",
					"Expected body: "+expected,
				),
			)
			return
		}

		filter := models.Filter.Create(c.Request.URL.Query())
		result := s.dao.PatchAll(body, filter)

		if result.IsErr() {
			handleErrorAnswer(c, result.Error())
			return
		}

		objects := utils.Map(result.Value(), M.ToResponse)
		c.JSON(types.Http.C200().Ok(),
			types.Response(
				objects,
				"",
			),
		)
	}
}

func (s iGenericServicesNS[C, R, M, DAO]) DeleteByID() func(*gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")
		logger.Debug("Deleting %s by ID: ", s.model.TableName(), id)
		filter := models.Filter.Create(c.Request.URL.Query())

		result := s.dao.DeleteByID(id, filter)

		if result.IsErr() {
			handleErrorAnswer(c, result.Error())
			return
		}

		data := result.Value().ToResponse()

		c.JSON(types.Http.C200().Accepted(),
			types.Response(
				data,
				s.model.TableName()+" marked for deletion",
			),
		)
	}
}
func (s iGenericServicesNS[C, R, M, DAO]) DeleteAll() func(*gin.Context) {
	return func(c *gin.Context) {
		logger.Debug("Deleting %s: ", s.model.TableName())
		filter := models.Filter.Create(c.Request.URL.Query())

		result := s.dao.DeleteAll(filter)

		if result.IsErr() {
			handleErrorAnswer(c, result.Error())
			return
		}

		objects := utils.Map(result.Value(), M.ToResponse)

		c.JSON(types.Http.C200().Accepted(),
			types.Response(
				objects,
				s.model.TableName()+" marked for deletion",
			),
		)
	}
}

func (s iGenericServicesNS[C, R, M, DAO]) ForceDeleteByID() func(c *gin.Context) {
	return func(c *gin.Context) {
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
		logger.Infof("Force deleting %s by ID: ", s.model.TableName(), id)
		filter := models.Filter.Create(c.Request.URL.Query())

		result := s.dao.DeletePermanentByID(id, filter)

		if result.IsErr() {
			handleErrorAnswer(c, result.Error())
			return
		}

		data := result.Value().ToResponse()

		c.JSON(types.Http.C200().Ok(),
			types.Response(
				data,
				s.model.TableName()+" deleted permanently",
			),
		)
	}
}
func (s iGenericServicesNS[C, R, M, DAO]) ForceDeleteAll() func(c *gin.Context) {
	return func(c *gin.Context) {
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

		logger.Infof("Force deleting %s: ", s.model.TableName())
		filter := models.Filter.Create(c.Request.URL.Query())

		results := s.dao.DeletePermanentAll(filter)

		if results.IsErr() {
			handleErrorAnswer(c, results.Error())
			return
		}

		objects := utils.Map(results.Value(), M.ToResponse)

		c.JSON(types.Http.C200().Ok(),
			types.Response(
				objects,
				s.model.TableName()+" deleted permanently",
			),
		)
	}
}

func handleErrorAnswer(c *gin.Context, err error) {
	switch resultError := err.(type) {
	case *types.HttpError:
		handleHttpError(c, *resultError)

	default:
		c.JSON(types.Http.C500().InternalServerError(),
			types.EmptyResponse(
				resultError.Error(),
			),
		)
	}
}

func handleHttpError(c *gin.Context, err types.HttpError) {
	c.JSON(err.Code,
		types.EmptyResponse(
			err.Msg(),
			err.Details(),
		),
	)
}

func tst() {
	GenericOf(models.AlertDB{}, dao.Alert)
	//GetAll()
}
