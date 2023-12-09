package logging

import (
	"testing"

	"github.com/charmbracelet/log"
)

func TestSetVerbose(t *testing.T) {
	myLog := GetLogManager()
	myLog.SetVerboseLevel()

	if myLog.(logManager).logger.GetLevel().String() != log.InfoLevel.String() {
		t.Errorf("expected: %s\ngot: %s", log.InfoLevel, myLog.(logManager).logger.GetLevel().String())
	}
}

func TestSetDebug(t *testing.T) {
	myLog := GetLogManager()
	myLog.SetDebugLevel()

	if myLog.(logManager).logger.GetLevel().String() != log.DebugLevel.String() {
		t.Errorf("expected: %s\ngot: %s", log.InfoLevel, myLog.(logManager).logger.GetLevel().String())
	}
}
