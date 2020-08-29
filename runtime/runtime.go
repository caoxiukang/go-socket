package runtime

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

var (
	Trace   *log.Logger // 记录所有日志
	Info    *log.Logger // 重要的信息
	Warning *log.Logger // 需要注意的信息
	Error   *log.Logger // 非常严重的问题
)

func init() {
	now := time.Now()

	//时间格式化输出 Printf输出
	fileName := fmt.Sprintf("runtime/error_log_%d%d%d.log", now.Year(), now.Month(), now.Day())

	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open error log file:", err)
	}

	Trace = log.New(io.MultiWriter(file, os.Stderr), "TRACE: ", log.Ldate|log.Ltime|log.Lshortfile)

	Info = log.New(io.MultiWriter(file, os.Stderr), "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(io.MultiWriter(file, os.Stderr), "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(io.MultiWriter(file, os.Stderr), "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}
