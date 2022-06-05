package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"runtime"
)

var logger *zap.Logger

func init() {
	ExitsFile("log")

	cfg := zap.Config{
		Level:            zap.NewAtomicLevelAt(zapcore.Level(zapcore.DebugLevel)),
		Development:      true,
		Encoding:         "json",
		EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
		OutputPaths:      []string{"log/debug.log"},
		ErrorOutputPaths: []string{"stderr", "log/error.log"},
	}

	if runtime.GOOS == "windows" {
		cfg.Encoding = "console"
		cfg.OutputPaths = append(cfg.OutputPaths, "stdout")
	}

	var err error
	logger, err = cfg.Build(zap.AddCallerSkip(1))

	if err != nil {
		panic(err)
	}
}

func ExitsFile(path string) {
	_, err := os.Stat(path)

	if os.IsExist(err) {
		return
	}

	err = os.MkdirAll(path, os.ModeDir)

	if err != nil {
		panic(err)
	}
}

// Debug logs a message at DebugLevel. The message includes any fields passed
func Debug(msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)
}

// Info logs a message at InfoLevel. The message includes any fields passed
func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

// Wran logs a message at WarnLevel. The message includes any fields passed
func Warn(msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}

// Error logs a message at ErrorLevel. The message includes any fields passed
func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}

// Dpanic logs a message at PanicLevel. The message includes any fields passed
func DPanic(msg string, fields ...zap.Field) {
	logger.DPanic(msg, fields...)
}

// Panic logs a message at PanicLevel. The message includes any fields passed
func Panic(msg string, fields ...zap.Field) {
	logger.Panic(msg, fields...)
}

// Fatal logs a message at FatalLevel. The message includes any fields passed
func Fatal(msg string, fields ...zap.Field) {
	logger.Fatal(msg, fields...)
}
