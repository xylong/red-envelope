package textx

import (
	"github.com/tietang/props/ini"
	"github.com/tietang/props/kvs"
	"red-envelope/infra"
	"red-envelope/infra/base"
)

func init() {
	file := kvs.GetCurrentFilePath("../../brun/config.ini", 1)
	conf := ini.NewIniFileCompositeConfigSource(file)
	base.InitLog(conf)

	infra.Register(&base.PropsStart{})
	infra.Register(&base.DbxDatabaseStarter{})
	infra.Register(&base.ValidatorStarter{})

	app := infra.New(conf)
	app.Start()

}
