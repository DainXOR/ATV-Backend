package routes

import (
	"dainxor/atv/controller"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine) {
	userRouter := router.Group("api/v0/user")
	{
		userRouter.GET("/:id", controller.User.GetByID)

		userRouter.POST("/", controller.User.Create)
	}
}
