package logger

import (
	"log"

	"go.uber.org/zap"
)

var logger *zap.SugaredLogger

func init() {
	baseLogger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("cam't start logger: %v", err)
	}

	logger = baseLogger.Sugar()
}

// Info - инфо лог
func Info(args ...interface{}) {
	logger.Info(args...)
}

// Infof - инфо лог с текстом
func Infof(template string, args ...interface{}) {
	logger.Infof(template, args...)
}

// Error - эррор лог
func Error(args ...interface{}) {
	logger.Error(args...)
}

// Errorf - эррор лог с текстом
func Errorf(template string, args ...interface{}) {
	logger.Errorf(template, args...)
}

// Fatalf - фатал
func Fatalf(template string, args ...interface{}) {
	logger.Fatalf(template, args...)
}
