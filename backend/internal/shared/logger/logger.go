// internal/shared/logger/logger.go
package logger

import (
	"log"
	"os"
)

var (
	infoLogger  = log.New(os.Stdout, "INFO: ", log.LstdFlags|log.Lshortfile)
	errorLogger = log.New(os.Stderr, "ОШИБКА: ", log.LstdFlags|log.Lshortfile)
)

func Info(msg string) {
	infoLogger.Println(msg)
}

func Error(msg string, err error) {
	errorLogger.Printf("%s: %v\n", msg, err)
}

func Fatal(msg string, err error) {
	errorLogger.Fatalf("%s: %v\n", msg, err)
}
