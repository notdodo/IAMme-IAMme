package logging

import (
	"os"
	"sync"
	"time"

	"github.com/charmbracelet/log"
)

type LogManager interface {
	SetVerboseLevel()
	SetDebugLevel()
	Debug(message interface{}, keyvals ...interface{})
	Info(message interface{}, keyvals ...interface{})
	Warn(message interface{}, keyvals ...interface{})
	Error(message interface{}, keyvals ...interface{})
}

type logManager struct {
	logger *log.Logger
}

var logger *logManager
var once sync.Once

func GetLogManager() LogManager {
	once.Do(func() {
		logger = &logManager{
			logger: log.NewWithOptions(os.Stdout, log.Options{
				CallerOffset:    1,
				Level:           log.WarnLevel,
				ReportCaller:    true,
				ReportTimestamp: true,
				TimeFormat:      time.RFC1123,
			}),
		}
	})

	return *logger
}

func (lm logManager) SetVerboseLevel() {
	lm.logger.SetLevel(log.InfoLevel)
}

func (lm logManager) SetDebugLevel() {
	lm.logger.SetLevel(log.DebugLevel)
}

func (lm logManager) Debug(message interface{}, keyvals ...interface{}) {
	lm.logger.Debug(message, keyvals...)
}

func (lm logManager) Info(message interface{}, keyvals ...interface{}) {
	lm.logger.Info(message, keyvals...)
}

func (lm logManager) Warn(message interface{}, keyvals ...interface{}) {
	lm.logger.Warn(message, keyvals...)
}

func (lm logManager) Error(message interface{}, keyvals ...interface{}) {
	lm.logger.Error(message, keyvals...)
	os.Exit(1)
}
