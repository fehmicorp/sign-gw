package logger

import (
	"gopkg.in/natefinch/lumberjack.v2"
)

func NewRotateWriter(cfg Config) *lumberjack.Logger {

	return &lumberjack.Logger{

		Filename: cfg.File,

		MaxSize: cfg.MaxSize,

		MaxBackups: cfg.MaxBackups,

		MaxAge: cfg.MaxAge,

		Compress: cfg.Compress,
	}
}
