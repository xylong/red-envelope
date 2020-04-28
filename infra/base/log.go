package base

import (
	"github.com/sirupsen/logrus"
	"github.com/x-cray/logrus-prefixed-formatter"
	"os"
)

func init() {
	// 定义日志格式
	formatter := &prefixed.TextFormatter{}
	formatter.FullTimestamp = true
	formatter.TimestampFormat = "2006-01-02.15:04:05.000000"
	formatter.ForceFormatting = true
	logrus.SetFormatter(formatter)
	// 定义日志级别
	level := os.Getenv("log.debug")
	if level == "true" {
		logrus.SetLevel(logrus.DebugLevel)
	}
	// 控制台高亮显示
	formatter.ForceColors = true
	formatter.DisableColors = false
	// 日志文件和滚动配置
}
