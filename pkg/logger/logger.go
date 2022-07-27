package logger

import (
	"go.uber.org/zap"
)

var singleton zap.Logger

type Config struct {
	Encoding string
	Level    string
}

func NewDefaultConfig() *Config {
	return &Config{
		Encoding: "console",
		Level:    "info",
	}
}

var Log *zap.Logger

// Init initializes singleton logger
func Init(c *Config) error {
	if Log == nil {
		zapConfig := zap.NewProductionConfig()
		zapConfig.Encoding = c.Encoding
		err := zapConfig.Level.UnmarshalText([]byte(c.Level))
		if err != nil {
			return err
		}

		log, err := zapConfig.Build()
		if err != nil {
			return err
		}
		Log = log
	} else {
		_ = Log.Sync()
	}

	defer func() {
		_ = Log.Sync()
	}()
	return nil
}

func Debug(message string, fields ...zap.Field) {
	singleton.Debug(message, fields...)
}

func Info(message string, fields ...zap.Field) {
	singleton.Info(message, fields...)
}

func Warn(message string, fields ...zap.Field) {
	singleton.Warn(message, fields...)
}

func Error(message string, fields ...zap.Field) {
	singleton.Error(message, fields...)
}

func Fatal(message string, fields ...zap.Field) {
	singleton.Fatal(message, fields...)
}
