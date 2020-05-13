package accounts

import (
	"errors"
	"github.com/segmentio/ksuid"
	"github.com/tietang/dbx"
	"red-envelope/infra/base"
	"red-envelope/services"
)

type accountDomain struct {
	account    Account
	accountLog AccountLog
}

// createAccountLogNo 创建logNo
func (d *accountDomain) createAccountLogNo() {
	d.accountLog.LogNo = ksuid.New().Next().String()
}

// createAccountNo 生成accountNo
func (d *accountDomain) createAccountNo() {
	d.account.AccountNo = ksuid.New().Next().String()
}

// createAccountLog 创建流水记录
func (d *accountDomain) createAccountLog() {
	// 通过account来创建流水，创建账户逻辑在前
	d.accountLog = AccountLog{}
	d.createAccountLogNo()
	d.accountLog.TradeNo = d.accountLog.LogNo
	// 流水中的交易主体信息
	d.accountLog.AccountNo = d.account.AccountNo
	d.accountLog.UserId = d.account.UserId
	d.accountLog.Username = d.account.Username.String
	// 交易对象信息
	d.accountLog.TargetAccountNo = d.account.AccountNo
	d.accountLog.TargetUserId = d.account.UserId
	d.accountLog.Username = d.account.Username.String
	// 交易金额
	d.accountLog.Amount = d.account.Balance
	d.accountLog.Balance = d.account.Balance
	// 交易变化属性
	d.accountLog.Decs = "创建账户"
	d.accountLog.ChangeType = services.AccountCreated
	d.accountLog.ChangeFlag = services.FlagAccountCreated
}

func (d *accountDomain) Create(dto services.AccountDTO) (*services.AccountDTO, error) {
	// 创建账户持久化对象
	d.account = Account{}
	d.account.FromDTO(&dto)
	d.createAccountNo()
	d.account.Username.Valid = true
	// 创建账户流水持久化对象
	d.createAccountLog()
	accountDao := AccountDao{}
	accountLogDao := AccountLogDao{}
	var rdto *services.AccountDTO
	err := base.Tx(func(runner *dbx.TxRunner) error {
		accountDao.runner = runner
		accountLogDao.runner = runner
		// 插入账户数据
		id, err := accountDao.Insert(&d.account)
		if err != nil {
			return err
		}
		if id < 0 {
			return errors.New("账户创建失败")
		}
		// 插入流水数据
		id, err = accountLogDao.Insert(&d.accountLog)
		if err != nil {
			return err
		}
		if id < 0 {
			return errors.New("账户流水创建失败")
		}
		d.account = *accountDao.GetOne(d.account.AccountNo)
		return nil
	})
	rdto = d.account.ToDTO()
	return rdto, err
}
