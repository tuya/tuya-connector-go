package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

type RollingFile struct {
	dir      string
	filename string
	filesize int64
	mu       *sync.Mutex
	logfile  *os.File
	_suffix  int
	_date    *time.Time
}

func (f *RollingFile) isMustRename() bool {
	if maxLogFileSize > 0 && f.filesize > maxLogFileSize {
		return true
	}
	return false
}

func (f *RollingFile) rename() {
	if maxLogFileSize > 0 && f.filesize > maxLogFileSize {
		err := f.rotate()
		if err != nil {
			fmt.Println("333", err.Error())
		}
	}
}

func (f *RollingFile) write(data []byte) (int, error) {
	n, err := f.logfile.Write(data)
	if err != nil {
		return n, err
	}
	f.filesize += int64(n)
	f.rename()
	return n, err
}

func (f *RollingFile) removeFiles() {
	fs, err := filepath.Glob(fmt.Sprintf("%s/%s.*", f.dir, f.filename))
	if err != nil {
		return
	}
	sort.Strings(fs)
}

func (f *RollingFile) rotate() error {
	f.removeFiles()
	if f != nil && f.logfile != nil {
		f.logfile.Sync()
		f.logfile.Close()
		os.Rename(f.dir+"/"+f.filename, f.dir+"/"+f.filename+time.Now().Format(".20060101090909"))
	}
	fd, err := os.OpenFile(f.dir+"/"+f.filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	fi, err := fd.Stat()
	if err != nil {
		return err
	}
	f.logfile = fd
	f.filesize = fi.Size()
	return nil
}
