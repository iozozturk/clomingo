package monitoring

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
)

func NewLogger() *zap.Logger {
	productionConfig := zap.NewProductionConfig()
	productionConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logger, err := productionConfig.Build()
	if err != nil {
		panic("cannot initiate logger")
	}
	return logger
}

func Sync(logger *zap.Logger) {
	if err := logger.Sync(); err != nil {
		log.Println("logger sync error: ", err.Error())
	}
}
