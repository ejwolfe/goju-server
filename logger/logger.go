package logger

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type Log = func(v ...any)

type Logger struct {
	Debug Log
	Info  Log
	Warn  Log
	Error Log
}

func createLogFile() *os.File {
	file, err := os.Create("application.log")
	if err != nil {
		fmt.Println("Error creating log file: ", err)
	}

	// Close file when application terminates
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		file.Sync()
		file.Close()
		os.Exit(0)
	}()

	return file
}

func CreateLogger() Logger {
	file := createLogFile()
	flags := log.LstdFlags | log.Lmicroseconds
	debugLogger := log.New(file, "DEBUG ", flags)
	infoLogger := log.New(file, "INFO ", flags)
	warnLogger := log.New(file, "WARN ", flags)
	errorLogger := log.New(file, "ERROR ", flags)
	return Logger{
		Debug: debugLogger.Println,
		Info:  infoLogger.Println,
		Warn:  warnLogger.Println,
		Error: errorLogger.Println,
	}
}
