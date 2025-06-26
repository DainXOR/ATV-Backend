package routes

import (
	"dainxor/atv/controller"

	"github.com/gin-gonic/gin"
)

func SessionRoutes(router *gin.Engine) {
	// Grouping the speciality routes under "api/v1/speciality"
	sessionRouter := router.Group("api/v1/session")
	{
		sessionRouter.POST("/", controller.Session.Create)

		sessionRouter.GET("/:id", controller.Session.GetByID)
		sessionRouter.GET("/student/:student_id", controller.Session.GetAllByStudentID)
		sessionRouter.GET("/all", controller.Session.GetAll)
	}
}
