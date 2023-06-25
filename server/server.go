package server

import (
	"io"
	"os"

	"github.com/ejwolfe/goju-server/logger"
	"github.com/gin-gonic/gin"
)

var serviceLogger = logger.CreateLogger()

func CreateServer() {
	serverFile, _ := os.Create(serverLogFile)
	gin.DefaultWriter = io.MultiWriter(serverFile)
	router := gin.Default()
	router.Use(Logger())

	loadEntries()

	v1 := router.Group(v1Path)
	{
		v1.GET(entriesPath, getEntries)
		v1.GET(entriesIdPath, getEntryByID)
		v1.POST(entriesPath, addEntry)
		v1.PUT(entriesIdPath, updateEntry)
		v1.DELETE(entriesIdPath, removeEntry)
	}

	router.Run(serverAddress)
}
