package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// loggingMiddleware writes its message just before HTTP request.
func loggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Info(fmt.Sprintf("%s %s %s (%s, %s)",
			c.Request.Proto, c.Request.Method, c.Request.RequestURI,
			c.ClientIP(), c.Request.UserAgent()))

		c.Next()
	}
}

