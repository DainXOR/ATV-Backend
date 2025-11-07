package service

import (
	"dainxor/atv/dao"
	"dainxor/atv/logger"
	"dainxor/atv/models"
	"dainxor/atv/types"
	"dainxor/atv/utils"

	"github.com/gin-gonic/gin"
)

func genericGetAll[M models.DBModelInterface, C any](model M, daoObject dao.DAOInterface[M, C]) func(*gin.Context) {
	return func(c *gin.Context) {
		filter := models.Filter.Create(c.Request.URL.Query())
		result := dao.Alert.GetAll(filter)

		if result.IsErr() {
			handleErrorAnswer(c, result.Error())
			return
		}

		objects := utils.Map(result.Value(), models.AlertDB.ToResponse)
		if len(objects) == 0 {
			logger.Warningf("No %s found in database", model.TableName())
			c.JSON(types.Http.C400().NotFound(),
				types.EmptyResponse(
					"No "+model.TableName()+" found",
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
