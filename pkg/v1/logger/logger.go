package logger

import "go.uber.org/zap"

var Log *zap.Logger

func Init(cfg Config) error {

	l, err := New(cfg)
	if err != nil {
		return err
	}

	Log = l

	return nil
}

func Sync() {
	if Log != nil {
		_ = Log.Sync()
	}
}
