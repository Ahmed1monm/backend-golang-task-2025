package logger

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type GormLogger struct {
	SlowThreshold time.Duration
	LogLevel      gormlogger.LogLevel
}

// NewGormLogger creates a new GORM logger instance
func NewGormLogger() *GormLogger {
	return &GormLogger{
		SlowThreshold: 200 * time.Millisecond,
		LogLevel:      gormlogger.Info,
	}
}

func (l *GormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Info {
		Info(ctx, fmt.Sprintf(msg, data...), zap.String("source", "gorm"))
	}
}

func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Warn {
		Warn(ctx, fmt.Sprintf(msg, data...), zap.String("source", "gorm"))
	}
}

func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Error {
		Error(ctx, fmt.Sprintf(msg, data...), zap.String("source", "gorm"))
	}
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= gormlogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()
	fields := []zap.Field{
		zap.String("source", "gorm"),
		zap.Duration("elapsed", elapsed),
		zap.String("sql", sql),
	}

	if rows != -1 {
		fields = append(fields, zap.Int64("rows", rows))
	}

	switch {
	case err != nil && l.LogLevel >= gormlogger.Error && !errors.Is(err, gorm.ErrRecordNotFound):
		fields = append(fields, zap.Error(err))
		Error(ctx, "Database error", fields...)
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= gormlogger.Warn:
		fields = append(fields, zap.Duration("slow_threshold", l.SlowThreshold))
		Warn(ctx, "Slow SQL query", fields...)
	case l.LogLevel >= gormlogger.Info:
		Info(ctx, "SQL query executed", fields...)
	}
}
