package logging

import (
	"burrowfs/core/config"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"time"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

var levelNames = map[LogLevel]string{
	DEBUG: "DEBUG",
	INFO:  "INFO",
	WARN:  "WARN",
	ERROR: "ERROR",
	FATAL: "FATAL",
}

type Logger struct {
	module     string
	level      LogLevel
	output     io.Writer
	timeFormat string
	mu         sync.Mutex
}

var (
	defaultLogger *Logger
	globalLevel   = INFO
)

func Init() {
	if config.CONFIG.Debug {
		setGlobalLevel(DEBUG)
	}

	defaultLogger = &Logger{
		module:     "BurroWFS",
		level:      globalLevel,
		output:     os.Stdout,
		timeFormat: "2006-01-02 15:04:05",
	}
}

func Get(module string) *Logger {
	if defaultLogger == nil {
		Init()
	}
	return &Logger{
		module:     module,
		level:      globalLevel,
		output:     defaultLogger.output,
		timeFormat: defaultLogger.timeFormat,
	}
}

func (l *Logger) SetLevel(level LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

func (l *Logger) SetOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.output = w
}

func setGlobalLevel(level LogLevel) {
	globalLevel = level
}

func (l *Logger) setTimeFormat(format string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.timeFormat = format
}

func (l *Logger) log(level LogLevel, msg string, args ...interface{}) {
	if level < l.level {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	_, _, line, _ := runtime.Caller(2)

	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}

	timestamp := time.Now().Format(l.timeFormat)
	logEntry := fmt.Sprintf("[%s] [%s] [%s:%d] %s\n", timestamp, levelNames[level], l.module, line, msg)
	_, err := fmt.Fprint(l.output, logEntry)
	if err != nil {
		panic(err)
	}

	if level == FATAL {
		os.Exit(1)
	}
}

func (l *Logger) Debug(msg string, args ...interface{}) {
	l.log(DEBUG, msg, args...)
}

func (l *Logger) Info(msg string, args ...interface{}) {
	l.log(INFO, msg, args...)
}

func (l *Logger) Warn(msg string, args ...interface{}) {
	l.log(WARN, msg, args...)
}

func (l *Logger) Error(msg string, args ...interface{}) {
	l.log(ERROR, msg, args...)
}

func (l *Logger) Fatal(msg string, args ...interface{}) {
	l.log(FATAL, msg, args...)
}
