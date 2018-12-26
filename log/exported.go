package log

import (
	"strings"
	"tcc_transaction/log/logger"
)
// 全局，为了避免在上下文中传递对象
var (
	// std is the name of the standard logger in stdlib `log`
	std = New()
)

func New() *logger.Logger {
	return logger.NewStdLogger(false, false, true, true, true)
}

// 凡是使用此log者，必须调用此方法，设置文件输出地址
func SetPath(path string) {
	std.SetPath(path)
}

func Close() {
	std.Close()
}

func SetLevel(level string) {
	switch strings.ToUpper(level) {
	case "DEBUG":
		std.SetLevelDebug(true)
		std.SetLevelTrace(true)
		std.SetLevelWarn(true)
		std.SetLevelError(true)
	case "WARN":
		std.SetLevelDebug(false)
		std.SetLevelTrace(false)
		std.SetLevelWarn(true)
		std.SetLevelError(true)
	case "ERROR":
		std.SetLevelDebug(false)
		std.SetLevelTrace(false)
		std.SetLevelWarn(false)
		std.SetLevelError(true)
	default:
		std.SetLevelDebug(false)
		std.SetLevelTrace(true)
		std.SetLevelWarn(true)
		std.SetLevelError(true)
	}
}

// Tracef logs a message at level Trace on the standard logger.
func Info(args ...interface{}) {
	std.Trace("", args...)
}

func Infof(format string, args ...interface{}) {
	std.Trace(format, args...)
}

func InfofWithField(head map[string]interface{}, format string, args ...interface{}) {
	std.TraceWithField(head, format, args...)
}

// Tracef logs a message at level Trace on the standard logger.
func Trace(args ...interface{}) {
	std.Trace("", args...)
}

func Tracef(format string, args ...interface{}) {
	std.Trace(format, args...)
}

func TracefWithField(head map[string]interface{}, format string, args ...interface{}) {
	std.TraceWithField(head, format, args...)
}

// Debugf logs a message at level Debug on the standard logger.
func Debug(args ...interface{}) {
	std.Debug("", args...)
}

func Debugf(format string, args ...interface{}) {
	std.Debug(format, args...)
}

func DebugfWithField(head map[string]interface{}, format string, args ...interface{}) {
	std.DebugWithField(head, format, args...)
}

// Warnf logs a message at level Warn on the standard logger.
func Warn(args ...interface{}) {
	std.Warning("", args...)
}

func Warnf(format string, args ...interface{}) {
	std.Warning(format, args...)
}

func WarnfWithField(head map[string]interface{}, format string, args ...interface{}) {
	std.WarningWithField(head, format, args...)
}

// Errorf logs a message at level Error on the standard logger.
func Error(args ...interface{}) {
	std.Error("", args...)
}

func Errorf(format string, args ...interface{}) {
	std.Error(format, args...)
}

func ErrorfWithField(head map[string]interface{}, format string, args ...interface{}) {
	std.ErrorWithField(head, format, args...)
}

// Fatalf logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
func Fatal(args ...interface{}) {
	std.Fatal("", args...)
}

func Fatalf(format string, args ...interface{}) {
	std.Fatal(format, args...)
}

func FatalfWithField(head map[string]interface{}, format string, args ...interface{}) {
	std.FatalWithField(head, format, args...)
}
