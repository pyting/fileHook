package filehook

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type FileCycle int

const (
	YEAR FileCycle = 4 + iota*2
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
}

func NewFileHook(path string, dirNameCycle, fileNameCycle FileCycle) (fileHook *FileHook, err error) {
	format := []rune("20060102150405")
	path = strings.TrimSuffix(path, string(os.PathSeparator))

	dirname := time.Now().Format(string(format[0:dirNameCycle]))
	err = os.MkdirAll(path+string(os.PathSeparator)+dirname, 0755)
	if err != nil {
		return
	}
	err = os.RemoveAll(path + string(os.PathSeparator) + dirname)
	if err != nil {
		return
	}
	fileHook = &FileHook{
		path:              path,
		dirNameFormatter:  string(format[0:dirNameCycle]),
		fileNameFormatter: string(format[0:fileNameCycle]),
		Formatter:         &logrus.JSONFormatter{},
	}
	return
}

func (h *FileHook) CloseConsole() {
	f, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	logrus.SetOutput(f)
}

func (h *FileHook) Fire(entry *logrus.Entry) error {
	pathSeparator := string(os.PathSeparator)
	path := h.path + pathSeparator + time.Now().Format(h.dirNameFormatter)
	filename := h.Prefix + time.Now().Format(h.fileNameFormatter) + h.Suffix
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
		h.writer, err = os.OpenFile(path+pathSeparator+filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to open file, %v", err)
			return err
		}
		h.tmpfile = filename
	}

	_, err = h.writer.Write(line)
	if err != nil {
		panic(err)
		fmt.Fprintf(os.Stderr, "Unable to write file, %v", err)
		return err
	}
	return nil
}

func (h *FileHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
