package middleware

import (
	"dainxor/atv/logger"
	"dainxor/atv/types"
	"dainxor/atv/utils"
	"fmt"

	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("Recovered from panic:", err)

				f0, l0 := utils.CallOrigin(4)
				f1, l1 := utils.CallOrigin(5)
				f2, l2 := utils.CallOrigin(6)
				f3, l3 := utils.CallOrigin(7)
				f4, l4 := utils.CallOrigin(8)

				logger.Error(fmt.Sprintf(
					"Error originated at: %s:%d > %s:%d > %s:%d > %s:%d > %s:%d",
					f4, l4,
					f3, l3,
					f2, l2,
					f1, l1,
					f0, l0,
				))

				c.AbortWithStatusJSON(types.Http.C500().InternalServerError(),
					types.EmptyResponse(
						"An unexpected error occurred. Please try again later.",
						"Check the server logs for more details.",
					))
				return
			}

		}()
		c.Next()
	}
}
