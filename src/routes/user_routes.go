package routes

import (
	"dainxor/atv/controller"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine) {
	// Grouping the user routes under "api/v0/user"
	// This allows for better organization and versioning of the API
	// Grouping can also be done inside other groups
	userRouter := router.Group("api/v0/user")
	{
		userRouter.GET("/:id", controller.User.GetByID)

		userRouter.POST("/", controller.User.Create)
	}
}
