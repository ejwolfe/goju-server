package server

import (
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/ejwolfe/goju-server/logger"
)

const CONTEXT = "/entries"

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
	serverFile, _ := os.Create("server.log")
	gin.DefaultWriter = io.MultiWriter(serverFile)
	router := gin.Default()
	router.Use(Logger())

	entries = readEntriesFile()

	v1 := router.Group("/v1")
	{
		v1.GET(CONTEXT, getEntries)
		v1.GET(CONTEXT+"/:id", getEntryByID)
		v1.POST(CONTEXT, addEntry)
		v1.DELETE(CONTEXT+"/:id", removeEntry)
		v1.PUT(CONTEXT+"/:id", updateEntry)
	}

	router.Run("localhost:8080")
}

func getEntries(c *gin.Context) {
	serviceLogger.Info("Request Headers", c.Request.Header)
	c.IndentedJSON(http.StatusOK, entries)
	serviceLogger.Info("Response Body", http.StatusOK, entries)
}

func getEntryByID(c *gin.Context) {
	id := c.Param("id")

	for _, entry := range entries {
		if entry.ID == id {
			c.IndentedJSON(http.StatusOK, entry)
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "entry not found"})
}

func addEntry(c *gin.Context) {
	var newEntry Entry

	if err := c.BindJSON(&newEntry); err != nil {
		serviceLogger.Error("Failed to parse request", err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}

	newEntry.ID = uuid.NewString()
	entries = append(entries, newEntry)
	locationHeader := CONTEXT + "/" + newEntry.ID
	c.Header("Location", locationHeader)
	c.IndentedJSON(http.StatusCreated, newEntry)
	writeEntriesFile(entries)
}

func removeEntry(c *gin.Context) {
	id := c.Param("id")

	for index, entry := range entries {
		if entry.ID == id {
			entries = removeElementByIndex(entries, index)
			c.Status(http.StatusNoContent)
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "entry not found"})
}

func updateEntry(c *gin.Context) {
	id := c.Param("id")

	for index, entry := range entries {
		if entry.ID == id {
			var newEntry Entry

			if err := c.BindJSON(&newEntry); err != nil {
				serviceLogger.Error("Failed to parse request", err)
				c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
				return
			}
			entries[index].Type = newEntry.Type
			entries[index].Message = newEntry.Message
			c.Status(http.StatusNoContent)
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "entry not found"})
}
