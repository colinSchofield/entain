package logging

import (
	"os"
	"sync"

	log "github.com/sirupsen/logrus"
)

var logger *log.Logger
var once sync.Once

// TODO configure for CloudWatch ingestion..
func initLogger() {
	logger = &log.Logger{
		Out:       os.Stdout,
		Level:     log.InfoLevel,
		Formatter: &log.JSONFormatter{},
	}
}

// Logger returns a lazily initialised singleton logger
func Logger() *log.Logger {
	once.Do(initLogger)
	return logger
}
