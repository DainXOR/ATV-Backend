package routes

import (
	"dainxor/atv/configs"
	"dainxor/atv/controller"
	"dainxor/atv/logger"
	"dainxor/atv/types"
	"fmt"

	"github.com/gin-gonic/gin"
)

func SpecialityRoutes(router *gin.Engine) {
	rv := configs.App.RoutesVersion()
	beforeRoute := fmt.Sprintf("/api/v%d/speciality", rv-1)
	lastRoute := fmt.Sprintf("/api/v%d/speciality", rv)

	router.Group(beforeRoute).Any("", func(ctx *gin.Context) {
		ctx.JSON(types.Http.C300().PermanentRedirect(),
			types.EmptyResponse(
				logger.DeprecateMsg(types.V("0.2.0"), types.V("0.3.0"), "Use", lastRoute, "instead"),
			),
		)
		ctx.Redirect(types.Http.C300().PermanentRedirect(), lastRoute)
	})

	specialityRouter := router.Group(lastRoute)
	{
		specialityRouter.POST("/", controller.Speciality.Create)

		specialityRouter.GET("/:id", controller.Speciality.GetByID)
		specialityRouter.GET("/all", controller.Speciality.GetAll)
	}
}
