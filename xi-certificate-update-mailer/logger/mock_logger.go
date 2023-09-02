package logger

import "go.uber.org/zap/zapcore"

type MockLogger struct{}

func (m *MockLogger) Info(msg string, fields ...zapcore.Field)  {}
func (m *MockLogger) Fatal(msg string, fields ...zapcore.Field) {}
