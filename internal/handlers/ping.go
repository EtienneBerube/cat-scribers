package handlers

import "github.com/gin-gonic/gin"

// Ping responds to a ping request. Useful to obtain the health of the system
func Ping(c *gin.Context) {
	c.String(200, "pong")
	return
}
