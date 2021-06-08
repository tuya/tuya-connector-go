package logger

import "runtime"

type UNIT int64

const (
	_       = iota
	KB UNIT = 1 << (iota * 10)
	MB
	GB
	TB
)

const (
	LOG int = iota
	DEBUG
	INFO
	WARN
	ERROR
	FATAL
)

const (
	OS_LINUX = iota
	OS_X
	OS_WIN
	OS_OTHERS
)

const (
	TimeFormat = "2006-01-02 15:04:05"
)

var logFormat = "%s [%s:%d] [%s] -> %s"
var maxLogFileSize int64

func getLogLevelName(level int) string {
	switch level {
	case LOG:
		return "LOG"
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

func getOsFlag() int {
	switch os := runtime.GOOS; os {
	case "darwin":
		return OS_X
	case "linux":
		return OS_LINUX
	case "windows":
		return OS_WIN
	default:
		return OS_OTHERS
	}
}

func getOsEol() string {
	if getOsFlag() == OS_WIN {
		return "\r\n"
	}
	return "\n"
}
