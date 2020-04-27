package infra

import (
	log "github.com/sirupsen/logrus"
	"github.com/tietang/props/kvs"
	"reflect"
)

const (
	KeyProps = "_conf"
)

// 资源启动器上下文，
// 用来在服务资源初始化、安装、启动和停止的生命周期中变量和对象的传递
type StarterContext map[string]interface{}

func (s StarterContext) Props() kvs.ConfigSource {
	p := s[KeyProps]
	if p == nil {
		panic("配置未初始化")
	}
	return p.(kvs.ConfigSource)
}

// 基础资源启动器接口
type Starter interface {
	// 1.系统启动，初始化基础资源
	Init(StarterContext)
	// 2.系统基础资源的安装
	Setup(StarterContext)
	// 3.启动基础资源
	Start(StarterContext)
	// 启动器是否可阻塞
	StartBlocking() bool
	// 4.资源的停止和销毁
	Stop(StarterContext)
}

type BaseStarter struct {
}

func (b *BaseStarter) Init(ctx StarterContext) {

}

func (b *BaseStarter) Setup(ctx StarterContext) {

}

func (b *BaseStarter) Start(ctx StarterContext) {

}

func (b *BaseStarter) StartBlocking() bool {
	return true
}

func (b *BaseStarter) Stop(ctx StarterContext) {

}

// 服务启动注册器
type starterRegister struct {
	nonBlockingStarters []Starter
	blockingStarters    []Starter
}

// Register 注册启动器
func (s *starterRegister) Register(starter Starter) {
	if starter.StartBlocking() {
		s.blockingStarters = append(s.blockingStarters, starter)
	} else {
		s.nonBlockingStarters = append(s.nonBlockingStarters, starter)
	}
	typ := reflect.TypeOf(starter)
	log.Infof("Register starter: %s", typ.String())
}

func (s *starterRegister) AllStarters() []Starter {
	starters := make([]Starter, 0)
	starters = append(starters, s.nonBlockingStarters...)
	starters = append(starters, s.blockingStarters...)
	return starters
}

var StarterRegister *starterRegister = &starterRegister{}

func Register(starter Starter) {
	StarterRegister.Register(starter)
}

// SystemRun 系统基础资源的启动管理
func SystemRun() {
	ctx := StarterContext{}
	for _, starter := range StarterRegister.AllStarters() {
		starter.Init(ctx)
	}
	for _, starter := range StarterRegister.AllStarters() {
		starter.Setup(ctx)
	}
	for _, starter := range StarterRegister.AllStarters() {
		starter.Start(ctx)
	}
}
