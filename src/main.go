package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"

	"dainxor/atv/configs"
	"dainxor/atv/controller"
	"dainxor/atv/logger"
	"dainxor/atv/middleware"
)

//var envErr = godotenv.Load()

func init() {
	logger.SetVersion(configs.App.ApiVersion())
	configs.DB.Use(configs.GetMongoAccessor()).Start()
	// configs.DB.Migrate(&models.StudentDBMongo{})
	logger.Info("Env configurations loaded")
	logger.Debug("Starting server")
}

func main() {
	defer configs.DB.Close()
	defer logger.Close()

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.Recovery()) // Middleware to recover from panics and logs a small trace
	router.Use(middleware.CORS())
	//router.Use(middleware.TokenMiddleware())

	// Root level routes
	controller.MainRoutes(router)

	// Versioned API routes
	controller.StudentsRoutes(router)
	controller.UniversitiesRoutes(router)
	controller.SpecialitiesRoutes(router)
	controller.CompanionsRoutes(router)
	controller.SessionTypesRoutes(router)
	controller.SessionsRoutes(router)
	controller.PrioritiesRoutes(router)
	controller.AlertsRoutes(router)

	// Api informative routes
	controller.TestRoutes(router) // Routes for testing purposes
	controller.InfoRoutes(router) // Routes for information about the API

	router.Run(configs.App.Address()) // listen and serve on 0.0.0.0:8080 (for windows ":8080")
}
