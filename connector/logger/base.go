package logger

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/callmegema/tuya-connector-go/connector/constant"
	"github.com/callmegema/tuya-connector-go/connector/env/extension"
	"github.com/callmegema/tuya-connector-go/connector/utils"
)

func init() {
	extension.SetLog(constant.TUYA_LOG, newLogInstance)
	fmt.Println("init log extension......")
}

func newLogInstance() extension.ILogger {
	return NewDefaultLogger("tysdk", false)
}

var Log extension.ILogger

type TyLogger struct {
	log          *RollingFile
	consolePrint bool
	logLevel     int
}

func NewDefaultLogger(appName string, isDebug bool) extension.ILogger {
	return NewLogger(appName, 10, MB, isDebug)
}

func NewLogger(appName string, maxFileSize int64, fileSizeUnit UNIT, isDebug bool) extension.ILogger {
	tyLog := &TyLogger{
		consolePrint: isDebug,
		logLevel:     ERROR,
	}
	if isDebug {
		tyLog.logLevel = DEBUG
	}
	maxLogFileSize = maxFileSize * int64(fileSizeUnit)
	lf := &RollingFile{dir: "logs", filename: appName + ".log", mu: new(sync.Mutex)}
	lf.mu.Lock()
	defer lf.mu.Unlock()
	if !utils.Exists(lf.dir) {
		utils.Mkdir(lf.dir)
	}
	lf.logfile, _ = os.OpenFile(lf.dir+"/"+lf.filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)
	fi, err := lf.logfile.Stat()
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	lf.filesize = fi.Size()
	tyLog.log = lf
	Log = tyLog
	return tyLog
}

func (t *TyLogger) SetLevel(level int) {
	t.logLevel = level
}

func (t *TyLogger) ConsolePrint(msg string) {
	if t.log == nil || t.log.logfile == nil || t.consolePrint {
		log.Print(msg)
	}
}

func (t *TyLogger) Output(level int, v ...interface{}) bool {
	defer func() {
		if err := recover(); err != nil {
			log.Println("err", err)
		}
	}()
	if t.log != nil {
		t.log.mu.Lock()
		defer t.log.mu.Unlock()
	}
	msg := t.ConcatArgs(" ", v...)
	logStr := t.AssemblyMsg(level, msg)
	t.ConsolePrint(logStr)
	if level >= t.logLevel {
		if t.log != nil {
			_, err := t.log.write([]byte(logStr))
			if err != nil {
				log.Println(err.Error())
				return false
			}
		}
	}
	return true
}

// concat args
func (t *TyLogger) ConcatArgs(delimiter string, input ...interface{}) string {
	buffer := bytes.Buffer{}
	l := len(input)
	for i := 0; i < l; i++ {
		buffer.WriteString(fmt.Sprint(input[i]))
		if i < l-1 {
			buffer.WriteString(delimiter)
		}
	}
	return buffer.String()
}

func (t *TyLogger) AssemblyMsg(level int, msg string) string {
	_, fpath, line, _ := runtime.Caller(5)
	dir := filepath.Dir(fpath)
	file := filepath.Base(fpath)
	if getOsFlag() == OS_WIN {
		dir = strings.Replace(dir, "\\", "/", -1)
	}
	logInfo := fmt.Sprintf(logFormat+getOsEol(), time.Now().Format(TimeFormat),
		path.Base(dir)+"/"+file, line, getLogLevelName(level), msg)
	return logInfo
}

// Log LOG
func (t *TyLogger) Log(v ...interface{}) {
	t.Output(LOG, v...)
}

func (t *TyLogger) Logf(format string, v ...interface{}) {
	t.Output(LOG, fmt.Sprintf(format, v...))
}

// Debug DEBUG
func (t *TyLogger) Debug(v ...interface{}) {
	t.Output(DEBUG, v...)
}

func (t *TyLogger) Debugf(format string, v ...interface{}) {
	t.Output(DEBUG, fmt.Sprintf(format, v...))
}

// Info INFO
func (t *TyLogger) Info(v ...interface{}) {
	t.Output(INFO, v...)
}

func (t *TyLogger) Infof(format string, v ...interface{}) {
	t.Output(INFO, fmt.Sprintf(format, v...))
}

// Warn WARN
func (t *TyLogger) Warn(v ...interface{}) {
	t.Output(WARN, v...)
}

func (t *TyLogger) Warnf(format string, v ...interface{}) {
	t.Output(WARN, fmt.Sprintf(format, v...))
}

// Error ERROR
func (t *TyLogger) Error(v ...interface{}) {
	t.Output(ERROR, v...)
}

func (t *TyLogger) Errorf(format string, v ...interface{}) {
	t.Output(ERROR, fmt.Sprintf(format, v...))
}

// Fatal FATAL
func (t *TyLogger) Fatal(v ...interface{}) {
	t.Output(FATAL, v...)
}

func (t *TyLogger) Fatalf(format string, v ...interface{}) {
	t.Output(FATAL, fmt.Sprintf(format, v...))
}
