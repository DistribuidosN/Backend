package handlers

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func RequestTrace() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		log.Printf("[REST] --> %s %s", c.Request.Method, c.Request.URL.Path)
		c.Next()
		log.Printf("[REST] <-- %s %s status=%d duration=%s", c.Request.Method, c.Request.URL.Path, c.Writer.Status(), time.Since(start).Round(time.Millisecond))
	}
}
