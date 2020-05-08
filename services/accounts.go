package services

import (
	"github.com/shopspring/decimal"
	"time"
)

type AccountService interface {
	CreateAccount(dto AccountCreatedDTO) (*AccountDTO, error)
	Transfer(dto AccountTransferDTO) (TransferedStatus, error)
	StoreValue(dto AccountTransferDTO) (TransferedStatus, error)
	GetEnvelopeAccountByUserId(userId string) *AccountDTO
}

// TradeParticipator 账户交易参与者
type TradeParticipator struct {
	AccountNo string
	UserId    string
	Username  string
}

// AccountTransferDTO 账户转账
type AccountTransferDTO struct {
	TradeNo     string
	TradeBody   TradeParticipator
	TradeTarget TradeParticipator
	AmountStr   string
	Amount      decimal.Decimal `` //交易金额,该交易涉及的金额
	ChangeType  ChangeType
	ChangeFlag  ChangeFlag
	Decs        string
}

// AccountCreatedDTO 账户创建
type AccountCreatedDTO struct {
	UserId       string
	Username     string
	AccountName  string
	AccountType  int
	CurrencyCode string
	Amount       string
	CreatedAt    time.Time
}

// AccountDTO 账户信息
type AccountDTO struct {
	AccountCreatedDTO
	AccountNo string          // 账户编号,账户唯一标识
	Balance   decimal.Decimal // 账户可用余额
	Status    int             // 账户状态，账户状态：0账户初始化，1启用，2停用
	UpdatedAt time.Time
}
