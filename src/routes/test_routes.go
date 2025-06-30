package routes

import (
	"dainxor/atv/controller"

	"github.com/gin-gonic/gin"
)

func TestRoutes(router *gin.Engine) {
	testRouter := router.Group("api/test")
	{
		testRouter.GET("/get", controller.Test.Get)
		testRouter.POST("/post", controller.Test.Post)
		testRouter.PUT("/put", controller.Test.Put)
		testRouter.PATCH("/patch", controller.Test.Patch)
		testRouter.DELETE("/del", controller.Test.Delete)
	}
}
