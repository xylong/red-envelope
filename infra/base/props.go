package base

import (
	"fmt"
	"github.com/tietang/props/kvs"
	"red-envelope/infra"
)

var props kvs.ConfigSource

func Props() kvs.ConfigSource {
	return props
}

type PropsStart struct {
	infra.BaseStarter
}

func (p *PropsStart) Init(ctx infra.StarterContext) {
	props = ctx.Props()
	fmt.Println("初始化配置...")
}
