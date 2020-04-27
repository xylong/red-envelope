package infra

import "github.com/tietang/props/kvs"

type BootApplication struct {
	conf           kvs.ConfigSource
	starterContext StarterContext
}

func (b *BootApplication) Start() {
	b.init()
	b.setup()
	b.start()
}

func New(conf kvs.ConfigSource) *BootApplication {
	b := &BootApplication{
		conf:           conf,
		starterContext: StarterContext{},
	}
	b.starterContext[KeyProps] = conf
	return b
}

func (b *BootApplication) init() {
	for _, starter := range StarterRegister.AllStarters() {
		starter.Init(b.starterContext)
	}
}

func (b *BootApplication) setup() {
	for _, starter := range StarterRegister.AllStarters() {
		starter.Setup(b.starterContext)
	}
}

func (b *BootApplication) start() {
	for index, starter := range StarterRegister.AllStarters() {
		if starter.StartBlocking() {
			if index+1 == len(StarterRegister.AllStarters()) {
				starter.Start(b.starterContext)
			} else {
				go starter.Start(b.starterContext)
			}
		} else {
			starter.Start(b.starterContext)
		}
	}
}
