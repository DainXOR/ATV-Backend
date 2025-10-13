package controller

import (
	"dainxor/atv/configs"
	"dainxor/atv/logger"
	"dainxor/atv/service"
	"dainxor/atv/types"
	"fmt"

	"github.com/gin-gonic/gin"
)

func AlertsRoutes(router *gin.Engine) {
	rv := configs.App.RoutesVersion()
	beforeRoute := fmt.Sprintf("/api/v%d/alerts", rv-1)
	lastRoute := fmt.Sprintf("/api/v%d/alerts", rv)

	router.Group(beforeRoute).Any("", func(ctx *gin.Context) {
		ctx.JSON(types.Http.C300().PermanentRedirect(),
			types.EmptyResponse(
				logger.DeprecateMsg(types.V("0.2.0"), types.V("0.3.0"), "Use", lastRoute, "instead"),
			),
		)
		ctx.Redirect(types.Http.C300().PermanentRedirect(), lastRoute)
	})

	alertsRouter := router.Group(lastRoute)
	{
		alertsRouter.POST("/", service.Alert.Create)

		alertsRouter.GET("/:id", service.Alert.GetByID)
		alertsRouter.GET("/all", service.Alert.GetAll)

		alertsRouter.PUT("/:id", service.Alert.UpdateByID)

		alertsRouter.PATCH("/:id", service.Alert.PatchByID)

		alertsRouter.DELETE("/:id", service.Alert.DeleteByID)
	}
}
