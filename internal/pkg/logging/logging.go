package logging

import (
	"log"
	"os"
)

var logFile *os.File

func init() {
	var err error
	logFile, err = os.OpenFile("import.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(logFile)
}

func Info(message string) {
	log.Printf("%s %s", "INFO", message)
}

func Error(message string) {
	log.Printf("%s %s", "ERROR", message)
}

func Panic(message string) {
	log.Panicf("%s %s", "PANIC", message)
}

func Close() {
	logFile.Close()
}
