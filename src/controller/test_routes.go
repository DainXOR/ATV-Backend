package controller

import (
	"dainxor/atv/service"

	"github.com/gin-gonic/gin"
)

func TestRoutes(router *gin.Engine) {
	testRouter := router.Group("api/test")
	{
		testRouter.GET("/get", service.Test.Get)
		testRouter.POST("/post", service.Test.Post)
		testRouter.PUT("/put", service.Test.Put)
		testRouter.PATCH("/patch", service.Test.Patch)
		testRouter.DELETE("/del", service.Test.Delete)
	}
}
