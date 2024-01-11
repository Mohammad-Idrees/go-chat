package logger

import (
	"context"
	"fmt"
	"os"
	"project/constants"
	"reflect"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger *zap.Logger
)

func Info(ctx context.Context, msg string, fields ...zapcore.Field) {
	logger.Info(msg, addContext(ctx, fields)...)

}

func Error(ctx context.Context, msg string, fields ...zapcore.Field) {
	logger.Error(msg, addContext(ctx, fields)...)
}

func addContext(ctx context.Context, fields []zap.Field) []zap.Field {
	// TODO: bind request headers to context.
	XRequestID := ""
	return append(
		fields,
		zap.String(constants.XRequestID, XRequestID),
	)
}

func init() {
	// The bundled Config struct only supports the most common configuration
	// options. More complex needs, like splitting logs between multiple files
	// or writing to non-file outputs, require use of the zapcore package.
	//
	// In this example, imagine we're both sending our logs to Kafka and writing
	// them to the console. We'd like to encode the console output and the Kafka
	// topics differently, and we'd also like special treatment for
	// high-priority logs.

	// First, define our level-handling logic.
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel
	})

	// High-priority output should also go to standard error, and low-priority
	// output should also go to standard out.
	consoleDebugging := zapcore.Lock(os.Stdout)
	consoleErrors := zapcore.Lock(os.Stderr)

	// Optimize the Kafka output for machine consumption and the console output
	// for human operators.
	ecfg := zap.NewProductionEncoderConfig()
	ecfg.EncodeTime = zapcore.ISO8601TimeEncoder
	ecfg.EncodeLevel = zapcore.CapitalLevelEncoder

	consoleEncoder := zapcore.NewJSONEncoder(ecfg)

	// Join the outputs, encoders, and level-handling functions into
	// zapcore.Cores, then tee the four cores together.
	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, consoleErrors, highPriority),
		zapcore.NewCore(consoleEncoder, consoleDebugging, lowPriority),
	)

	// From a zapcore.Core, it's easy to construct a Logger.
	logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	defer logger.Sync()

	zap.RedirectStdLog(logger)

	logger.Info("constructed a logger")
}

func Field(k string, v interface{}) zap.Field {
	if v != nil && reflect.ValueOf(v).Kind() == reflect.Ptr {
		return zap.Any(k, fmt.Sprintf("%+v", v))
	}
	return zap.Any(k, v)
}
