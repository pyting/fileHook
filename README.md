# local file system hook for [logrus](https://github.com/sirupsen/logrus)

1. 二级目录按时间分割
2. 日志文件按时间分割
3. 日志文件按大小分割

```go
package main

import (
	"github.com/sirupsen/logrus"
	"github.com/pyting/filehook"
	"strconv"
)

func main() {
    fhook,err := filehook.NewFileHook("/log",filehook.MONTH,filehook.HOUR)
    if err != nil {
        panic(err)
    }
    fhook.Suffix = ".log"
    fhook.Size = 5 * 1024 * 1024
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