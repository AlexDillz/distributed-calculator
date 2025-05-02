package logging

import (
	"log"
	"os"
)

var (
	Logger *log.Logger
)

func InitLogger() {
	Logger = log.New(os.Stdout, "[CALC] ", log.Ldate|log.Ltime|log.Lshortfile)
}

func GetLogger() *log.Logger {
	return Logger
}
