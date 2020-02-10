package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Logger provides definition of the logger
type Logger struct {
	log *logrus.Logger
}

// New provides initialization of the logger
func New() *Logger {
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
	log.SetOutput(os.Stdout)
	return &Logger{
		log: log,
	}
}

// Fatalf for fatal errors
func (l *Logger) Fatalf(format string, data ...interface{}) {
	l.log.Fatalf(format, data...)
}

// Infof for info errors
func (l *Logger) Infof(format string, data ...interface{}) {
	l.log.Infof(format, data...)
}

// Errorf for errors with "Error" level
func (l *Logger) Errorf(format string, data ...interface{}) {
	l.log.Errorf(format, data...)
}
