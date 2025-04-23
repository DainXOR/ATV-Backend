package main

import (
	"dainxor/atv/configs"
	"dainxor/atv/logger"
	"dainxor/atv/middleware"
	"dainxor/atv/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var address = "localhost:8080"

func init() {
	logger.Init()

	err := godotenv.Load()
	if err != nil {
		logger.Error("Error loading .env file")
	}

	logger.EnvInit()
	configs.DB.EnvInit()
	logger.Info("Env configurations loaded")
	logger.Info("Starting server")
}

func main() {
	router := gin.Default()
	router.Use(middleware.CORSMiddleware())

	routes.MainRoutes(router)
	routes.InfoRoutes(router)
	routes.TestRoutes(router)

	routes.UserRoutes(router)

	router.Run("localhost:8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
