// Package logger provides a pre-configured application logger.
package logger

import (
	"github.com/sirupsen/logrus"
)

// NewLogger returns the standard logrus logger instance.
func NewLogger() *logrus.Logger {
	return logrus.StandardLogger()
}
