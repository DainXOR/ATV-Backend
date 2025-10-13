package controller

import (
	"dainxor/atv/configs"
	"dainxor/atv/service"
	"fmt"

	"github.com/gin-gonic/gin"
)

func SessionsRoutes(router *gin.Engine) {
	rv := configs.App.RoutesVersion()
	beforeRoute := fmt.Sprintf("/api/v%d/session", rv-1)
	lastRoute := fmt.Sprintf("/api/v%d/sessions", rv)

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
		sessionRouter.POST("/", service.Session.Create)

		sessionRouter.GET("/:id", service.Session.GetByID)
		sessionRouter.GET("/all", service.Session.GetAll)
		sessionRouter.GET("/student/:student_id", service.Session.GetAllByStudentID)

		sessionRouter.PUT("/:id", service.Session.UpdateByID)

		sessionRouter.PATCH("/:id", service.Session.PatchByID)

		sessionRouter.DELETE("/:id", service.Session.DeleteByID)
	}

	sessionRouter = router.Group(lastRoute)
	{
		sessionRouter.POST("/", service.Session.Create)

		sessionRouter.GET("/:id", service.Session.GetByID)
		sessionRouter.GET("/all", service.Session.GetAll)
		sessionRouter.GET("/student/:student_id", service.Session.GetAllByStudentID)

		sessionRouter.PUT("/:id", service.Session.UpdateByID)

		sessionRouter.PATCH("/:id", service.Session.PatchByID)

		sessionRouter.DELETE("/:id", service.Session.DeleteByID)
	}
}
