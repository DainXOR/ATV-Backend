package middleware

import (
	"dainxor/atv/logger"
	"dainxor/atv/models"
	"dainxor/atv/types"

	"github.com/gin-gonic/gin"
)

func RecoverMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("Recovered from panic:", err)

				c.AbortWithStatusJSON(types.Http.C500().InternalServerError(),
					models.Response(
						gin.H{},
						"An unexpected error occurred. Please try again later.",
						"Check the server logs for more details.",
					))
				return
			}

		}()
		c.Next()
	}
}
