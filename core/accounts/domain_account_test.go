package accounts

import (
	"github.com/segmentio/ksuid"
	"github.com/shopspring/decimal"
	"github.com/smartystreets/goconvey/convey"
	"red-envelope/services"
	"testing"
)

func TestAccountDomain_Create(t *testing.T) {
	dto := services.AccountDTO{
		UserId:   ksuid.New().Next().String(),
		Username: "测试用户",
		Balance:  decimal.NewFromFloat(0),
		Status:   1,
	}
	domain := new(accountDomain)
	convey.Convey("创建账户", t, func() {
		rdto, err := domain.Create(dto)
		convey.So(err, convey.ShouldBeNil)
		convey.So(rdto, convey.ShouldNotBeNil)
		convey.So(rdto.Balance.String(), convey.ShouldEqual, dto.Balance.String())
		convey.So(rdto.UserId, convey.ShouldEqual, dto.UserId)
		convey.So(rdto.Username, convey.ShouldEqual, dto.Username)
		convey.So(rdto.Status, convey.ShouldEqual, dto.Status)
	})
}

func TestAccountDomain_Transfer(t *testing.T) {
	// 两个账户，交易账户主题要有余额
	account1, account2 := &services.AccountDTO{
		UserId:   ksuid.New().Next().String(),
		Username: "张三",
		Balance:  decimal.NewFromFloat(100),
		Status:   1,
	}, &services.AccountDTO{
		UserId:   ksuid.New().Next().String(),
		Username: "李四",
		Balance:  decimal.NewFromFloat(100),
		Status:   1,
	}
	domain := accountDomain{}
	convey.Convey("转账测试", t, func() {
		// 创建账户1
		a1, err := domain.Create(*account1)
		convey.So(err, convey.ShouldBeNil)
		convey.So(a1, convey.ShouldNotBeNil)
		convey.So(a1.Balance.String(), convey.ShouldEqual, account1.Balance.String())
		convey.So(a1.UserId, convey.ShouldEqual, account1.UserId)
		convey.So(a1.Username, convey.ShouldEqual, account1.Username)
		convey.So(a1.Status, convey.ShouldEqual, account1.Status)
		account1 = a1
		// 创建账户2
		a2, err := domain.Create(*account2)
		convey.So(err, convey.ShouldBeNil)
		convey.So(a2, convey.ShouldNotBeNil)
		convey.So(a2.Balance.String(), convey.ShouldEqual, account1.Balance.String())
		convey.So(a2.UserId, convey.ShouldEqual, account2.UserId)
		convey.So(a2.Username, convey.ShouldEqual, account2.Username)
		convey.So(a2.Status, convey.ShouldEqual, account2.Status)
		account2 = a2
		// 1.余额充足，资金转入其他账户
		convey.Convey("余额充足，给其他账户转账", func() {
			amount := decimal.NewFromFloat(1)
			body, target := services.TradeParticipator{
				AccountNo: account1.AccountNo,
				UserId:    account1.UserId,
				Username:  account1.Username,
			}, services.TradeParticipator{
				AccountNo: account2.AccountNo,
				UserId:    account2.UserId,
				Username:  account2.Username,
			}
			dto := services.AccountTransferDTO{
				TradeBody:   body,
				TradeTarget: target,
				TradeNo:     ksuid.New().Next().String(),
				Amount:      amount,
				ChangeType:  services.ChangeType(-1),
				ChangeFlag:  services.FlagTransferOut,
				Decs:        "转出",
			}
			status, err := domain.Transfer(dto)
			convey.So(err, convey.ShouldBeNil)
			convey.So(status, convey.ShouldEqual, services.TransferedStatusSuccess)
			// 实际余额更新后的预期值
			account := domain.GetAccount(account1.AccountNo)
			convey.So(account, convey.ShouldNotBeNil)
			convey.So(account.Balance.String(), convey.ShouldEqual, account1.Balance.Sub(amount).String())
		})
		// 2.余额不足，资金转出
		convey.Convey("余额不足，资金转出", func() {
			amount := account1.Balance
			amount = amount.Add(decimal.NewFromFloat(200))
			body, target := services.TradeParticipator{
				AccountNo: account1.AccountNo,
				UserId:    account1.UserId,
				Username:  account1.Username,
			}, services.TradeParticipator{
				AccountNo: account2.AccountNo,
				UserId:    account2.UserId,
				Username:  account2.Username,
			}
			dto := services.AccountTransferDTO{
				TradeBody:   body,
				TradeTarget: target,
				TradeNo:     ksuid.New().Next().String(),
				Amount:      amount,
				ChangeType:  services.ChangeType(-1),
				ChangeFlag:  services.FlagTransferOut,
				Decs:        "转账",
			}
			status, err := domain.Transfer(dto)
			convey.So(err, convey.ShouldNotBeNil)
			convey.So(status, convey.ShouldEqual, services.TransferedStatusSufficientFunds)
			//实际余额更新后的预期值
			account := domain.GetAccount(account1.AccountNo)
			convey.So(account, convey.ShouldNotBeNil)
			convey.So(account.Balance.String(), convey.ShouldEqual, account1.Balance.String())
		})
		// 3.充值
	})
}
