package controller

import (
	"dainxor/atv/configs"
	"dainxor/atv/logger"
	"dainxor/atv/service"
	"dainxor/atv/types"
	"fmt"

	"github.com/gin-gonic/gin"
)

func StudentsRoutes(router *gin.Engine) {
	rv := configs.App.RoutesVersion()
	beforeRoute := fmt.Sprintf("/api/v%d/student", rv-1)
	lastRoute := fmt.Sprintf("/api/v%d/students", rv)

	router.Group(beforeRoute).Any("", func(ctx *gin.Context) {
		ctx.JSON(types.Http.C300().PermanentRedirect(),
			types.EmptyResponse(
				logger.DeprecateMsg(types.V("0.2.0"), types.V("0.3.0"), "Use", lastRoute, "instead"),
			),
		)
		ctx.Redirect(types.Http.C300().PermanentRedirect(), lastRoute)
	})

	//studentRouter := router.Group(beforeRoute)
	//{
	//	studentRouter.GET("/:id", controller.Student.GetByID)
	//	studentRouter.GET("/all", controller.Student.GetAll)
	//
	//	studentRouter.POST("/", controller.Student.Create)
	//
	//	studentRouter.PUT("/:id", controller.Student.UpdatebyID)
	//
	//	studentRouter.PATCH("/:id", controller.Student.PatchByID)
	//
	//	studentRouter.DELETE("/:id", controller.Student.DeleteByID)
	//	//studentRouter.DELETE("/permanent-delete/:id/:confirm", controller.Student.ForceDeleteByID)
	//
	//}

	studentRouter := router.Group(lastRoute)
	{
		studentRouter.GET("/:id", service.Student.GetByID)
		studentRouter.GET("/all", service.Student.GetAll)

		studentRouter.POST("/", service.Student.Create)

		studentRouter.PUT("/:id", service.Student.UpdateByID)

		studentRouter.PATCH("/:id", service.Student.PatchByID)

		studentRouter.DELETE("/:id", service.Student.DeleteByID)
		//studentRouter.DELETE("/permanent-delete/:id/:confirm", controller.Student.ForceDeleteByID)

	}
}
