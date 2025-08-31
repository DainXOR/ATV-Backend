package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"

	"dainxor/atv/configs"
	"dainxor/atv/logger"
	"dainxor/atv/middleware"
	"dainxor/atv/routes"
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
	routes.MainRoutes(router)

	// Versioned API routes
	routes.StudentRoutes(router)
	routes.UniversityRoutes(router)
	routes.SpecialityRoutes(router)
	routes.CompanionRoutes(router)
	routes.SessionTypeRoutes(router)
	routes.SessionRoutes(router)

	// Api informative routes
	routes.TestRoutes(router) // Routes for testing purposes
	routes.InfoRoutes(router) // Routes for information about the API

	router.Run(configs.App.Address()) // listen and serve on 0.0.0.0:8080 (for windows ":8080")
}
