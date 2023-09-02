package logger

import "go.uber.org/zap"

var Log *zap.Logger

func Initialize() {
	var err error
	Log, err = zap.NewProduction()
	if err != nil {
		panic(err)
	}
}

func String(key string, value string) zap.Field {
	return zap.String(key, value)
}

func Int(key string, value int) zap.Field {
	return zap.Int(key, value)
}
