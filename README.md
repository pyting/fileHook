# local file system hook for [logrus](https://github.com/sirupsen/logrus)

1. 二级目录按时间生成，日志文件按时间生成
2. 日志文件按大小分割

```go
package main

import (
	"github.com/sirupsen/logrus"
	"github.com/pyting/filehook"
	"strconv"
)

func main() {
    fhook,err := filehook.NewFileHook("/log",filehook.MONTH,filehook.HOUR)
    if err != nil {//
        panic(err)
    }
    fhook.Suffix = ".log"
    logrus.SetLevel(logrus.InfoLevel)
    logrus.AddHook(fhook)
    
    for i:=0;true;i++{
    	logrus.Info(strconv.Itoa(i)+" info")
    	logrus.Warn(strconv.Itoa(i)+" warn")
        logrus.Error(strconv.Itoa(i)+" error")
        logrus.Fatal(strconv.Itoa(i)+" fatal")
    }
}
```