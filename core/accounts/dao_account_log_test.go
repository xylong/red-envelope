package accounts

import (
	"github.com/segmentio/ksuid"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/smartystreets/goconvey/convey"
	"github.com/tietang/dbx"
	"red-envelope/infra/base"
	"red-envelope/services"
	_ "red-envelope/test/textx"
	"testing"
)

func TestAccountLogDao_GetOne(t *testing.T) {
	err := base.Tx(func(runner *dbx.TxRunner) error {
		dao := &AccountLogDao{
			runner: runner,
		}
		convey.Convey("通过log编号查询流水数据", t, func() {
			a := &AccountLog{
				LogNo:      ksuid.New().Next().String(),
				TradeNo:    ksuid.New().Next().String(),
				Status:     1,
				AccountNo:  ksuid.New().Next().String(),
				UserId:     ksuid.New().Next().String(),
				Username:   "测试用户",
				Amount:     decimal.NewFromFloat(1),
				Balance:    decimal.NewFromFloat(100),
				ChangeFlag: services.FlagAccountCreated,
				ChangeType: services.AccountCreated,
			}
			convey.Convey("通过log_no查询", func() {
				id, err := dao.Insert(a)
				convey.So(err, convey.ShouldBeNil)
				convey.So(id, convey.ShouldBeGreaterThan, 0)
				na := dao.GetOne(a.LogNo)
				convey.So(na, convey.ShouldNotBeNil)
				convey.So(na.Balance.String(), convey.ShouldEqual, a.Balance.String())
				convey.So(na.Amount.String(), convey.ShouldEqual, a.Amount.String())
				convey.So(na.CreatedAt, convey.ShouldNotBeNil)
			})
			convey.Convey("通过trade_no查询", func() {
				id, err := dao.Insert(a)
				convey.So(err, convey.ShouldBeNil)
				convey.So(id, convey.ShouldBeGreaterThan, 0)
				na := dao.GetByTradeNo(a.TradeNo)
				convey.So(na, convey.ShouldNotBeNil)
				convey.So(na.Balance.String(), convey.ShouldEqual, a.Balance.String())
				convey.So(na.Amount.String(), convey.ShouldEqual, a.Amount.String())
				convey.So(na.CreatedAt, convey.ShouldNotBeNil)
			})
		})
		return nil
	})
	if err != nil {
		logrus.Error(err)
	}
}
