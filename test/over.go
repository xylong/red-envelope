package test

import (
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/tietang/dbx"
	"net/http"
)

var db *dbx.Database

func init() {
	settings := dbx.Settings{
		DriverName: "mysql",
		Host:       "127.0.0.1:3306",
		User:       "root",
		Password:   "root",
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

// UpdateForLock 事务锁方案
func UpdateForLock(g Goods) {
	err := db.Tx(func(run *dbx.TxRunner) error {
		query := "select * from goods where envelope_no=? for update"
		out := &GoodsSigned{}
		ok, err := run.Get(out, query, g.EnvelopeNo)
		if !ok || err != nil {
			return err
		}
		// 计算金额和剩余数量
		subAmount := decimal.NewFromFloat(0.01)
		remainAmount := out.RemainAmount.Sub(subAmount)
		remainQuantity := out.RemainQuantity - 1

		sql := "update goods set remain_amount=?,remain_quantity=? where envelope_no=?"
		_, row, err := db.Execute(sql, remainAmount, remainQuantity, g.EnvelopeNo)
		if err != nil {
			return err
		}
		if row < 1 {
			return errors.New("库存扣减失败")
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
}

// UpdateForUnsigned 无符号行直接更新方案
func UpdateForUnsigned(goods Goods) {
	sql := "update goods_unsigned set remain_amount=remain_amount-?,remain_quantity=remain_quantity-? where envelope_no=?"
	_, row, err := db.Execute(sql, 0.01, 1, goods.EnvelopeNo)
	if err != nil {
		fmt.Println(err)
	}
	if row < 1 {
		fmt.Println("扣减失败")
	}
}

// UpdateForOptimistic 乐观锁方案
func UpdateForOptimistic(goods Goods) {
	sql := "update goods set remain_amount=remain_amount-?,remain_quantity=remain_quantity-? where envelope_no=? and remain_amount>=? and remain_quantity>=?"
	_, row, err := db.Execute(sql, 0.01, 1, goods.EnvelopeNo, 0.01, 1)
	if err != nil {
		fmt.Println(err)
	}
	if row < 1 {
		fmt.Println("扣减失败")
	}
}

// UpdateForOptimisticUnsigned 乐观锁+无符号方案
func UpdateForOptimisticUnsigned(goods Goods) {
	sql := "update goods_unsigned set remain_amount=remain_amount-?,remain_quantity=remain_quantity-? where envelope_no=? and remain_amount>=? and remain_quantity>=?"
	_, row, err := db.Execute(sql, 0.01, 1, goods.EnvelopeNo, 0.01, 1)
	if err != nil {
		fmt.Println(err)
	}
	if row < 1 {
		fmt.Println("扣减失败")
	}
}
