package filehook

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"errors"
	"path/filepath"
	"syscall"
)

type FileCycle int

const (
	YEAR   FileCycle = 4 + iota*2
	MONTH
	DAY
	HOUR
	Minute
	Second
)

type FileHook struct {
	sync.Mutex
	path              string // 目录路径
	dirNameFormatter  string // 目录名按时间划分
	fileNameFormatter string // 文件名按时间划分
	tmppath           string //
	tmpfile           string //
	Prefix            string // 文件名前缀
	Suffix            string // 文件名后缀
	Formatter         logrus.Formatter
	writer            *os.File
	MaxSize           int64        // 文件大小限制
	Level             logrus.Level // 日志等级
}

func NewFileHook(path string, dirCycle, fileCycle FileCycle) (fileHook *FileHook, err error) {
	if dirCycle >= fileCycle {
		err = errors.New("dirNameCycle must less fileNameCycle.")
		return
	}
	format := []rune("20060102150405")
	path = filepath.Clean(path)
	err = os.MkdirAll(path, os.ModeDir|0755)
	if err != nil {
		return
	}
	var fin os.FileInfo
	fin, err = os.Stat(path)
	if err != nil {
		return
	}
	sys_stat, ok := fin.Sys().(*syscall.Stat_t)
	if !ok {
		err = errors.New("*syscall.Stat_t reflect failed")
		return
	}

	// isowner,isgrouper,isother
	var permIsWRX bool
	if os.Getuid() == int(sys_stat.Uid) {
		// owner 查看高三位
		//fmt.Println("is owner")
		if sys_stat.Mode&0700 == 0700 {
			permIsWRX = true
		}
	} else if os.Getgid() == int(sys_stat.Gid) {
		// grouper 查看中三位
		//fmt.Println("is grouper")
		if sys_stat.Mode&0070 == 0070 {
			permIsWRX = true
		}
	} else {
		// other 查看低三位
		//fmt.Println("is other")
		if sys_stat.Mode&0007 == 0007 {
			permIsWRX = true
		}
	}
	if !permIsWRX {
		err = errors.New(path + " Permission denied")
		return
	}
	fileHook = &FileHook{
		path:              path,
		dirNameFormatter:  string(format[0:dirCycle]),
		fileNameFormatter: string(format[0:fileCycle]),
		Formatter:         &logrus.JSONFormatter{},
		Level:             logrus.InfoLevel,
	}
	return
}

func (h *FileHook) CloseConsole() error {
	f, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	logrus.SetOutput(f)
	return nil
}

func (h *FileHook) Fire(entry *logrus.Entry) error {
	if entry.Level > h.Level {
		return nil
	}
	pathSeparator := string(os.PathSeparator)
	path := h.path + pathSeparator + time.Now().Format(h.dirNameFormatter)
	filename := time.Now().Format(h.fileNameFormatter)

	// 当目录发生变化时才创建目录
	if path != h.tmppath {
		os.MkdirAll(path, 0755)
		h.tmppath = path
	}

	line, err := h.Formatter.Format(entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read entry, %v", err)
		return err
	}
	h.Lock()
	defer h.Unlock()
	// 当文件名发生变化时才重新打开文件
	if filename != h.tmpfile {
		h.writer.Close()
		h.writer, err = os.OpenFile(path+pathSeparator+h.Prefix+time.Now().Format("20060102150405")+h.Suffix,
			os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to open file, %v", err)
			return err
		}
		h.tmpfile = filename
	} else {
		fin, err := h.writer.Stat()
		if err != nil {
			return err
		}
		size := fin.Size()
		if h.MaxSize != 0 && h.MaxSize <= size {
			h.writer.Close()
			h.writer, err = os.OpenFile(path + pathSeparator + h.Prefix+
				time.Now().Format("20060102150405")+ h.Suffix,
				os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Unable to open file, %v", err)
				return err
			}
			h.tmpfile = filename
		}
	}

	_, err = h.writer.Write(line)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to write file, %v", err)
		return err
	}
	return nil
}

func (h *FileHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
