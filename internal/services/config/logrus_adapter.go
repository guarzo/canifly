// services/config/logrus_adapter.go
package config

import (
	"github.com/sirupsen/logrus"

	"github.com/guarzo/canifly/internal/services/interfaces"
)

type LogrusAdapter struct {
	entry *logrus.Entry
}

func NewLogrusAdapter(l *logrus.Logger) interfaces.Logger {
	return &LogrusAdapter{entry: logrus.NewEntry(l)}
}

func (l *LogrusAdapter) Debug(args ...interface{}) {
	l.entry.Debug(args...)
}

func (l *LogrusAdapter) Debugf(format string, args ...interface{}) {
	l.entry.Debugf(format, args...)
}

func (l *LogrusAdapter) Info(args ...interface{}) {
	l.entry.Info(args...)
}

func (l *LogrusAdapter) Infof(format string, args ...interface{}) {
	l.entry.Infof(format, args...)
}

func (l *LogrusAdapter) Warn(args ...interface{}) {
	l.entry.Warn(args...)
}

func (l *LogrusAdapter) Warnf(format string, args ...interface{}) {
	l.entry.Warnf(format, args...)
}

func (l *LogrusAdapter) Error(args ...interface{}) {
	l.entry.Error(args...)
}

func (l *LogrusAdapter) Errorf(format string, args ...interface{}) {
	l.entry.Errorf(format, args...)
}

func (l *LogrusAdapter) Fatal(args ...interface{}) {
	l.entry.Fatal(args...)
}

func (l *LogrusAdapter) Fatalf(format string, args ...interface{}) {
	l.entry.Fatalf(format, args...)
}

// WithError creates a new LogrusAdapter with the error field set, returning Logger
func (l *LogrusAdapter) WithError(err error) interfaces.Logger {
	return &LogrusAdapter{entry: l.entry.WithError(err)}
}

func (l *LogrusAdapter) WithField(key string, value interface{}) interfaces.Logger {
	return &LogrusAdapter{entry: l.entry.WithField(key, value)}
}

func (l *LogrusAdapter) WithFields(fields map[string]interface{}) interfaces.Logger {
	return &LogrusAdapter{
		entry: l.entry.WithFields(logrus.Fields(fields)),
	}
}
