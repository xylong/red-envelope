package base

import (
	"fmt"
	"github.com/tietang/props/kvs"
	"red-envelope/infra"
	"sync"
)

var props kvs.ConfigSource

func Props() kvs.ConfigSource {
	Check(props)
	return props
}

type PropsStart struct {
	infra.BaseStarter
}

func (p *PropsStart) Init(ctx infra.StarterContext) {
	props = ctx.Props()
	fmt.Println("初始化配置...")
}

type SystemAccount struct {
	AccountNo   string
	AccountName string
	UserId      string
	Username    string
}

var systemAccount *SystemAccount
var systemAccountOnce sync.Once

func getSystemAccount() *SystemAccount {
	systemAccountOnce.Do(func() {
		systemAccount = new(SystemAccount)
		err := kvs.Unmarshal(Props(), systemAccount, "system.account")
		if err != nil {
			panic(err)
		}
	})
	return systemAccount
}

func GetEnvelopeActivityLink() string {
	return Props().GetDefault("envelope.link", "/v1/envelope/link")
}

func GetEnvelopeActivityDomain() string {
	return Props().GetDefault("envelope.domain", "http://localhost")
}
