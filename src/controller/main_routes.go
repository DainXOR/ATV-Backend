package controller

import (
	"dainxor/atv/types"

	"github.com/gin-gonic/gin"
)

func MainRoutes(router *gin.Engine) {
	router.GET("/", func(c *gin.Context) {
		c.JSON(types.Http.C200().Ok(), gin.H{
			"message": "Go to /api/info/ to see available routes",
		})
	})
}
