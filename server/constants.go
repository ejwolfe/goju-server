package server

// File names & directories

const forwardSlash string = "/"
const serverLogFile string = "server.log"

// Endpoints

const colon string = ":"
const v1Path = forwardSlash + "v1"
const entriesPath = forwardSlash + "entries"
const idParam string = "id"
const entriesIdPath = entriesPath + forwardSlash + colon + idParam

// Headers

const locationHeader string = "Location"

// Server

const serverAddress string = "localhost:8080"

// Error messages

const error001 string = "Entry not found"
const error002 string = "Invalid request"
const error003 string = "Failed to parse request"
