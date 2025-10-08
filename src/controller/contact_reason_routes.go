package controller

import (
	"dainxor/atv/configs"
	"dainxor/atv/logger"
	"dainxor/atv/service"
	"dainxor/atv/types"
	"fmt"

	"github.com/gin-gonic/gin"
)

func ContactReasonsRoutes(router *gin.Engine) {
	rv := configs.App.RoutesVersion()
	beforeRoute := fmt.Sprintf("/api/v%d/contact-reasons", rv-1)
	lastRoute := fmt.Sprintf("/api/v%d/contact-reasons", rv)

	router.Group(beforeRoute).Any("", func(ctx *gin.Context) {
		ctx.JSON(types.Http.C300().PermanentRedirect(),
			types.EmptyResponse(
				logger.DeprecateMsg(types.V("0.2.0"), types.V("0.3.0"), "Use", lastRoute, "instead"),
			),
		)
		ctx.Redirect(types.Http.C300().PermanentRedirect(), lastRoute)
	})

	contactReasonRouter := router.Group(lastRoute)
	{
		contactReasonRouter.POST("/", service.ContactReason.Create)

		contactReasonRouter.GET("/:id", service.ContactReason.GetByID)
		contactReasonRouter.GET("/all", service.ContactReason.GetAll)

		contactReasonRouter.PUT("/:id", service.ContactReason.UpdateByID)

		contactReasonRouter.PATCH("/:id", service.ContactReason.PatchByID)

		contactReasonRouter.DELETE("/:id", service.ContactReason.DeleteByID)
	}
}
