package middleware

import (
	"github.com/EtienneBerube/cat-scribers/pkg/auth"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

/* Auth is a middleware to check if requests is from an authenticated user. The request will be aborted if the request
is not authenticated by a user's token
*/
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		bearToken := c.GetHeader("Authorization")

		strArr := strings.Split(bearToken, " ")
		if len(strArr) == 2 {
			if ok, err := auth.ValidateToken(strArr[1]); !ok || err != nil {
				c.AbortWithStatus(http.StatusUnauthorized)
			}

			id, err := auth.ExtractUserId(strArr[1])
			if err != nil {
				c.AbortWithStatus(http.StatusUnauthorized)
			}

			c.Set("user_id", id)
			c.Next()
		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}
