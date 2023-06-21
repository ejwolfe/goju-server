package server

import (
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/ejwolfe/goju-server/logger"
)

type EntryType int

const (
	Task EntryType = iota
	Completed
	Migrated
	Scheduled
	Cancelled
	Note
	Event
	Priority
)

type Entry struct {
	ID      string    `json:"id"`
	Type    EntryType `json:"type"`
	Message string    `json:"message"`
}

var entries []Entry

var serviceLogger = logger.CreateLogger()

func CreateServer() {
	serverFile, _ := os.Create(serverLogFile)
	gin.DefaultWriter = io.MultiWriter(serverFile)
	router := gin.Default()
	router.Use(Logger())

	entries = readEntriesFile()

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

func getEntries(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, entries)
}

func getEntryByID(c *gin.Context) {
	id := c.Param(idParam)

	for _, entry := range entries {
		if entry.ID == id {
			c.IndentedJSON(http.StatusOK, entry)
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, createResponseErrorMessage(error001))
}

func addEntry(c *gin.Context) {
	var newEntry Entry

	if err := c.BindJSON(&newEntry); err != nil {
		serviceLogger.Error(error003, err)
		c.IndentedJSON(http.StatusBadRequest, createResponseErrorMessage(error002))
		return
	}

	newEntry.ID = uuid.NewString()
	entries = append(entries, newEntry)
	location := v1Path + entriesPath + newEntry.ID
	c.Header(locationHeader, location)
	c.IndentedJSON(http.StatusCreated, newEntry)
	writeEntriesFile(entries)
}

func removeEntry(c *gin.Context) {
	id := c.Param(idParam)

	for index, entry := range entries {
		if entry.ID == id {
			entries = removeElementByIndex(entries, index)
			c.Status(http.StatusNoContent)
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, createResponseErrorMessage(error001))
}

func updateEntry(c *gin.Context) {
	id := c.Param(idParam)

	for index, entry := range entries {
		if entry.ID == id {
			var newEntry Entry

			if err := c.BindJSON(&newEntry); err != nil {
				serviceLogger.Error(error003, err)
				c.IndentedJSON(http.StatusBadRequest, createResponseErrorMessage(error001))
				return
			}
			entries[index].Type = newEntry.Type
			entries[index].Message = newEntry.Message
			c.Status(http.StatusNoContent)
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, createResponseErrorMessage(error001))
}
