package test

import (
	"github.com/tietang/dbx"
	"golang.org/x/tools/go/ssa/interp/testdata/src/fmt"
	"net/http"
)

var db *dbx.Database

func init() {
	settings := dbx.Settings{
		DriverName: "mysql",
		Host:       "127.0.0.1:3306",
		User:       "root",
		Password:   "123456",
		Database:   "red",
		Options: map[string]string{
			"parseTime": "true",
		},
	}
	var err error
	db, err = dbx.Open(settings)
	if err != nil {
		fmt.Println(err)
	}
	db.SetLogging(false)
	db.RegisterTable(&GoodsSigned{}, "goods")
	db.RegisterTable(&GoodsSigned2{}, "red_envelope_goods3")
	db.RegisterTable(&GoodsUnsigned{}, "goods_unsigned")
	pprof()
}

// pprof 分析器
func pprof() {
	go func() {
		fmt.Println(http.ListenAndServe(":16060", nil))
	}()
}
