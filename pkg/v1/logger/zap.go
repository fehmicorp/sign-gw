package logger

import (
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(cfg Config) (*zap.Logger, error) {

	level := zapcore.InfoLevel

	switch strings.ToLower(cfg.Level) {

	case "debug":
		level = zapcore.DebugLevel

	case "warn":
		level = zapcore.WarnLevel

	case "error":
		level = zapcore.ErrorLevel
	}

	encoderCfg := zap.NewProductionEncoderConfig()

	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	encoder := zapcore.NewConsoleEncoder(
		encoderCfg,
	)

	fileWriter := zapcore.AddSync(
		NewRotateWriter(cfg),
	)

	var core zapcore.Core

	if cfg.Console {

		consoleWriter := zapcore.AddSync(os.Stdout)

		core = zapcore.NewTee(

			zapcore.NewCore(
				encoder,
				consoleWriter,
				level,
			),

			zapcore.NewCore(
				encoder,
				fileWriter,
				level,
			),
		)

	} else {

		core = zapcore.NewCore(
			encoder,
			fileWriter,
			level,
		)

	}

	return zap.New(
		core,
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	), nil
}
