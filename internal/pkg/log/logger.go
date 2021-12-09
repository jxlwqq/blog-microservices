package log

import (
	"go.uber.org/zap"
)

type Logger struct {
	*zap.SugaredLogger
}

func New() *Logger {
	l, _ := zap.NewProduction()
	return &Logger{l.Sugar()}
}
