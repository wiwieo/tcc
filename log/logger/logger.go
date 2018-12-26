// Copyright 2012-2015 Apcera Inc. All rights reserved.

//Package logger provides logging facilities for the NATS server
package logger

import (
	"fmt"
	"os"
	"tcc_transaction/log/writer"
	"time"
)

// Logger is the server logger
type Logger struct {
	logger     writer.Writer
	debug      bool
	trace      bool
	warn       bool
	err        bool
	warnLabel  string
	errorLabel string
	fatalLabel string
	debugLabel string
	traceLabel string
}

// NewStdLogger creates a logger with output directed to Stderr
func NewStdLogger(debug, trace, warn, err, colors bool) *Logger {
	l := &Logger{
		debug: debug,
		trace: trace,
		warn:  warn,
		err:   err,
	}

	if colors {
		setColoredLabelFormats(l)
	} else {
		setPlainLabelFormats(l)
	}
	return l
}

func (l *Logger) SetLevelDebug(flg bool) {
	l.debug = flg
}

func (l *Logger) SetLevelTrace(flg bool) {
	l.trace = flg
}

func (l *Logger) SetLevelWarn(flg bool) {
	l.warn = flg
}

func (l *Logger) SetLevelError(flg bool) {
	l.err = flg
}

func (l *Logger) SetPath(path string) {
	l.logger = writer.NewWriter(path, 1<<20)
}

func (l *Logger) Close() {
	l.logger.Close()
}

func setPlainLabelFormats(l *Logger) {
	l.debugLabel = "[DBG]"
	l.traceLabel = "[TRC]"
	l.warnLabel = "[WAR]"
	l.errorLabel = "[ERR]"
	l.fatalLabel = "[FTL]"
}

func setColoredLabelFormats(l *Logger) {
	colorFormat := "[\x1b[%dm%s\x1b[0m]"
	l.debugLabel = fmt.Sprintf(colorFormat, 36, "DBG")
	l.traceLabel = fmt.Sprintf(colorFormat, 33, "TRC")
	l.warnLabel = fmt.Sprintf(colorFormat, 32, "WAR")
	l.errorLabel = fmt.Sprintf(colorFormat, 31, "ERR")
	l.fatalLabel = fmt.Sprintf(colorFormat, 31, "FTL")
}

// Debug logs a debug statement
func (l *Logger) Debug(format string, v ...interface{}) {
	if l.debug {
		l.logger.Write([]byte(getContent(nil, format, l.debugLabel, v...)))
	}
}

func (l *Logger) DebugWithField(head map[string]interface{}, format string, v ...interface{}) {
	if l.debug {
		l.logger.Write([]byte(getContent(head, format, l.debugLabel, v...)))
	}
}

// Trace logs a trace statement
func (l *Logger) Trace(format string, v ...interface{}) {
	if l.trace {
		l.logger.Write([]byte(getContent(nil, format, l.traceLabel, v...)))
	}
}

func (l *Logger) TraceWithField(head map[string]interface{}, format string, v ...interface{}) {
	if l.trace {
		l.logger.Write([]byte(getContent(head, format, l.traceLabel, v...)))
	}
}

// Warning logs a notice statement
func (l *Logger) Warning(format string, v ...interface{}) {
	if l.warn {
		l.logger.Write([]byte(getContent(nil, format, l.warnLabel, v...)))
	}
}

// Warning logs a notice statement
func (l *Logger) WarningWithField(head map[string]interface{}, format string, v ...interface{}) {
	if l.warn {
		l.logger.Write([]byte(getContent(head, format, l.warnLabel, v...)))
	}
}

// Error logs an error statement
func (l *Logger) Error(format string, v ...interface{}) {
	if l.err {
		l.logger.Write([]byte(getContent(nil, format, l.errorLabel, v...)))
	}
}

func (l *Logger) ErrorWithField(head map[string]interface{}, format string, v ...interface{}) {
	if l.err {
		l.logger.Write([]byte(getContent(head, format, l.errorLabel, v...)))
	}
}

// Fatal logs a fatal error
func (l *Logger) Fatal(format string, v ...interface{}) {
	l.logger.Write([]byte(getContent(nil, format, l.fatalLabel, v...)))
	l.exit(1)
}

func (l *Logger) FatalWithField(head map[string]interface{}, format string, v ...interface{}) {
	l.logger.Write([]byte(getContent(head, format, l.fatalLabel, v...)))
	l.exit(1)
}

func (l *Logger) exit(code int) {
	l.Close()
	os.Exit(code)
}

func getContent(head map[string]interface{}, format, label string, v ...interface{}) string {
	if len(head) > 0 {
		return fmt.Sprintf("%s [%+v] %+v %s%s", label, time.Now().Format("2006-01-02 15:04:05.0000"), head, f(format, v...), fmt.Sprintln())
	}
	return fmt.Sprintf("%s [%s] %s%s", label, time.Now().Format("2006-01-02 15:04:05.0000"), f(format, v...), fmt.Sprintln())
}

func f(format string, v ...interface{}) string {
	if len(format) > 0 {
		return fmt.Sprintf(format, v...)
	} else {
		return fmt.Sprint(v...)
	}
}
