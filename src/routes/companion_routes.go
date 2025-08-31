package routes

import (
	"dainxor/atv/configs"
	"dainxor/atv/controller"
	"dainxor/atv/logger"
	"dainxor/atv/types"
	"fmt"

	"github.com/gin-gonic/gin"
)

func CompanionRoutes(router *gin.Engine) {
	rv := configs.App.RoutesVersion()
	beforeRoute := fmt.Sprintf("/api/v%d/companion", rv-1)
	lastRoute := fmt.Sprintf("/api/v%d/companion", rv)

	router.Group(beforeRoute).Any("", func(ctx *gin.Context) {
		ctx.JSON(types.Http.C300().PermanentRedirect(),
			types.EmptyResponse(
				logger.DeprecateMsg(types.V("0.2.0"), types.V("0.3.0"), "Use", lastRoute, "instead"),
			),
		)
		ctx.Redirect(types.Http.C300().PermanentRedirect(), lastRoute)
	})

	companionRouter := router.Group(lastRoute)
	{
		companionRouter.GET("/:id", controller.Companion.GetByID)
		companionRouter.GET("/all", controller.Companion.GetAll)

		companionRouter.POST("/", controller.Companion.Create)

		companionRouter.PUT("/:id", controller.Companion.UpdateByID)

		companionRouter.PATCH("/:id", controller.Companion.PatchByID)

		companionRouter.DELETE("/:id", controller.Companion.DeleteByID)
		//companionRouter.DELETE("/permanent-delete/:id/:confirm", controller.Student.ForceDeleteByID)

	}
}
