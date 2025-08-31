package routes

import (
	"dainxor/atv/configs"
	"dainxor/atv/controller"
	"dainxor/atv/logger"
	"dainxor/atv/types"
	"fmt"

	"github.com/gin-gonic/gin"
)

func UniversityRoutes(router *gin.Engine) {
	rv := configs.App.RoutesVersion()
	beforeRoute := fmt.Sprintf("/api/v%d/university", rv-1)
	lastRoute := fmt.Sprintf("/api/v%d/university", rv)

	router.Group(beforeRoute).Any("", func(ctx *gin.Context) {
		ctx.JSON(types.Http.C300().PermanentRedirect(),
			types.EmptyResponse(
				logger.DeprecateMsg(types.V("0.2.0"), types.V("0.3.0"), "Use", lastRoute, "instead"),
			),
		)
		ctx.Redirect(types.Http.C300().PermanentRedirect(), lastRoute)
	})

	universityRouter := router.Group(lastRoute)
	{
		universityRouter.POST("/", controller.University.Create)

		universityRouter.GET("/:id", controller.University.GetByID)
		universityRouter.GET("/all", controller.University.GetAll)
	}
}
