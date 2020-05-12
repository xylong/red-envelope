package base

import (
	"fmt"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/mattn/go-colorable"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"github.com/tietang/go-utils"
	"github.com/tietang/props/kvs"
	"github.com/x-cray/logrus-prefixed-formatter"
	"io"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var (
	formatter *prefixed.TextFormatter
	lfh       *utils.LineNumLogrusHook
	// 将滚动日志writer共享给iris glog output
	log_writer io.Writer
)

func init() {
	// 定义日志格式
	formatter = &prefixed.TextFormatter{}
	// 设置高亮显示的色彩样式
	formatter.ForceColors = true
	formatter.DisableColors = false
	formatter.ForceFormatting = true
	formatter.SetColorScheme(&prefixed.ColorScheme{
		InfoLevelStyle:  "green",
		WarnLevelStyle:  "yellow",
		ErrorLevelStyle: "red",
		FatalLevelStyle: "41",
		PanicLevelStyle: "41",
		DebugLevelStyle: "blue",
		PrefixStyle:     "cyan",
		TimestampStyle:  "37",
	})
	// 开启完整时间戳输出和时间戳格式
	formatter.FullTimestamp = true
	// 设置时间格式
	formatter.TimestampFormat = "2006-01-02.15:04:05.000000"
	// 设置日志formatter
	logrus.SetFormatter(formatter)
	logrus.SetOutput(colorable.NewColorableStdout())
	// 定义日志级别
	if os.Getenv("log.debug") == "true" {
		logrus.SetLevel(logrus.DebugLevel)
	}
	// 开启调用函数、文件、代码行信息的输出
	logrus.SetReportCaller(true)
	SetLineNumLogrusHook()
}

func SetLineNumLogrusHook() {
	lfh = utils.NewLineNumLogrusHook()
	lfh.EnableFileNameLog = true
	lfh.EnableFuncNameLog = true
	logrus.AddHook(lfh)
}

// InitLog 初始化log配置，配置logrus日志文件滚动生成
func InitLog(source kvs.ConfigSource) {
	// 设置日志输出级别
	level, err := logrus.ParseLevel(source.GetDefault("log.level", "info"))
	if err != nil {
		level = logrus.InfoLevel
	}
	logrus.SetLevel(level)
	if source.GetBoolDefault("log.enableLineLog", true) {
		lfh.EnableFuncNameLog = true
		lfh.EnableFuncNameLog = true
	} else {
		lfh.EnableFuncNameLog = false
		lfh.EnableFileNameLog = false
	}
	// 配置日志输出目录
	logDir := source.GetDefault("log.dir", "./logs")
	logTestDir, err := source.Get("log.test.dir")
	if err == nil {
		logDir = logTestDir
	}
	logPath := logDir
	logFilePath, _ := filepath.Abs(logPath)
	logrus.Infof("log dir: %s", logFilePath)
	logFileName := source.GetDefault("log.file.name", "red-envelop")
	maxAge := source.GetDurationDefault("log.max.age", time.Hour*24)
	rotationTime := source.GetDurationDefault("log.rotation.time", time.Hour*1)
	os.MkdirAll(logPath, os.ModePerm)

	baseLogPath := path.Join(logPath, logFileName)
	// 设置滚动日志输出writer
	writer, err := rotatelogs.New(
		strings.TrimSuffix(baseLogPath, ".log")+".%Y%m%d%H.log",
		rotatelogs.WithLinkName(baseLogPath),      // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(maxAge),             // 文件最大保存时间
		rotatelogs.WithRotationTime(rotationTime), // 日志切割时间间隔
	)
	if err != nil {
		logrus.Errorf("config local file system logger error. %+v", err)
	}

	// 设置日志文件输出的日志格式
	formatter := &logrus.TextFormatter{}
	formatter.CallerPrettyfier = func(frame *runtime.Frame) (function string, file string) {
		function = frame.Function
		dir, filename := path.Split(frame.File)
		f := path.Base(dir)
		return function, fmt.Sprintf("%s/%s:%d", f, filename, frame.Line)
	}
	lfHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: writer,
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	}, formatter)

	logrus.AddHook(lfHook)
	log_writer = writer
}
