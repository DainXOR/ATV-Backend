package controller

import (
	"dainxor/atv/configs"
	"dainxor/atv/models"
	"dainxor/atv/service"
	"dainxor/atv/types"
	"fmt"

	"github.com/gin-gonic/gin"
)

func FormsRoutes(router *gin.Engine) {
	rv := configs.App.RoutesVersion()
	// beforeRoute := fmt.Sprintf("/api/v%d/forms", rv-1)
	lastRoute := fmt.Sprintf("/api/v%d/forms", rv)

	//companionRouterOld := router.Group(beforeRoute)
	//{ }
	companionRouter := router.Group(lastRoute)
	{
		companionRouter.GET("/:id/questions", func(c *gin.Context) {
			// TO DO: implement pagination, filtering, etc.
			c.JSON(200, types.Response(
				gin.H{
					"id":   "97sdha7hd983h9h98hd9283hd981h3d",
					"name": "Sample Form Question",
					"questions": []models.FormQuestionResponse{
						{
							ID:             "q1",
							Name:           "Question 1",
							Question:       "What is your favorite color?",
							Options:        []string{"Red", "Blue", "Green", "Yellow"},
							IDQuestionType: "type1",
							CreatedAt:      models.Time.Now(),
							UpdatedAt:      models.Time.Now(),
						},
						{
							ID:             "q2",
							Name:           "Question 2",
							Question:       "What is your preferred mode of transport?",
							Options:        []string{"Car", "Bike", "Public Transport", "Walking"},
							IDQuestionType: "type2",
							CreatedAt:      models.Time.Now(),
							UpdatedAt:      models.Time.Now(),
						},
						{
							ID:             "q3",
							Name:           "Question 3",
							Question:       "Which cuisine do you like the most?",
							Options:        []string{"Italian", "Chinese", "Mexican", "Indian"},
							IDQuestionType: "type3",
							CreatedAt:      models.Time.Now(),
							UpdatedAt:      models.Time.Now(),
						},
					},
				},
				"",
			))
		})

		companionRouter.POST("/:id/questions", func(c *gin.Context) {
			// TO DO: implement validation, authorization, etc.
			service.FormQuestions.Create(c)
		})
	}
}
