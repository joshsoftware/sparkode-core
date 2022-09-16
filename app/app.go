package app

import (
	"go.uber.org/zap"
)

var (
	logger *zap.SugaredLogger
)

func Init() {
	InitLogger()
}

func InitLogger() {
	zapLogger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	logger = zapLogger.Sugar()
}

func GetLogger() *zap.SugaredLogger {
	Init()
	return logger
}

func Close() {
	logger.Sync()
}
