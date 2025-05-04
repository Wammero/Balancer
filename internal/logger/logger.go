package logger

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

func New() *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logger.SetOutput(os.Stdout)
	return logger
}

func LogRequest(logger *logrus.Logger, level logrus.Level, reqID, method, path, backend string, statusCode int, duration time.Duration, err error) {
	fields := logrus.Fields{
		"request_id": reqID,
		"method":     method,
		"path":       path,
		"backend":    backend,
		"duration":   duration,
	}

	if err != nil {
		fields["error"] = err.Error()
	}

	if statusCode != 0 {
		fields["status"] = statusCode
	}

	logger.WithFields(fields).Log(level, "Request processed")
}
