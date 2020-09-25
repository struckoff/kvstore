package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger = defaultLog()

func InitLog(cfg zap.Config) error {
	var err error
	log, err = cfg.Build()
	if err != nil {
		return err
	}
	return nil
}

func Logger() *zap.Logger {
	if log == nil {
		cfg := DefaultConfig()
		InitLog(cfg)
	}
	return log
}

func defaultLog() *zap.Logger {
	cfg := DefaultConfig()
	l, _ := cfg.Build()
	return l
}

func DefaultConfig() zap.Config {
	return zap.Config{
		Encoding:         "console",
		Level:            zap.NewAtomicLevelAt(zapcore.DebugLevel),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message",

			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,

			TimeKey:    "time",
			EncodeTime: zapcore.ISO8601TimeEncoder,

			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,

			StacktraceKey: "stacktrace",
		},
	}
}
