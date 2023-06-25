package server

import (
	"bytes"
	"github.com/gin-gonic/gin"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func Logger() gin.HandlerFunc {
	return func(context *gin.Context) {
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: context.Writer}
		context.Writer = blw
		logRequest(context)
		context.Next()
		logResponse(context, blw)
	}
}

func logRequest(c *gin.Context) {
	headers := c.Request.Header
	method := c.Request.Method
	url := c.Request.URL
	body := c.Request.Body
	serviceLogger.Info("Request", method, url, headers, body)
}

func logResponse(c *gin.Context, writer *bodyLogWriter) {
	headers := c.Writer.Header()
	body := writer.body.String()
	serviceLogger.Info("Response", headers, body)
}

func removeElementByIndex[T any](slice []T, index int) []T {
	return append(slice[:index], slice[index+1:]...)
}

func createResponseErrorMessage(message string) map[string]any {
	return gin.H{"message": message}
}
