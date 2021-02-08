package logger

import (
	"errors"
	"os"
	"testing"
)

func InitTestLog() *Logger {
	os.Setenv("TEST_LOG_SYSLOG", "true")
	os.Setenv("TEST_LOG_LEVEL", "trace")
	os.Setenv("TEST_LOG_FILE", "/dev/stdout")
	os.Setenv("TEST_LOG_PREFIX", "test -- ")
	l := GetLoggerFromEnv("TEST_", true)
	l.Logger.SetFlags(0)
	return l
}

func TestInitWithoutEnv(t *testing.T) {
	l := GetLoggerFromEnv("", true)
	l.Info("OK")
}

func TestSetPrefix(t *testing.T) {
	l := InitTestLog()
	l.SetPrefix("Something else :: ")
	l.Info("OK")
}

func TestLog(t *testing.T) {
	l := InitTestLog()
	l.Trace("test trace")
	l.Debug("test debug")
	l.Info("test info")
	l.Sys("test syslog")
	l.Warning("test warning")
	l.Error("test error")
}

func TestLogf(t *testing.T) {
	l := InitTestLog()
	l.Tracef("%s trace", "test")
	l.Debugf("%s debug", "test")
	l.Infof("%s info", "test")
	l.Sysf("%s syslog", "test")
	l.Warningf("%s warning", "test")
	l.Errorf("%s error", "test")
}

func TestErrore(t *testing.T) {
	l := InitTestLog()
	e := errors.New("Fake error 1")
	l.Warninge(e, "Test warning : %s", e)
	l.Errore(e, "Test error : %s", e)
}

func TestPanice(t *testing.T) {
	l := InitTestLog()
	defer func(l *Logger) {
		if r := recover(); r != nil {
			l.Recoverf("Recover from panic : \"%s\"", r)
		}
	}(l)
	e := errors.New("Fake error 2")
	l.Panice(e, "Test panic : %s", e)
}

func TestLevel(t *testing.T) {
	os.Setenv("LOG_LEVEL", "error")
	l := GetLoggerFromEnv("", true)
	l.Trace("test trace")
	l.Debug("test debug")
	l.Info("test info")
	l.Sys("test syslog")
	l.Warning("test warning")
	l.Error("test error")
}

func TestNoError(t *testing.T) {
	l := InitTestLog()
	var e error
	l.Panice(e, "Test no panic : %s", e)
	l.Errore(e, "Test no error : %s", e)
	l.Fatale(e, "Test no fatal : %s", e)
}

func TestLogFile(t *testing.T) {
	f := "test-tmp/test.log"
	os.Setenv("LOG_FILE", f)
	l := GetLoggerFromEnv("", true)
	l.Error("test error file")
	os.RemoveAll("test-tmp")
}

/*
func TestFatale(t *testing.T) {
	l := InitTestLog()
	e := errors.New("Fake error 3")
	l.Fatale(e, "Fatal error : \"%s\"", e)
}
*/
