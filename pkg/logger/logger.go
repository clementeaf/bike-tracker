package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func InitLogger() {
	log.Out = os.Stdout
	log.SetFormatter(&logrus.JSONFormatter{})
}

func Info(msg string, fields logrus.Fields) {
	log.WithFields(fields).Info(msg)
}

func Error(msg string, fields logrus.Fields) {
	log.WithFields(fields).Error(msg)
}
