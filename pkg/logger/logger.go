package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Logger provides definition of the logger
type Logger struct {
	log *logrus.Logger
}

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

func (l *Logger) Fatalf(format string, data ...interface{}) {
	l.log.Fatalf(format, data...)
}

func (l *Logger) Infof(format string, data ...interface{}) {
	l.log.Infof(format, data...)
}
