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
	_ = gin.H{
		"info": gin.H{
			"root":          "/api/info/",
			"ping":          "/api/info/ping",
			"api version":   "/api/info/api-version",
			"route version": "/api/info/route-version",
		},

		"test": gin.H{
			"get":    "/api/test/get",
			"post":   "/api/test/post",
			"put":    "/api/test/put",
			"patch":  "/api/test/patch",
			"delete": "/api/test/del",
		},

		"v1": gin.H{
			"student": gin.H{
				"post":         "/api/v1/student/",
				"get by id":    "/api/v1/student/:id",
				"get all":      "/api/v1/student/all",
				"put":          "/api/v1/student/:id",
				"patch":        "/api/v1/student/:id",
				"delete by id": "/api/v1/student/:id",
				//"force delete by id": "/api/v1/student/permanent-delete/:id/:confirm",
			},
			"university": gin.H{
				"post":      "/api/v1/university/",
				"get by id": "/api/v1/university/:id",
				"get all":   "/api/v1/university/all",
			},
			"speciality": gin.H{
				"post":      "/api/v1/speciality/",
				"get by id": "/api/v1/speciality/:id",
				"get all":   "/api/v1/speciality/all",
			},
			"session type": gin.H{
				"post":      "/api/v1/session-type/",
				"get by id": "/api/v1/session-type/:id",
				"get all":   "/api/v1/session-type/all",
			},
			"session": gin.H{
				"post":      "/api/v1/session/",
				"get by id": "/api/v1/session/:id",
				"get all":   "/api/v1/session/all",
			},
			"companion": gin.H{
				"post":         "/api/v1/companion/",
				"get by id":    "/api/v1/companion/:id",
				"get all":      "/api/v1/companion/all",
				"put":          "/api/v1/companion/:id",
				"patch":        "/api/v1/companion/:id",
				"delete by id": "/api/v1/companion/:id",
			},
		},
	}

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
