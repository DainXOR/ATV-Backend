package routes

import (
	"dainxor/atv/controller"
	"dainxor/atv/logger"
	"dainxor/atv/types"

	"github.com/gin-gonic/gin"
)

func StudentRoutes(router *gin.Engine) {
	// Grouping the student routes under "api/v#/student"
	// This allows for better organization and versioning of the API
	// Grouping can also be done inside other groups
	router.Group("api/v0/user").Any("", func(ctx *gin.Context) {
		ctx.JSON(types.Http.C300().MovedPermanently(),
			types.EmptyResponse(
				logger.DeprecateMsg(types.V("0.0.3"), types.V("0.1.0"), "Use /api/v1/student/ instead"),
			),
		)
	})

	router.Group("api/v0/student").Any("", func(ctx *gin.Context) {
		ctx.JSON(types.Http.C300().MovedPermanently(),
			types.EmptyResponse(
				logger.DeprecateMsg(types.V("0.0.3"), types.V("0.1.0"), "Use /api/v1/student/ instead"),
			),
		)
	})
	studentRouter := router.Group("api/v1/student")
	{
		studentRouter.GET("/:id", controller.Student.GetByID)
		studentRouter.GET("/all", controller.Student.GetAll)

		studentRouter.POST("/", controller.Student.Create)

		studentRouter.PUT("/:id", controller.Student.UpdatebyID)

		studentRouter.PATCH("/:id", controller.Student.PatchByID)

		studentRouter.DELETE("/:id", controller.Student.DeleteByID)
		//studentRouter.DELETE("/permanent-delete/:id/:confirm", controller.Student.ForceDeleteByID)

	}
}
