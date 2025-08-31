package routes

import (
	"dainxor/atv/configs"
	"dainxor/atv/controller"
	"fmt"

	"github.com/gin-gonic/gin"
)

func SessionRoutes(router *gin.Engine) {
	rv := configs.App.RoutesVersion()
	beforeRoute := fmt.Sprintf("/api/v%d/session", rv-1)
	lastRoute := fmt.Sprintf("/api/v%d/session", rv)

	//router.Group(beforeRoute).Any("", func(ctx *gin.Context) {
	//	ctx.JSON(types.Http.C300().PermanentRedirect(),
	//		types.EmptyResponse(
	//			logger.DeprecateMsg(types.V("0.3.0"), types.V("0.4.0"), "Use", lastRoute, "instead"),
	//		),
	//	)
	//	ctx.Redirect(types.Http.C300().PermanentRedirect(), lastRoute)
	//})

	sessionRouter := router.Group(beforeRoute)
	{
		sessionRouter.POST("/", controller.Session.Create)

		sessionRouter.GET("/:id", controller.Session.GetByID)
		sessionRouter.GET("/all", controller.Session.GetAll)
		sessionRouter.GET("/student/:student_id", controller.Session.GetAllByStudentID)

		sessionRouter.PUT("/:id", controller.Session.UpdateByID)

		sessionRouter.PATCH("/:id", controller.Session.PatchByID)

		sessionRouter.DELETE("/:id", controller.Session.DeleteByID)
	}

	sessionRouter = router.Group(lastRoute)
	{
		sessionRouter.POST("/", controller.Session.Create)

		sessionRouter.GET("/:id", controller.Session.GetByID)
		sessionRouter.GET("/all", controller.Session.GetAll)
		sessionRouter.GET("/student/:student_id", controller.Session.GetAllByStudentID)

		sessionRouter.PUT("/:id", controller.Session.UpdateByID)

		sessionRouter.PATCH("/:id", controller.Session.PatchByID)

		sessionRouter.DELETE("/:id", controller.Session.DeleteByID)
	}
}
