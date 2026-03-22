package logger

import (
	"testing"

	"github.com/sirupsen/logrus"
)

func TestNewLogger_ReturnsStandardLogger(t *testing.T) {
	log := NewLogger()
	if log == nil {
		t.Fatal("NewLogger() returned nil")
	}
	if log != logrus.StandardLogger() {
		t.Error("NewLogger() should return logrus.StandardLogger()")
	}
}

func TestNewLogger_ReturnsSameInstance(t *testing.T) {
	a := NewLogger()
	b := NewLogger()
	if a != b {
		t.Error("NewLogger() should return the same logger instance each time")
	}
}
