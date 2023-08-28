package logging

import (
	"fmt"
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
	message = fmt.Sprintf("%s %s", "INFO", message)
	fmt.Println(message)
	log.Println(message)
}

func Error(message string) {
	message = fmt.Sprintf("%s %s", "ERROR", message)
	fmt.Println(message)
	log.Println(message)
}

func Panic(message string) {
	message = fmt.Sprintf("%s %s", "PANIC", message)
	fmt.Println(message)
	log.Panicf(message)
}

func Close() {
	logFile.Close()
}
