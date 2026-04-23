package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

func Init(env string) {
	var config zap.Config

	if env == "production" {
		config = zap.NewProductionConfig()
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	var err error
	log, err = config.Build(zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}
}

func Info(msg string, fields ...zap.Field)  { log.Info(msg, fields...) }
func Error(msg string, fields ...zap.Field) { log.Error(msg, fields...) }
func Debug(msg string, fields ...zap.Field) { log.Debug(msg, fields...) }
func Warn(msg string, fields ...zap.Field)  { log.Warn(msg, fields...) }
func Fatal(msg string, fields ...zap.Field) { log.Fatal(msg, fields...) }

func Sync() error {
	return log.Sync()
}
