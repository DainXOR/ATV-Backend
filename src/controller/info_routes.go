package controller

import (
	"dainxor/atv/configs"
	"dainxor/atv/logger"
	"dainxor/atv/types"
	"dainxor/atv/utils"
	"fmt"
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
			c.JSON(types.Http.C200().Ok(), gin.H{
				"message": "Available routes",
				"routes":  availableRoutes,
			})
		})
		infoRouter.GET("/ping", func(c *gin.Context) {
			c.JSON(types.Http.C200().Ok(), gin.H{
				"message": "pong",
			})
		})
		infoRouter.GET("/api-version", func(c *gin.Context) {
			c.JSON(types.Http.C200().Ok(), gin.H{
				"version": configs.App.ApiVersion(),
			})
		})
		infoRouter.GET("/route-version", func(c *gin.Context) {
			c.JSON(types.Http.C200().Ok(), gin.H{
				"version": configs.App.RoutesVersion(),
			})
		})
	}
}

var includeMethods = []string{
	"POST",
	"GET",
	"put",
	"patch",
	"delete",
}

func BuildRoutesInfo(router *gin.Engine) gin.H {
	result := gin.H{}

	for _, routeInfo := range router.Routes() {
		path := routeInfo.Path // e.g. /api/v1/companion/:id
		method := routeInfo.Method

		if !utils.Any(includeMethods, func(m string) bool { return m == method }) {
			continue
		}

		pathParts := strings.Split(path, "/")[1:]
		current := result

		if pathLen := len(pathParts) - 1; pathLen > 1 && pathParts[pathLen] == "" {
			pathParts[pathLen] = pathParts[pathLen-1]
		}

		for len(pathParts) > 1 {
			if _, ok := current[pathParts[0]]; !ok {
				current[pathParts[0]] = gin.H{}
			}

			current = current[pathParts[0]].(gin.H)
			pathParts = pathParts[1:]
		}

		var opName string
		var ok bool

		if len(pathParts) != 0 {
			opName, ok = operationName(method, pathParts[0])
		} else {
			current["root"] = gin.H{}
			current = current["root"].(gin.H)
			opName = strings.ToLower(method)

			ok = true
		}

		if ok {
			current[opName] = path
		}
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

	} else if operation == "post" {
		return fmt.Sprintf("%s %s", operation, rest), true

	} else {
		logger.Warningf("No rest: %s, %s", method, rest)
		return fmt.Sprintf("%s %s", operation, rest), false
	}
}
