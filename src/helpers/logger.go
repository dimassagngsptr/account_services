package helpers

import (
	"os"

	"github.com/sirupsen/logrus"
)

func InitLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)

	logger.SetLevel(logrus.DebugLevel)

	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})
	return logger
}

func LogWithFields(level logrus.Level, event string, method string, data interface{}, message string) {
	logger := InitLogger()
	logger.WithFields(logrus.Fields{
		"event":  event,
		"method": method,
		"data":   data,
	}).Log(level, message)
}
