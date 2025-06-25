package routes

import (
	"dainxor/atv/controller"

	"github.com/gin-gonic/gin"
)

func UniversityRoutes(router *gin.Engine) {
	// Grouping the university routes under "api/v1/university"
	universityRouter := router.Group("api/v1/university")
	{
		universityRouter.POST("/", controller.University.Create)

		universityRouter.GET("/:id", controller.University.GetByID)
		universityRouter.GET("/all", controller.University.GetAll)
	}
}
