package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

type Logger struct {
	warningLogger *log.Logger
	infoLogger    *log.Logger
	errorLogger   *log.Logger

	file os.File
}

func NewLogger(name string) (Logger, error) {
	var l Logger
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	basepath += "/.."

	path := filepath.Join(basepath, "log")
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return l, err
	}

	path = filepath.Join(basepath, "log", "log.txt")
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return l, err
	}

	l.file = *file

	l.infoLogger = log.New(&l.file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	l.warningLogger = log.New(&l.file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	l.errorLogger = log.New(&l.file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	return l, nil
}

func (l Logger) LogInfo(trasactionId int, msg string) {
	l.infoLogger.Println(fmt.Sprintf("T: %d, Msg : %s", trasactionId, msg))
}

func (l Logger) LogWarning(trasactionId int, msg string) {
	l.warningLogger.Println(fmt.Sprintf("T: %d, Msg : %s", trasactionId, msg))
}

func (l Logger) LogError(trasactionId int, msg string) {
	l.errorLogger.Println(fmt.Sprintf("T: %d, Msg : %s", trasactionId, msg))
}
