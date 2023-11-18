package logging

import (
	"os"
	"time"

	"github.com/charmbracelet/log"
)

type LogManager interface {
	Debug(message interface{}, keyvals ...interface{})
	Info(message interface{}, keyvals ...interface{})
	Warn(message interface{}, keyvals ...interface{})
	Error(message interface{}, keyvals ...interface{})
}

type logManager struct {
	logger *log.Logger
}

func NewLogManager() LogManager {
	return &logManager{
		logger: log.NewWithOptions(os.Stdout, log.Options{
			CallerOffset:    1,
			Fields:          []interface{}{"err", "flag"},
			Level:           log.WarnLevel,
			ReportCaller:    true,
			ReportTimestamp: true,
			TimeFormat:      time.RFC1123,
		}),
	}
}

func (lm *logManager) Debug(message interface{}, keyvals ...interface{}) {
	lm.logger.Debug(message, keyvals)
}

func (lm *logManager) Info(message interface{}, keyvals ...interface{}) {
	lm.logger.Info(message, keyvals...)
}

func (lm *logManager) Warn(message interface{}, keyvals ...interface{}) {
	lm.logger.Info(message, keyvals...)
}

func (lm *logManager) Error(message interface{}, keyvals ...interface{}) {
	lm.logger.Error(message, keyvals...)
	os.Exit(1)
}
