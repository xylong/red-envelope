package main

import (
	"github.com/tietang/props/ini"
	"github.com/tietang/props/kvs"
	_ "red-envelope"
	"red-envelope/infra"
)

func main() {
	file := kvs.GetCurrentFilePath("config.ini", 1)
	conf := ini.NewIniFileCompositeConfigSource(file)
	app := infra.New(conf)
	app.Start()
}
