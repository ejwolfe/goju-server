package server

import (
	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(context *gin.Context) {
		logRequest(context)
		context.Next()
		logResponse(context)
	}
}

func logRequest(c *gin.Context) {
	headers := c.Request.Header
	method := c.Request.Method
	url := c.Request.URL
	body := c.Request.Body
	serviceLogger.Info("Request", method, url, headers, body)
}

func logResponse(c *gin.Context) {
	headers := c.Request.Response.Header
	body := c.Request.Response.Body
	serviceLogger.Info("Response", headers, body)
}
