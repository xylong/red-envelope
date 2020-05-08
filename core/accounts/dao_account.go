package accounts

import (
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
)

type AccountDao struct {
	runner *dbx.TxRunner
}

func (dao *AccountDao) GetOne(AccountNo string) *Account {
	a := &Account{AccountNo: AccountNo}
	ok, err := dao.runner.GetOne(a)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	if !ok {
		return nil
	}
	return a
}

func (dao *AccountDao) GetByUserId(userId string, accountType int) *Account {
	a := &Account{}
	sql := "select * from account where user_id=? and account_type=?"
	ok, err := dao.runner.Get(a, sql, userId, accountType)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	if !ok {
		return nil
	}
	return a
}

func (dao *AccountDao) Insert(a *Account) (id int64, err error) {
	r, err := dao.runner.Insert(a)
	if err != nil {
		return 0, err
	}
	return r.LastInsertId()
}

// UpdateBalance 账户余额更新
// amount是正数就是增加，是负数就是扣减
func (dao *AccountDao) UpdateBalance(accountNo string, amount decimal.Decimal) (rows int64, err error) {
	// 乐观锁判断余额一定够
	sql := "update account set balance=balance + cast(? as decimal(30, 6)) where account_no=? and balance>=-1*cast(? as decimal(30, 6))"
	r, err := dao.runner.Exec(sql, amount.String(), accountNo, amount.String())
	if err != nil {
		return 0, err
	}
	return r.RowsAffected()
}

// 账户状态更新
func (dao *AccountDao) UpdateStatus(accountNo string, status int) (rows int64, err error) {
	sql := "update account set status=? where account_no=?"
	r, err := dao.runner.Exec(sql, status, accountNo)
	if err != nil {
		return 0, err
	}
	return r.RowsAffected()
}
