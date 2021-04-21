package logger

import (
	"testing"
)

func TestLogger(t *testing.T) {
	log1 := GetLoggerModule("help")
	log2 := GetLoggerModule("log2")
	log1.Info("test")
	log2.Error("hello world")
	log2.Warning("world")

	// logger.Debug("bug")
}
