package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

func setupLogging() {
	if config.LogFile != "" {
		f, err := os.OpenFile(config.LogFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		errFatal(err, "Failed to open "+config.LogFile)
		log.SetOutput(f)
	}
	if config.LogFormat == "json" {
		log.SetFormatter(&log.JSONFormatter{})
	}
}

// errPanic something went wrong, panic.
func errPanic(err error) {
	if err != nil {
		panic(err)
	}
}

// nilPanic does so when it's parameter is such.
func nilPanic(val interface{}) {
	if val == nil {
		fmt.Printf("val is nil")
	}
}

// errPanic something went wrong, panic.
func errFatal(err error, str string) {
	if err != nil {
		log.Fatalln(err, str)
	}
}

// nilPanic does so when it's parameter is such.
func nilFatal(val interface{}, str string) {
	if val == nil {
		log.Fatalln(str)
	}
}
