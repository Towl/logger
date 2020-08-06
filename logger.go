package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime/debug"
	"strconv"
)

// Level of the log for output
type Level int

// List of available levels
const (
	Error Level = iota + 1
	Warning
	Info
	Debug
	Trace
)

var logLevelLookup = map[string]Level{
	"trace":   Trace,
	"debug":   Debug,
	"info":    Info,
	"warning": Warning,
	"error":   Error,
}

// Logger object for logs
type Logger struct {
	EnvPrefix string
	Logger    *log.Logger
	Level     Level
	SysLog    bool
}

var logger *Logger

func init() {
	log.SetFlags(0)
	logger = GetLoggerFromEnv("", false)
}

// GetLoggerFromEnv use LOG_SYSLOG, LOG_FILE, LOG_LEVEL and LOG_PREFIX to
// instantiate the logger
// Create a new one if n = true
func GetLoggerFromEnv(w string, n bool) *Logger {
	if n == true || logger == nil || logger.EnvPrefix != w {
		s, e := strconv.ParseBool(os.Getenv(w + "LOG_SYSLOG"))
		if e != nil {
			log.Printf("[WARN] Unable to parse %sLOG_SYSLOG : %s\n", w, e)
			s = false
			log.Printf("[WARN] Using default value \"%v\"\n", s)
		}

		f := os.Getenv(w + "LOG_FILE")
		if len(f) < 3 {
			log.Printf("[WARN] Invalid value for %sLOG_FILE : \"%s\"\n", w, f)
			f = "/dev/stdout"
			log.Printf("[WARN] Using default value \"%s\"\n", f)
		}

		l := os.Getenv(w + "LOG_LEVEL")
		if _, ok := logLevelLookup[l]; !ok {
			log.Printf("[WARN] Invalid value for %sLOG_LEVEL \"%s\"\n", w, l)
			l = "info"
			log.Printf("[WARN] Using default value \"%s\"\n", l)
		}

		p := os.Getenv(w + "LOG_PREFIX")
		logger = NewLogger(f, p, l, s)
		logger.EnvPrefix = w
	}

	return logger
}

// NewLogger return a Logger instance from configs
func NewLogger(f string, p string, v string, b bool) *Logger {
	i := PrepareFile(f)
	l := &Logger{
		Logger: log.New(i, p, log.Ldate|log.Ltime|log.Lmicroseconds),
		Level:  logLevelLookup[v],
		SysLog: b,
	}
	l.Debug("New logger created")
	return l
}

// PrepareFile for writing logs
func PrepareFile(f string) *os.File {
	d := filepath.Dir(f)
	if _, e := os.Stat(d); os.IsNotExist(e) {
		if e := os.MkdirAll(d, 0775); e != nil {
			log.Panic(e)
		}
	}
	if _, e := os.Stat(f); os.IsNotExist(e) {
		if _, e := os.Create(f); e != nil {
			log.Panic(e)
		}
	}
	l, e := os.OpenFile(f, os.O_RDWR|os.O_APPEND, 0666)
	if e != nil {
		log.Panic(e)
	}
	return l
}

// SetPrefix message on logs
func (l *Logger) SetPrefix(p string) {
	l.Logger.SetPrefix(p)
}

// Trace logs Trace string
func (l *Logger) Trace(s string) {
	if l.Level >= Trace {
		l.log(s, "TRACE")
	}
}

// Tracef logs Trace string using Sprintf
func (l *Logger) Tracef(s string, a ...interface{}) {
	l.Trace(fmt.Sprintf(s, a...))
}

// Debug logs Debug string
func (l *Logger) Debug(s string) {
	if l.Level >= Debug {
		l.log(s, "DEBUG")
	}
}

// Debugf logs Debug string using Sprintf
func (l *Logger) Debugf(s string, a ...interface{}) {
	l.Debug(fmt.Sprintf(s, a...))
}

// Info logs Info string
func (l *Logger) Info(s string) {
	if l.Level >= Info {
		l.log(s, "INFO")
	}
}

// Infof logs string using Sprintf
func (l *Logger) Infof(s string, a ...interface{}) {
	l.Info(fmt.Sprintf(s, a...))
}

// Warning logs Warning string
func (l *Logger) Warning(s string) {
	l.Sys(s)
	if l.Level >= Warning {
		l.log(s, "WARN")
	}
}

// Warningf logs Warning Sprintf
func (l *Logger) Warningf(s string, a ...interface{}) {
	l.Warning(fmt.Sprintf(s, a...))
}

// Warninge logs string if Error not nil
func (l *Logger) Warninge(e error, s string, a ...interface{}) {
	if e != nil {
		l.Warningf(s, a...)
	}
}

// Error logs error on output and file
func (l *Logger) Error(s string) {
	l.Sys(s)
	l.log(s, "ERROR")
}

// Errorf logs Sprintf on output and file
func (l *Logger) Errorf(s string, a ...interface{}) {
	l.Error(fmt.Sprintf(s, a...))
}

// Errore logs Sprintf if Error not nil on output and file
func (l *Logger) Errore(e error, s string, a ...interface{}) {
	if e != nil {
		l.Errorf(s, a...)
	}
}

// Panic logs string then Panic
func (l *Logger) Panic(s string) {
	l.Sys(s)
	l.log(s, "PANIC")
	log.Panic(s)
}

// Panicf logs Sprintf then Panic
func (l *Logger) Panicf(s string, a ...interface{}) {
	l.Panic(fmt.Sprintf(s, a...))
}

// Recover logs string
func (l *Logger) Recover(s string) {
	l.Sys(s)
	l.log(s, "RECOVER")
	l.Debugf("Panic was : %s", debug.Stack())
}

// Recoverf logs Sprintf
func (l *Logger) Recoverf(s string, a ...interface{}) {
	l.Recover(fmt.Sprintf(s, a...))
}

// Panice logs Error if not nil then Panic
func (l *Logger) Panice(e error, s string, a ...interface{}) {
	if e != nil {
		l.Panicf(s, a...)
	}
}

// Fatal logs string then exit
func (l *Logger) Fatal(s string) {
	l.log(s, "FATAL")
	l.Sys(s)
	log.Fatal(s)
}

// Fatalf logs Sprintf then exit
func (l *Logger) Fatalf(s string, a ...interface{}) {
	l.Fatal(fmt.Sprintf(s, a...))
}

// Fatale logs Error if not nil then exit
func (l *Logger) Fatale(e error, s string, a ...interface{}) {
	if e != nil {
		l.Fatalf(s, a...)
	}
}

// Sys logs string on output
func (l *Logger) Sys(s string) {
	if l.SysLog {
		log.Println(s)
	}
}

// Sysf logs Sprintf on output
func (l *Logger) Sysf(s string, a ...interface{}) {
	l.Sys(fmt.Sprintf(s, a...))
}

func (l *Logger) log(s string, p string) {
	l.Logger.Println(fmt.Sprintf("[%s] \"%s\"", p, s))
}
