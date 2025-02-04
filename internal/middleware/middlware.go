package middleware

import (
	"github.com/gin-gonic/gin"
)

func SetCommonHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// CORS
		c.Header("Access-Control-Allow-Origin", "*")
		// Future headers if needed
		c.Next()
	}
}
