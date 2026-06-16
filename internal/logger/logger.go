package logger

import (
	"sync"

	"go.uber.org/zap"
)

var (
	log  *zap.Logger
	once sync.Once
)

// Init initializes the global Zap logger (production config).
// Safe to call multiple times; only the first call takes effect.
func Init() {
	once.Do(func() {
		var err error
		log, err = zap.NewProduction()
		if err != nil {
			panic("failed to initialize zap logger: " + err.Error())
		}
	})
}

// Get returns the global Zap logger instance.
// Panics if Init() has not been called.
func Get() *zap.Logger {
	if log == nil {
		panic("logger not initialized — call logger.Init() first")
	}
	return log
}

// Sync flushes any buffered log entries. Call before application exits.
func Sync() {
	if log != nil {
		_ = log.Sync()
	}
}
