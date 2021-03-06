package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

// WithLogging is a custom formatting for the servers logging system
func WithLogging(param gin.LogFormatterParams) string {

	// your custom format
	return fmt.Sprintf("%s - [%s] --%s %s %s %d %s \"%s\" %s\n",
		param.ClientIP,
		param.TimeStamp.Format(time.RFC1123),
		param.Method,
		param.Path,
		param.Request.Proto,
		param.StatusCode,
		param.Latency,
		param.Request.UserAgent(),
		param.ErrorMessage,
	)
}
