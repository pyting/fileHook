package main

import (
	"github.com/sirupsen/logrus"
	"strconv"
	"time"
	"github.com/pyting/filehook"
)

func main() {
	// 目录按月分割 文件按分种分割
	fileHook, err := filehook.NewFileHook("/log", filehook.MONTH, filehook.Minute)
	if err != nil {
		panic(err)
	}
	fileHook.MaxSize = 5 * 1024 * 1024
	fileHook.Suffix = ".log"
	fileHook.CloseConsole()
	logrus.SetLevel(logrus.WarnLevel)
	logrus.AddHook(fileHook)

	go do(1 * time.Second)
	go do(1 * time.Second)
	do(1 * time.Second)
}

func do(d time.Duration) {
	for i := 0; true; i++ {
		logrus.Info(strconv.Itoa(i) + "abcdefg")
		logrus.Warn(strconv.Itoa(i) + "abcdefg")
		logrus.Error(strconv.Itoa(i) + "abcdefg")
		//logrus.WithFields(logrus.Fields{
		//	"error": i,
		//}).Error("error")
		time.Sleep(d)
	}
}
