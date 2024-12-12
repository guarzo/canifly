package server

import (
	"github.com/sirupsen/logrus"

	flyLogger "github.com/guarzo/canifly/internal/services/config"
	"github.com/guarzo/canifly/internal/services/interfaces"
)

func SetupLogger() interfaces.Logger {

	logrusLogger := logrus.New()

	// logrusLogger.SetReportCaller(true) // Enable caller reporting

	logrusLogger.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: false,
		//CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
		//	// Extract only the file name and line number
		//	filename := filepath.Base(frame.File)
		//	return frame.Function, fmt.Sprintf("%s:%d", filename, frame.Line)
		//},
	})

	logrusLogger.SetLevel(logrus.InfoLevel)
	return flyLogger.NewLogrusAdapter(logrusLogger)
}
