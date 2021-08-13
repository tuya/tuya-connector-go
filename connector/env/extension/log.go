package extension

// log interface
type ILogger interface {
	SetLevel(level int)
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Warn(v ...interface{})
	Warnf(format string, v ...interface{})
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
}

var (
	logs = make(map[string]ILogger)
)

func SetLog(name string, v func() ILogger) {
	logs[name] = v()
}

func GetLog(name string) ILogger {
	if logs[name] == nil {
		panic("log for " + name + " is not existing, make sure you have import the package.")
	}
	return logs[name]
}
