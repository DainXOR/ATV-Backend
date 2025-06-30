package middleware

import (
	"dainxor/atv/logger"

	"github.com/gin-gonic/gin"
)

func TokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Debug("Header:", c.Request.Header)

		c.Next()
	}
}
