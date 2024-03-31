package logger

import (
	"log"
	"sync"

	"go.uber.org/zap"
)

type Logger struct {
	ZapLogger *zap.Logger
}

var (
	Log  *Logger
	once sync.Once
)

func MustInitLogger() *Logger {
	once.Do(
		func() {
			zaplogger, err := zap.NewDevelopment()
			if err != nil {
				log.Fatal(err)
			}
			Log = &Logger{
				ZapLogger: zaplogger,
			}
		})
	return Log
}
