package routes

import (
	"dainxor/atv/controller"

	"github.com/gin-gonic/gin"
)

func CompanionRoutes(router *gin.Engine) {
	// Grouping the companion routes under "api/v#/companion"
	// This allows for better organization and versioning of the API
	// Grouping can also be done inside other groups
	companionRouter := router.Group("api/v1/companion")
	{
		companionRouter.GET("/:id", controller.Companion.GetByID)
		companionRouter.GET("/all", controller.Companion.GetAll)

		companionRouter.POST("/", controller.Companion.Create)

		companionRouter.PUT("/:id", controller.Companion.UpdateByID)

		companionRouter.PATCH("/:id", controller.Companion.PatchByID)

		companionRouter.DELETE("/:id", controller.Companion.DeleteByID)
		//companionRouter.DELETE("/permanent-delete/:id/:confirm", controller.Student.ForceDeleteByID)

	}
}
