package logger

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.Logger
}

type contextKey struct{}

var (
	globalLogger *Logger
	loggerKey    = &contextKey{}
)

func Init(env string) error {
	l, err := New(env)
	if err != nil {
		return err
	}
	globalLogger = l
	return nil
}

func Std() *Logger {
	return globalLogger
}

func Sync() error {
	if globalLogger == nil {
		return nil
	}
	return globalLogger.Sync()
}

func NewContext(ctx context.Context, l *Logger) context.Context {
	return context.WithValue(ctx, loggerKey, l)
}

func FromContext(ctx context.Context) *Logger {
	if ctx == nil {
		return globalLogger
	}
	if l, ok := ctx.Value(loggerKey).(*Logger); ok && l != nil {
		return l
	}
	return globalLogger
}

func New(env string) (*Logger, error) {
	var cfg zap.Config

	if env == "production" {
		cfg = zap.NewProductionConfig()
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	l, err := cfg.Build(zap.AddCallerSkip(1))
	if err != nil {
		return nil, err
	}

	return &Logger{Logger: l}, nil
}

func (l *Logger) With(fields ...zap.Field) *Logger {
	return &Logger{Logger: l.Logger.With(fields...)}
}

func Debug(msg string, fields ...zap.Field) {
	if globalLogger != nil {
		globalLogger.Debug(msg, fields...)
	}
}

func Info(msg string, fields ...zap.Field) {
	if globalLogger != nil {
		globalLogger.Info(msg, fields...)
	}
}

func Warn(msg string, fields ...zap.Field) {
	if globalLogger != nil {
		globalLogger.Warn(msg, fields...)
	}
}

func Error(msg string, fields ...zap.Field) {
	if globalLogger != nil {
		globalLogger.Error(msg, fields...)
	}
}

func Fatal(msg string, fields ...zap.Field) {
	if globalLogger != nil {
		globalLogger.Fatal(msg, fields...)
	}
}
