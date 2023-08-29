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
	defer baseLogger.Sync()

	logger = baseLogger.Sugar()
}

func Info(args ...interface{}) {
	logger.Info(args...)
}

func Infof(template string, args ...interface{}) {
	logger.Infof(template, args...)
}

func Error(args ...interface{}) {
	logger.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	logger.Errorf(template, args...)
}

func Fatalf(template string, args ...interface{}) {
	logger.Fatalf(template, args...)
}
