package envelopes

import (
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
	"red-envelope/services"
	"time"
)

type RedEnvelopeGoodsDao struct {
	runner *dbx.TxRunner
}

// Insert 写入
func (dao *RedEnvelopeGoodsDao) Insert(goods RedEnvelopeGoods) (int64, error) {
	res, err := dao.runner.Insert(goods)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// GetOne 根据红包编号查询红包
func (dao *RedEnvelopeGoodsDao) GetOne(envelopeNo string) *RedEnvelopeGoods {
	good := &RedEnvelopeGoods{EnvelopeNo: envelopeNo}
	ok, err := dao.runner.GetOne(good)
	if err != nil || !ok {
		logrus.Error(err)
		return nil
	}
	return good
}

// UpdateBalance 更新红包数量和余额
func (dao *RedEnvelopeGoodsDao) UpdateBalance(envelopeNo string, amount decimal.Decimal) (int64, error) {
	sql := "update red_envelope_goods set remain_amount=remain_amount-CAST(? AS DECIMAL(30,6)),remain_quantity=remain_quantity-1" +
		" where envelope_no=? " +
		// 乐观锁代替事务行锁
		"and remain_quantity>0 and remain_amount >= CAST(? AS DECIMAL(30,6))"
	res, err := dao.runner.Exec(sql, amount.String(), envelopeNo, amount.String())
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// UpdateOrderStatus 更新订单状态
func (dao *RedEnvelopeGoodsDao) UpdateOrderStatus(envelopeNo string, status services.OrderStatus) (int64, error) {
	sql := "update red_envelope_goods set status=? where envelope_no=?"
	res, err := dao.runner.Exec(sql)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// FindExpired 查询过期红包
func (dao *RedEnvelopeGoodsDao) FindExpired(offset, size int) []RedEnvelopeGoods {
	var goods []RedEnvelopeGoods
	sql := "select * from red_envelope_goods where remain_quantity>0 and expired_at>? and status<>4 limit ?,?"
	err := dao.runner.Find(&goods, sql, time.Now(), offset, size)
	if err != nil {
		logrus.Error(err)
	}
	return goods
}
