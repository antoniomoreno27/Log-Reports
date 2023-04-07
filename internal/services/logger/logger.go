package logger

import (
	"log"
	"os"
	"sync"
)

var (
	once sync.Once

	// singleton instance
	logger Logger
)

type Logger struct {
	general *log.Logger
	err     *log.Logger
}

func init() {
	once.Do(func() {
		logger = Logger{
			general: log.New(os.Stdout, "Warning: ", log.Ldate|log.Ltime|log.Lmicroseconds),
			err:     log.New(os.Stdout, "Error: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile),
		}
	})
}

func Warnf(message string, args ...interface{}) {
	logger.general.Printf(message+"\n", args...)
}

func Errorf(message string, args ...interface{}) {
	logger.err.Printf(message+"\n", args...)
}

func Panic(message string) {
	logger.err.Panic(message)
}
