package utils

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

// PathExists 判断文件或目录是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	// 如果error为nil, 则文件或目录存在
	if err == nil {
		return true, nil
	}

	// 如果error为IsNotExist错误，则表示文件或目录不存在
	if os.IsNotExist(err) {
		return false, nil
	}

	// 如果error为其他错误，则传递error
	return false, err
}

// CreateHardLink 创建硬链接
func CreateHardLink(source string) (string, error) {
	hlink := source + ".hlink"
	if err := os.Link(source, hlink); err != nil {
		return "", err
	}

	return hlink, nil
}

// DropHardLink 删除硬链接
func DropHardLink(source string) (string, error) {
	hlink := source + ".hlink"
	if err := os.Remove(source + ".hlink"); err != nil {
		return "", err
	}

	return hlink, nil
}

// CheckHardLink 检查硬链接
// func CheckHardLink(source string) (bool, error) {
// 	src, err := os.Stat(source)
// 	if err != nil {
// 		return false, err
// 	}
//
// 	hlk, err := os.Stat(source + ".hlink")
// 	if err != nil {
// 		return false, err
// 	}
//
// 	if src.Size() == hlk.Size() {
// 		return true, nil
// 	}
//
// 	return false, nil
// }

// LogLog 日志记录
func Log(logfile string, level string, reportcaller bool) (logger *logrus.Logger, file *os.File, err error) {
	// 创建logger实例，可以创建多个
	var log = logrus.New()

	// 设置log打印目的地
	if logfile != "" && logfile != "os.Stdout" {
		// 不为""则打开文件，打印日志到文件
		file, err := os.OpenFile(logfile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
		//defer file.Close()

		if err != nil {
			return nil, nil, err
		}

		log.SetOutput(file)
	} else {
		// 为""则打印到标准输出
		log.SetOutput(os.Stdout)
	}

	// 格式化日志时间
	log.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})

	// 是否打印行号
	log.SetReportCaller(reportcaller)

	// 设置日志级别为Info及以上
	switch strings.ToUpper(level) {
	case "TRACE":
		log.SetLevel(logrus.TraceLevel)
	case "DEBUG":
		log.SetLevel(logrus.DebugLevel)
	case "INFO":
		log.SetLevel(logrus.InfoLevel)
	case "WARN":
		log.SetLevel(logrus.WarnLevel)
	case "ERROR":
		log.SetLevel(logrus.ErrorLevel)
	case "FATAL":
		log.SetLevel(logrus.FatalLevel)
	case "PANIC":
		log.SetLevel(logrus.PanicLevel)
	}

	return log, file, nil
}

// LogLogLog 合并打印日志
// func LogLogLog(file string, text interface{}) error {
// 	logFile, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
// 	defer logFile.Close()
//
// 	if err != nil {
// 		return err
// 	}
//
// 	mw := io.MultiWriter(os.Stdout, logFile)
// 	log.SetOutput(mw)
// 	log.Println(text)
//
// 	return nil
// }
