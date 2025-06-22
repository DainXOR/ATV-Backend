package routes

import (
	"dainxor/atv/controller"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine) {
	// Grouping the user routes under "api/v0/user"
	// This allows for better organization and versioning of the API
	// Grouping can also be done inside other groups
	userRouter := router.Group("api/v0/student")
	{
		userRouter.GET("/:id", controller.Student.GetByIDGorm)
		userRouter.GET("/all", controller.Student.GetAllGorm)

		userRouter.POST("/", controller.Student.CreateGorm)
	}
	userRouter = router.Group("api/v1/student")
	{
		userRouter.GET("/:id", controller.Student.GetByIDMongo)
		userRouter.GET("/all", controller.Student.GetAllMongo)

		userRouter.POST("/", controller.Student.CreateMongo)

		//userRouter.PUT("/:id", controller.Student.UpdateMongo)

		//userRouter.PATCH("/:id", controller.Student.PatchMongo)

		//userRouter.DELETE("/:id", controller.Student.Delete)

	}
}
