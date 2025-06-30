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
				logger.DeprecateMsg("0.0.3", "0.1.0", "Use /api/v1/student/ instead"),
			),
		)
	})

	studentRouter := router.Group("api/v0/student")
	{
		studentRouter.GET("/:id", controller.Student.GetByIDGorm)
		studentRouter.GET("/all", controller.Student.GetAllGorm)

		studentRouter.POST("/", controller.Student.CreateGorm)
	}
	studentRouter = router.Group("api/v1/student")
	{
		studentRouter.GET("/:id", controller.Student.GetByIDMongo)
		studentRouter.GET("/all", controller.Student.GetAllMongo)

		studentRouter.POST("/", controller.Student.CreateMongo)

		studentRouter.PUT("/:id", controller.Student.UpdateMongo)

		studentRouter.PATCH("/:id", controller.Student.PatchMongo)

		studentRouter.DELETE("/:id", controller.Student.DeleteByID)
		//studentRouter.DELETE("/permanent-delete/:id/:confirm", controller.Student.ForceDeleteByID)

	}
}
