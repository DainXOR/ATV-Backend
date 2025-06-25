package routes

import (
	"dainxor/atv/controller"

	"github.com/gin-gonic/gin"
)

func SpecialityRoutes(router *gin.Engine) {
	// Grouping the speciality routes under "api/v1/speciality"
	specialityRouter := router.Group("api/v1/speciality")
	{
		specialityRouter.POST("/", controller.Speciality.Create)

		specialityRouter.GET("/:id", controller.Speciality.GetByID)
		specialityRouter.GET("/all", controller.Speciality.GetAll)
	}
}
