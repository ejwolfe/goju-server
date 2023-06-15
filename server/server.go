package server

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

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
	ID      int       `json:"id"`
	Type    EntryType `json:"type"`
	Message string    `json:"message"`
}

var entries []Entry

var serviceLogger = logger.CreateLogger()

func CreateServer() {
	router := gin.Default()
	router.Use(Logger())

	entries = readEntriesFile()

	router.GET(CONTEXT, getEntries)
	router.GET(CONTEXT+"/:id", getEntryByID)
	router.POST(CONTEXT, addEntry)

	router.Run("localhost:8080")
}

func getEntries(c *gin.Context) {
	serviceLogger.Info("Request Headers", c.Request.Header)
	c.IndentedJSON(http.StatusOK, entries)
	serviceLogger.Info("Response Body", http.StatusOK, entries)
}

func getEntryByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		serviceLogger.Error("Failed to parse id", err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}

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

	newEntry.ID = len(entries)

	entries = append(entries, newEntry)
	c.IndentedJSON(http.StatusCreated, newEntry)
	writeEntriesFile(entries)
}
