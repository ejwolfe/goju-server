package server

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
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

func loadEntries() {
	entries = readEntriesFile()

	if len(entries) == 0 {
		entries = append(entries,
			Entry{ID: uuid.NewString(), Type: Task, Message: "Simple task"},
			Entry{ID: uuid.NewString(), Type: Completed, Message: "Completed task"},
			Entry{ID: uuid.NewString(), Type: Migrated, Message: "Migrated task"},
			Entry{ID: uuid.NewString(), Type: Scheduled, Message: "Schedule task"},
			Entry{ID: uuid.NewString(), Type: Cancelled, Message: "Cancelled task"},
			Entry{ID: uuid.NewString(), Type: Note, Message: "Simple note"},
			Entry{ID: uuid.NewString(), Type: Event, Message: "Simple event"},
			Entry{ID: uuid.NewString(), Type: Priority, Message: "Priority task"},
		)
		writeEntriesFile(entries)
	}
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
