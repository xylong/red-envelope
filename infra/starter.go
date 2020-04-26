package infra

// 基础资源上下结构体
type StarterContext map[string]interface{}

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

// 注册器
type starterRegister struct {
	starters []Starter
}

// Register 注册启动器
func (s *starterRegister) Register(starter Starter) {
	s.starters = append(s.starters, starter)
}

func (s *starterRegister) AllStarters() []Starter {
	return s.starters
}

var StarterRegister *starterRegister = new(starterRegister)

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
