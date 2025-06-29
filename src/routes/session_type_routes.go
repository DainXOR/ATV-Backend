package routes

import (
	"dainxor/atv/controller"

	"github.com/gin-gonic/gin"
)

func SessionTypeRoutes(router *gin.Engine) {
	// Grouping the session type routes under "api/v1/session-type"
	sessionTypeRouter := router.Group("api/v1/session-type")
	{
		sessionTypeRouter.POST("/", controller.SessionType.Create)

		sessionTypeRouter.GET("/:id", controller.SessionType.GetByID)
		sessionTypeRouter.GET("/all", controller.SessionType.GetAll)
	}
}
