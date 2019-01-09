package main

import (
	"log"
	"os"
)

// nolint: gochecknoglobals
var (
	Logger *log.Logger

	// Fd is a logfile declared global to be closed in main()
	Fd *os.File
)

// nolint: gochecknoinits
func init() {

	// flag for opening LogFile
	var flag int

	_, err := os.Stat(LogFilePath)

	// if logfile does not exist - will also create the file
	switch os.IsNotExist(err) {
	case true:
		flag = os.O_WRONLY | os.O_CREATE | os.O_APPEND
	case false:
		flag = os.O_WRONLY | os.O_APPEND
	}

	// create/open log file
	Fd, err := os.OpenFile(LogFilePath, flag, 0666)
	if err != nil {
		log.Fatalf("error opening log file: %v", err)
	}

	// configure logger
	Logger = log.New(Fd, "", log.Ldate|log.Ltime|log.Lshortfile)

}
