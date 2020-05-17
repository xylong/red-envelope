package base

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	recover2 "github.com/kataras/iris/middleware/recover"
	"github.com/sirupsen/logrus"
	"red-envelope/infra"
	"time"
)

var irisApplication *iris.Application

func Iris() *iris.Application {
	Check(irisApplication)
	return irisApplication
}

type IrisServerStarter struct {
	infra.BaseStarter
}

func (i *IrisServerStarter) Init(ctx infra.StarterContext) {
	// 创建iris实例
	irisApplication = initIris()
	// 日志组建配置和扩展
	logger := irisApplication.Logger()
	logger.Install(logrus.StandardLogger())
}

func (i *IrisServerStarter) Start(ctx infra.StarterContext) {
	// 和logrus日志级别保持一致
	Iris().Logger().SetLevel(ctx.Props().GetDefault("log.level", "info"))
	// 控制台打印路由
	routes := Iris().GetRoutes()
	// 把路由信息打印到控制台
	for _, route := range routes {
		logrus.Info(route.Trace())
	}
	// 启动
	port := ctx.Props().GetDefault("app.server.port", "8080")
	Iris().Run(iris.Addr(":" + port))
}

func (i *IrisServerStarter) StartBlocking() bool {
	return true
}

func initIris() *iris.Application {
	app := iris.New()
	app.Use(recover2.New())
	conf := logger.Config{
		Status:             true,
		IP:                 true,
		Method:             true,
		Path:               true,
		Query:              true,
		Columns:            true,
		MessageContextKeys: []string{"logger_message"},
		MessageHeaderKeys:  []string{"User-Agent"},
		LogFunc: func(now time.Time, latency time.Duration, status, ip, method, path string, message, headerMessage interface{}) {
			app.Logger().Infof("| %s | %s | %s | %s | %s | %s | %s | %s |",
				now.Format("2016-01-02.15:04:05.000000"),
				latency.String(),
				status,
				ip,
				method,
				path,
				message,
				headerMessage,
			)
		},
	}
	app.Use(logger.New(conf))
	return app
}
