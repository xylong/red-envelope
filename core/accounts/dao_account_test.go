package accounts

import (
	"database/sql"
	"github.com/segmentio/ksuid"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/smartystreets/goconvey/convey"
	"github.com/tietang/dbx"
	"red-envelope/infra/base"
	_ "red-envelope/test/textx"
	"testing"
)

func TestAccountDao_GetOne(t *testing.T) {
	err := base.Tx(func(runner *dbx.TxRunner) error {
		dao := &AccountDao{
			runner: runner,
		}
		convey.Convey("通过编号查询账号数据", t, func() {
			a := &Account{
				Balance:     decimal.NewFromFloat(100),
				Status:      1,
				AccountNo:   ksuid.New().Next().String(),
				AccountName: "测试资金账户",
				UserId:      ksuid.New().Next().String(),
				Username: sql.NullString{
					String: "测试用户",
					Valid:  true,
				},
			}
			id, err := dao.Insert(a)
			convey.So(err, convey.ShouldBeNil)
			convey.So(id, convey.ShouldBeGreaterThan, 0)
			na := dao.GetOne(a.AccountNo)
			convey.So(na, convey.ShouldNotBeNil)
			convey.So(na.Balance.String(), convey.ShouldEqual, a.Balance.String())
			convey.So(na.CreatedAt, convey.ShouldNotBeNil)
			convey.So(na.UpdatedAt, convey.ShouldNotBeNil)
		})
		return nil
	})
	if err != nil {
		logrus.Error(err)
	}
}
