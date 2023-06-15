package server

import (
	"bytes"
	"encoding/json"
	"io/ioutil"

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

func readEntriesFile() []Entry {
	data, err := ioutil.ReadFile("entries.json")
	if err != nil {
		serviceLogger.Error("Error reading file", err)
	}

	var fileEntries []Entry
	jsonParseErr := json.Unmarshal(data, &fileEntries)

	if jsonParseErr != nil {
		serviceLogger.Error("Error parsing json", jsonParseErr)
	}

	return fileEntries
}

func writeEntriesFile(entries []Entry) {
	json, err := json.Marshal(entries)
	if err != nil {
		serviceLogger.Error("Error encoding json", err)
	}

	writeErr := ioutil.WriteFile("entries.json", json, 0644)

	if writeErr != nil {
		serviceLogger.Error("Error writing json file", writeErr)
	}
}
