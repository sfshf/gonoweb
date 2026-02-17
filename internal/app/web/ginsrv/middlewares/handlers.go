package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// NoMethod middleware
func NoMethod() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, nil)
		c.Abort()
	}
}

// NoRoute middleware
func NoRoute() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotFound, nil)
		c.Abort()
	}
}
