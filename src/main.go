package main

import (
	"dainxor/atv/configs"
	"dainxor/atv/logger"
	"dainxor/atv/middleware"
	"dainxor/atv/models"
	"dainxor/atv/routes"

	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file: " + err.Error())
	}

	logger.EnvInit()
	configs.DB.EnvInit()
	configs.DB.Migrate(&models.UserDBGorm{})
	logger.Info("Env configurations loaded")
	logger.Info("Starting server")

}

// address returns the server address from the environment variable
func address() string {
	address, exist := os.LookupEnv("SERVER_ADDRESS")
	if !exist {
		logger.Warning("SERVER_ADDRESS not found, using default")
		address = ":8080"
	}

	return address
}

func main() {
	router := gin.Default()
	router.Use(middleware.CORSMiddleware())

	routes.MainRoutes(router) // Main route for the API
	routes.InfoRoutes(router) // Routes for information about the API
	routes.TestRoutes(router) // Routes for testing purposes

	routes.UserRoutes(router) // Routes for user management

	router.Run(address()) // listen and serve on 0.0.0.0:8080 (for windows ":8080")
}
