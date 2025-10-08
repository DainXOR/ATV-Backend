package controller

import (
	"dainxor/atv/configs"
	"dainxor/atv/logger"
	"dainxor/atv/service"
	"dainxor/atv/types"
	"fmt"

	"github.com/gin-gonic/gin"
)

func SessionTypesRoutes(router *gin.Engine) {
	rv := configs.App.RoutesVersion()
	beforeRoute := fmt.Sprintf("/api/v%d/session-type", rv-1)
	lastRoute := fmt.Sprintf("/api/v%d/session-types", rv)

	router.Group(beforeRoute).Any("", func(ctx *gin.Context) {
		ctx.JSON(types.Http.C300().PermanentRedirect(),
			types.EmptyResponse(
				logger.DeprecateMsg(types.V("0.2.0"), types.V("0.3.0"), "Use", lastRoute, "instead"),
			),
		)
		ctx.Redirect(types.Http.C300().PermanentRedirect(), lastRoute)
	})

	sessionTypeRouter := router.Group(lastRoute)
	{
		sessionTypeRouter.POST("/", service.SessionType.Create)

		sessionTypeRouter.GET("/:id", service.SessionType.GetByID)
		sessionTypeRouter.GET("/all", service.SessionType.GetAll)

		sessionTypeRouter.PUT("/:id", service.SessionType.UpdateByID)

		sessionTypeRouter.PATCH("/:id", service.SessionType.PatchByID)

		sessionTypeRouter.DELETE("/:id", service.SessionType.DeleteByID)
	}
}
