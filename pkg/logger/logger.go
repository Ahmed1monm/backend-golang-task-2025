package logger

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ctxKey string

const (
	traceIDKey     ctxKey = "traceID"
	defaultTraceID        = "unknown"
)

var globalLogger *zap.Logger

// Init initializes the global logger
func Init(env string) {
	config := zap.NewProductionConfig()
	if env == "development" {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
	
	logger, err := config.Build(
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
	if err != nil {
		panic("failed to initialize logger: " + err.Error())
	}

	globalLogger = logger
}

// WithTraceID adds traceID to the logger context
func WithTraceID(ctx context.Context) *zap.Logger {
	if traceID, ok := ctx.Value(traceIDKey).(string); ok {
		return globalLogger.With(zap.String("trace_id", traceID))
	}
	return globalLogger.With(zap.String("trace_id", defaultTraceID))
}

// Debug logs a debug message with trace ID
func Debug(ctx context.Context, msg string, fields ...zap.Field) {
	WithTraceID(ctx).Debug(msg, fields...)
}

// Info logs an info message with trace ID
func Info(ctx context.Context, msg string, fields ...zap.Field) {
	WithTraceID(ctx).Info(msg, fields...)
}

// Warn logs a warning message with trace ID
func Warn(ctx context.Context, msg string, fields ...zap.Field) {
	WithTraceID(ctx).Warn(msg, fields...)
}

// Error logs an error message with trace ID
func Error(ctx context.Context, msg string, fields ...zap.Field) {
	WithTraceID(ctx).Error(msg, fields...)
}

// Fatal logs a fatal message with trace ID and exits
func Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	WithTraceID(ctx).Fatal(msg, fields...)
}

// Sync flushes any buffered log entries
func Sync() error {
	return globalLogger.Sync()
}
