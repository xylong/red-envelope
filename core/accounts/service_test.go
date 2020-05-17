package accounts

import (
	"github.com/segmentio/ksuid"
	"github.com/shopspring/decimal"
	"github.com/smartystreets/goconvey/convey"
	"red-envelope/services"
	"testing"
)

func TestAccountService_CreateAccount(t *testing.T) {
	dto := services.AccountCreatedDTO{
		UserId:       ksuid.New().Next().String(),
		Username:     "测试用户",
		Amount:       "100",
		AccountName:  "测试账户",
		AccountType:  2,
		CurrencyCode: "CNY",
	}
	service := new(accountService)
	convey.Convey("创建账户", t, func() {
		account, err := service.CreateAccount(dto)
		convey.So(err, convey.ShouldBeNil)
		convey.So(account, convey.ShouldNotBeNil)
		convey.So(account.Balance.String(), convey.ShouldEqual, dto.Amount)
		convey.So(account.UserId, convey.ShouldEqual, dto.UserId)
		convey.So(account.Username, convey.ShouldEqual, dto.Username)
		convey.So(account.Status, convey.ShouldEqual, 1)
	})
}

func TestAccountService_Transfer(t *testing.T) {
	convey.Convey("转账", t, func() {
		a1, a2 := services.AccountCreatedDTO{
			UserId:       ksuid.New().Next().String(),
			Username:     "王五",
			Amount:       "100",
			AccountName:  "测试账户1",
			AccountType:  2,
			CurrencyCode: "CNY",
		}, services.AccountCreatedDTO{
			UserId:       ksuid.New().Next().String(),
			Username:     "赵六",
			Amount:       "100",
			AccountName:  "测试账户2",
			AccountType:  2,
			CurrencyCode: "CNY",
		}
		service := new(accountService)
		account1, err := service.CreateAccount(a1)
		convey.So(err, convey.ShouldBeNil)
		convey.So(account1, convey.ShouldNotBeNil)
		account2, err := service.CreateAccount(a2)
		convey.So(err, convey.ShouldBeNil)
		convey.So(account2, convey.ShouldNotBeNil)

		convey.Convey("余额足够", func() {
			body, target := services.TradeParticipator{
				AccountNo: account1.AccountNo,
				UserId:    account1.UserId,
				Username:  account1.Username,
			}, services.TradeParticipator{
				AccountNo: account2.AccountNo,
				UserId:    account2.UserId,
				Username:  account2.Username,
			}

			amount := decimal.NewFromFloat(10)
			dto := services.AccountTransferDTO{
				TradeBody:   body,
				TradeTarget: target,
				TradeNo:     ksuid.New().Next().String(),
				AmountStr:   amount.String(),
				ChangeType:  services.ChangeType(-1),
				ChangeFlag:  services.FlagTransferOut,
				Decs:        "转出",
			}
			status, err := service.Transfer(dto)
			convey.So(err, convey.ShouldBeNil)
			convey.So(status, convey.ShouldEqual, services.TransferedStatusSuccess)
			account := service.GetAccount(account1.AccountNo)
			convey.So(account, convey.ShouldNotBeNil)
			convey.So(account.Balance.String(), convey.ShouldEqual, account1.Balance.Sub(amount).String())
		})

		convey.Convey("余额不足", func() {
			body, target := services.TradeParticipator{
				AccountNo: account1.AccountNo,
				UserId:    account1.UserId,
				Username:  account1.Username,
			}, services.TradeParticipator{
				AccountNo: account2.AccountNo,
				UserId:    account2.UserId,
				Username:  account2.Username,
			}

			amount := account1.Balance.Add(decimal.NewFromFloat(200))
			dto := services.AccountTransferDTO{
				TradeBody:   body,
				TradeTarget: target,
				TradeNo:     ksuid.New().Next().String(),
				AmountStr:   amount.String(),
				ChangeType:  services.ChangeType(-1),
				ChangeFlag:  services.FlagTransferOut,
				Decs:        "转出",
			}
			status, err := service.Transfer(dto)
			convey.So(err, convey.ShouldNotBeNil)
			convey.So(status, convey.ShouldEqual, services.TransferedStatusSufficientFunds)
			account := service.GetAccount(account1.AccountNo)
			convey.So(account, convey.ShouldNotBeNil)
			convey.So(account.Balance.String(), convey.ShouldEqual, account1.Balance.String())
		})

		convey.Convey("储值", func() {
			body := services.TradeParticipator{
				AccountNo: account1.AccountNo,
				UserId:    account1.UserId,
				Username:  account1.Username,
			}
			target := body
			amount := decimal.NewFromFloat(10)
			dto := services.AccountTransferDTO{
				TradeBody:   body,
				TradeTarget: target,
				TradeNo:     ksuid.New().Next().String(),
				AmountStr:   amount.String(),
				ChangeType:  services.AccountStoreValue,
				ChangeFlag:  services.FlagTransferIn,
				Decs:        "储值",
			}
			status, err := service.Transfer(dto)
			convey.So(err, convey.ShouldBeNil)
			convey.So(status, convey.ShouldEqual, services.TransferedStatusSuccess)

			account := service.GetAccount(account1.AccountNo)
			convey.So(account, convey.ShouldNotBeNil)
			convey.So(account.Balance.String(), convey.ShouldEqual, account1.Balance.Add(amount).String())
		})
	})
}
