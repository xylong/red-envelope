package infra

// Initializer 初始化接口
type Initializer interface {
	// Init 对象实例化后的初始化操作
	Init()
}

// InitializerRegister 初始化注册器
type InitializerRegister struct {
	Initializers []Initializer
}

// Register 注册初始化对象
func (i *InitializerRegister) Register(ini Initializer) {
	i.Initializers = append(i.Initializers, ini)
}
