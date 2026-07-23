package logger

import (
	"log"

	"go.uber.org/zap"
)

func Error(msg string, fields ...zap.Field) {

	if Log == nil {
		log.Println("[ERROR]", msg)
		return
	}

	Log.Error(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {

	if Log == nil {
		log.Println("[INFO]", msg)
		return
	}

	Log.Info(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {

	if Log == nil {
		log.Fatal(msg)
	}

	Log.Fatal(msg, fields...)
}
