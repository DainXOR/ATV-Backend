package routes

import (
	"dainxor/atv/configs"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InfoRoutes(router *gin.Engine) {
	//routes := gin.H{
	//	"info":               "/api/info/",
	//	"info ping":          "/api/info/ping",
	//	"info api version":   "/api/info/api-version",
	//	"info route version": "/api/info/route-version",
	//
	//	"info test get":    "/api/test/get",
	//	"info test post":   "/api/test/post",
	//	"info test put":    "/api/test/put",
	//	"info test patch":  "/api/test/patch",
	//	"info test delete": "/api/test/del",
	//
	//	//"register email":        "/api/v0/user/register/:email",
	//	//"user create":           "/api/v0/user/",
	//	//"user get all":          "/api/v0/user/all/",
	//	//"user get by id":        "/api/v0/user/id/:id",
	//	//"user get by status id": "/api/v0/user/id-status/:id",
	//	//"user update by id":     "/api/v0/user/id/:id",
	//	//"user delete by id":     "/api/v0/user/id/:id",
	//}
	newRoutes := gin.H{
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
				"post":                  "/api/v1/session/",
				"get by id":             "/api/v1/session/:id",
				"get all by student id": "/api/v1/session/student/:student_id",
				"get all":               "/api/v1/session/all",
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

	infoRouter := router.Group("api/info")
	{
		infoRouter.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Available routes",
				"routes":  newRoutes,
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
