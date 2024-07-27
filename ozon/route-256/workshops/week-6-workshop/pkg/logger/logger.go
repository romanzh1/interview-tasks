package logger

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sync"
)

type loggerCtx string

const (
	loggerCtxValue loggerCtx = "logger"
)

type Logger struct {
	logger *zap.SugaredLogger
}

var globalLogger *Logger

func New() (*Logger, error) {
	loggerConfig := zap.NewProductionConfig()
	loggerConfig.Level.SetLevel(zapcore.DebugLevel)
	loggerConfig.ErrorOutputPaths = []string{"stdout"}

	logger, err := loggerConfig.Build()
	if err != nil {
		return nil, fmt.Errorf("loggerConfig.Build: %w", err)
	}

	once := sync.Once{}
	once.Do(func() {
		globalLogger = &Logger{logger: logger.Sugar()}
	})

	return &Logger{logger: logger.Sugar()}, nil
}

func ToContext(ctx context.Context, logger *Logger) context.Context {
	return context.WithValue(ctx, loggerCtxValue, logger)
}

func Infow(ctx context.Context, msg string, keysAndValues ...interface{}) {
	if loggerC, ok := ctx.Value(loggerCtxValue).(*Logger); ok {
		loggerC.logger.Infow(msg, keysAndValues...)

		return
	}

	globalLogger.logger.Infow(msg, keysAndValues...)
}

func Errorw(ctx context.Context, msg string, keysAndValues ...interface{}) {
	if loggerC, ok := ctx.Value(loggerCtxValue).(*Logger); ok {
		loggerC.logger.Errorw(msg, keysAndValues...)

		return
	}

	globalLogger.logger.Errorw(msg, keysAndValues...)
}

func Panicw(ctx context.Context, msg string, keysAndValues ...interface{}) {
	if loggerC, ok := ctx.Value(loggerCtxValue).(*Logger); ok {
		loggerC.logger.Panicw(msg, keysAndValues...)

		return
	}

	globalLogger.logger.Panicw(msg, keysAndValues...)
}

func With(args ...interface{}) *Logger {
	return &Logger{globalLogger.logger.With(args...)}

}
