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
