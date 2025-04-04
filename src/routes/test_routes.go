package routes

import (
	"dainxor/atv/test"

	"github.com/gin-gonic/gin"
)

func TestRoutes(router *gin.Engine) {
	testRouter := router.Group("api/test")
	{
		testRouter.GET("/get", test.Get)
		testRouter.POST("/post", test.Post)
		testRouter.PUT("/put", test.Put)
		testRouter.PATCH("/patch", test.Patch)
		testRouter.DELETE("/del", test.Delete)
	}
}
