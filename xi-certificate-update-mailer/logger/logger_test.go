package logger

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestInitialize(t *testing.T) {
	Initialize()
	assert.NotNil(t, Log, "Logger should be initialized")
}
