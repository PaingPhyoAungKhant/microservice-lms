package logger

import "go.uber.org/zap"


func NewNop() *Logger {
	return &Logger{zap.NewNop()}
}
