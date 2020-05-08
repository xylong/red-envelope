package services

import "time"

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
	UserName  string
}

// AccountTransferDTO 账户转账
type AccountTransferDTO struct {
	TradeNo     string
	TradeBody   TradeParticipator
	TradeTarget TradeParticipator
	AmountStr   string
	ChangeType  ChangeType
	ChangeFlag  ChangeFlag
	Decs        string
}

// AccountCreatedDTO 账户创建
type AccountCreatedDTO struct {
	UserId       string
	UserName     string
	AccountName  string
	AccountType  int
	CurrencyType string
	Amount       string
}

// AccountDTO 账户信息
type AccountDTO struct {
	AccountCreatedDTO
	AccountNo string
	CreatedAt time.Time
}
