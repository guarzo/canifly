// services/logger/logrus_adapter.go
package logger

import (
	"github.com/guarzo/canifly/internal/services/interfaces"
	"github.com/sirupsen/logrus"
)

type LogrusAdapter struct {
	*logrus.Logger
}

func NewLogrusAdapter(l *logrus.Logger) interfaces.Logger {
	return &LogrusAdapter{l}
}
