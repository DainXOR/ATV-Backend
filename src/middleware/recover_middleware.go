package middleware

import (
	"dainxor/atv/logger"
	"dainxor/atv/types"
	"dainxor/atv/utils"
	"fmt"

	"github.com/gin-gonic/gin"
)

func RecoverMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("Recovered from panic:", err)

				origin1 := utils.CallOrigin(5)
				origin2 := utils.CallOrigin(6)
				origin3 := utils.CallOrigin(7)

				logger.Error(fmt.Sprintf("Error originated at: %s > %s > %s", origin3, origin2, origin1))

				c.AbortWithStatusJSON(types.Http.C500().InternalServerError(),
					types.Response(
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
