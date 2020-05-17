package infra

var apiInitializerRegister *InitializerRegister = new(InitializerRegister)

func RegisterApi(i Initializer) {
	apiInitializerRegister.Register(i)
}

// GetApiInitializers 获取注册的api初始化对象
func GetApiInitializers() []Initializer {
	return apiInitializerRegister.Initializers
}

type WebApiStarter struct {
	BaseStarter
}

func (w *WebApiStarter) Setup(ctx StarterContext) {
	for _, v := range GetApiInitializers() {
		v.Init()
	}
}
