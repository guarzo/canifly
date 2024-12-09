package server

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"path/filepath"
	"runtime"
)

func SetupLogger() *logrus.Logger {
	logger := logrus.New()

	logger.SetReportCaller(true) // Enable caller reporting

	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			// Extract only the file name and line number
			filename := filepath.Base(frame.File)
			return frame.Function, fmt.Sprintf("%s:%d", filename, frame.Line)
		},
	})

	logger.SetLevel(logrus.InfoLevel)
	return logger
}
