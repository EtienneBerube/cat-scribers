package middleware

import (
	"github.com/EtienneBerube/only-cats/pkg/auth"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		bearToken := c.GetHeader("Authorization")

		strArr := strings.Split(bearToken, " ")
		if len(strArr) == 2 {
			if ok, err := auth.ValidateToken(strArr[1]); !ok || err != nil{
				c.AbortWithStatus(http.StatusUnauthorized)
			}

			id, err := auth.ExtractUserId(strArr[1])
			if err != nil {
				c.AbortWithStatus(http.StatusUnauthorized)
			}

			c.Set("user_id", id)
			c.Next()
		}

		c.AbortWithStatus(http.StatusUnauthorized)
	}
}

