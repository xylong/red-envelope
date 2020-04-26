package test

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/segmentio/ksuid"
	"github.com/shopspring/decimal"
	"testing"
)

func BenchmarkUpdateForLock(b *testing.B) {
	g := GoodsSigned{}
	g.EnvelopeNo = ksuid.New().Next().String()
	g.RemainQuantity = 100000
	g.RemainAmount = decimal.NewFromFloat(100000)
	_, err := db.Insert(g)
	if err != nil {
		fmt.Println(err)
		return
	}
	for i := 0; i < b.N; i++ {
		UpdateForLock(g.Goods)
	}
}
