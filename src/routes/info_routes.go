package routes

import (
	"dainxor/atv/configs"
	"dainxor/atv/logger"
	"dainxor/atv/utils"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func InfoRoutes(router *gin.Engine) {
	availableRoutes := BuildRoutesInfo(router)
	availableRoutes["info"] = gin.H{
		"root":          "/api/info/",
		"ping":          "/api/info/ping",
		"api version":   "/api/info/api-version",
		"route version": "/api/info/route-version",
	}

	infoRouter := router.Group("api/info")
	{
		infoRouter.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Available routes",
				"routes":  availableRoutes,
			})
		})
		infoRouter.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})
		infoRouter.GET("/api-version", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"version": configs.App.ApiVersion(),
			})
		})
		infoRouter.GET("/route-version", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"version": configs.App.RoutesVersion(),
			})
		})
	}
}

var omitMethods = []string{
	"connect",
	"options",
	"trace",
	"head",
}

func BuildRoutesInfo(router *gin.Engine) gin.H {
	result := gin.H{}

	for _, routeInfo := range router.Routes() {
		path := routeInfo.Path // e.g. /api/v1/companion/:id
		method := routeInfo.Method

		if utils.Any(omitMethods, func(m string) bool { return m == method }) || path == "/" {
			continue
		}
		logger.Debug("Route:", method, path)

		pathParts := strings.Split(strings.Trim(path, "/"), "/")[1:]
		logger.Debug("Route parts:", pathParts)

		current := result
		for len(pathParts) > 1 {
			logger.Debug("Entering group:", pathParts[0])

			if _, ok := current[pathParts[0]]; !ok {
				current[pathParts[0]] = gin.H{}
			}

			current = current[pathParts[0]].(gin.H)
			pathParts = pathParts[1:]
		}

		opName, ok := operationName(method, pathParts[0])
		if !ok {
			continue
		}
		current[opName] = path
		logger.Debug("Result:", result)

	}

	return result
}

func operationName(method, rest string) (string, bool) {
	operation := strings.ToLower(method)

	if strings.HasPrefix(operation, rest) {
		return operation, true

	} else if rest == "all" {
		return fmt.Sprintf("%s all", operation), true

	} else if pathVariable, ok := strings.CutPrefix(rest, ":"); ok {

		if pathVariable == "id" {
			return fmt.Sprintf("%s by id", operation), true
		}
		return fmt.Sprintf("%s by %s", operation, strings.ReplaceAll(pathVariable, "_", " ")), true

	} else {
		return fmt.Sprintf("%s %s", operation, rest), false
	}
}
