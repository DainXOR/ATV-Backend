package main

import (
	"dainxor/atv/middleware"
	"dainxor/atv/routes"

	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Hello, World!")
	router := gin.Default()
	router.Use(middleware.CORSMiddleware())

	routes.MainRoutes(router)
	routes.InfoRoutes(router)
	routes.TestRoutes(router)

	router.Run("localhost:8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
