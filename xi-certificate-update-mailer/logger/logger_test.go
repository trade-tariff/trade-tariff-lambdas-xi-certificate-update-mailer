package logger

import (
	"testing"
)

func TestLoggerInitialization(t *testing.T) {
	if Log == nil {
		t.Error("Logger was not initialized")
	}

	if _, ok := Log.(*zapLogger); !ok {
		t.Error("Logger is not of type *logger.zapLogger")
	}
}
