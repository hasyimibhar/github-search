package common

import (
	"io"
	"log"
	"strings"

	"github.com/hashicorp/logutils"
)

// Logger defines an interface for logging messages.
type Logger interface {
	// Trace logs a trace message.
	Trace(v ...interface{})

	// Tracef logs a trace message similar to fmt.Printf.
	Tracef(format string, v ...interface{})

	// Info logs an info message.
	Info(v ...interface{})

	// Infof logs an info message similar to fmt.Printf.
	Infof(format string, v ...interface{})

	// Warn logs a warning message.
	Warn(err error)

	// Error logs an error message.
	Error(err error)
}

var validLevels = []logutils.LogLevel{"TRACE", "INFO", "WARN", "ERROR"}

// StandardLogger is an implementation of Logger.
type StandardLogger struct {
	logger *log.Logger
	filter *logutils.LevelFilter
}

// NewStandardLogger creates a StandardLogger which writes to the
// specified writer, and with the specified minimum log level.
func NewStandardLogger(writer io.Writer, minLevel string) *StandardLogger {
	// Default to TRACE
	if !isValidLogLevel(minLevel) {
		minLevel = "TRACE"
	}

	filter := &logutils.LevelFilter{
		Levels:   validLevels,
		MinLevel: logutils.LogLevel(strings.ToUpper(minLevel)),
		Writer:   writer,
	}

	return &StandardLogger{
		logger: log.New(filter, "", log.LstdFlags|log.LUTC),
		filter: filter,
	}
}

// Trace implements the Logger Trace method.
func (l *StandardLogger) Trace(v ...interface{}) {
	l.logger.Println(append([]interface{}{"[TRACE]"}, v...)...)
}

// Tracef implements the Logger Tracef method.
func (l *StandardLogger) Tracef(format string, v ...interface{}) {
	l.logger.Printf("[TRACE] "+format, v...)
}

// Info implements the Logger Info method.
func (l *StandardLogger) Info(v ...interface{}) {
	l.logger.Println(append([]interface{}{"[INFO]"}, v...)...)
}

// Infof implements the Logger Infof method.
func (l *StandardLogger) Infof(format string, v ...interface{}) {
	l.logger.Printf("[INFO] "+format, v...)
}

// Warn implements the Logger Warn method.
func (l *StandardLogger) Warn(err error) {
	l.logger.Println(append([]interface{}{"[WARN]"}, err.Error())...)
}

// Error implements the Logger Error method.
func (l *StandardLogger) Error(err error) {
	l.logger.Println(append([]interface{}{"[ERROR]"}, err.Error())...)
}

// Errorf implements the Logger Errorf method.
func (l *StandardLogger) Errorf(format string, v ...interface{}) {
	l.logger.Printf("[ERROR] "+format, v...)
}

func isValidLogLevel(level string) bool {
	for _, l := range validLevels {
		if strings.ToUpper(level) == string(l) {
			return true
		}
	}
	return false
}
