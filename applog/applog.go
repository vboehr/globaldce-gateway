package applog

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
)

var consoleLogger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

////////////////////////

var fileLogger *log.Logger
var once sync.Once
var fileLoggerFile *os.File
var displayunlocked bool
var displaytrace = true

////////////////////////

func LockDisplay() {
	displayunlocked = false
}
func UnlockDisplay() {
	displayunlocked = true
}
func EnableDisplayTrace() {
	displaytrace = true
}

func Init(tmpapppath string) {
	displayunlocked = true

	once.Do(func() {
		infologpath := filepath.Join(tmpapppath, "info.log")
		file, err := os.OpenFile(infologpath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			log.Println(err)
		}
		fileLogger = log.New(file, "", log.Ldate|log.Ltime|log.Lshortfile)
		fileLoggerFile = file

	})

}

func Close() {

	fileLoggerFile.Sync()
	fileLoggerFile.Close()

}

func Trace(pattern string, args ...interface{}) {
	if displayunlocked && displaytrace {
		consoleLogger.Output(2, fmt.Sprintf(pattern+"\n", args...))
	}

	fileLogger.Output(2, fmt.Sprintf(pattern+"\n", args...))
}

func Notice(pattern string, args ...interface{}) {
	if displayunlocked {
		//consoleLogger.Output(2, fmt.Sprintf(pattern+"\n", args...))
		fmt.Printf(pattern+"\n", args...)
	}

	fileLogger.Output(2, fmt.Sprintf(pattern+"\n", args...))
}

func Warning(pattern string, args ...interface{}) {
	if displayunlocked {
		consoleLogger.Output(2, fmt.Sprintf(pattern+"\n", args...))
	}

	fileLogger.Output(2, fmt.Sprintf(pattern+"\n", args...))
}

func Fatal(pattern string, args ...interface{}) {
	if displayunlocked {
		consoleLogger.Output(2, fmt.Sprintf(pattern+"\n", args...))
	}

	fileLogger.Output(2, fmt.Sprintf(pattern+"\n", args...))
}
