package accounts

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"red-envelope/infra/base"
	"red-envelope/services"
)

var _ services.AccountService = new(accountService)

type accountService struct {
}

func (a *accountService) CreateAccount(dto services.AccountCreatedDTO) (*services.AccountDTO, error) {
	domain := accountDomain{}
	err := base.Validate().Struct(&dto)
	if err != nil {
		_, ok := err.(*validator.InvalidValidationError)
		if ok {
			logrus.Error("验证错误", err)
		}
		errs, ok := err.(validator.ValidationErrors)
		if ok {
			for _, e := range errs {
				logrus.Error(e.Translate(base.Translate()))
			}
		}
		return nil, err
	}

	// 创建账户
	amount, err := decimal.NewFromString(dto.Amount)
	if err != nil {
		return nil, err
	}
	account := services.AccountDTO{
		UserId:       dto.UserId,
		Username:     dto.Username,
		AccountType:  dto.AccountType,
		AccountName:  dto.AccountName,
		CurrencyCode: dto.CurrencyCode,
		Status:       1,
		Balance:      amount,
	}
	return domain.Create(account)
}

func (a *accountService) Transfer(dto services.AccountTransferDTO) (services.TransferedStatus, error) {
	domain := accountDomain{}
	err := base.Validate().Struct(&dto)
	if err != nil {
		_, ok := err.(*validator.InvalidValidationError)
		if ok {
			logrus.Error("验证错误", err)
		}
		errs, ok := err.(validator.ValidationErrors)
		if ok {
			for _, e := range errs {
				logrus.Error(e.Translate(base.Translate()))
			}
		}
		return services.TransferedStatusFailure, err
	}
	// 转账
	amount, err := decimal.NewFromString(dto.AmountStr)
	if err != nil {
		return services.TransferedStatusFailure, err
	}
	dto.Amount = amount
	if dto.ChangeFlag == services.FlagTransferOut {
		if dto.ChangeType > 0 {
			return services.TransferedStatusFailure, errors.New("如果changeFlag为支出，那么changeType必须小于0")
		}
	} else {
		if dto.ChangeType < 0 {
			return services.TransferedStatusFailure, errors.New("如果changeFlag为收入,那么changeType必须大于0")
		}
	}
	status, err := domain.Transfer(dto)
	/*if status == services.TransferedStatusSuccess {
		backwardDto := dto
		backwardDto.TradeBody = dto.TradeTarget
		backwardDto.TradeTarget = dto.TradeBody
		backwardDto.ChangeType = -dto.ChangeType
		backwardDto.ChangeFlag = -dto.ChangeFlag
		status, err := domain.Transfer(backwardDto)
		return status, err
	}*/
	return status, err
}

func (a *accountService) StoreValue(dto services.AccountTransferDTO) (services.TransferedStatus, error) {
	dto.TradeTarget = dto.TradeBody
	dto.ChangeFlag = services.FlagTransferIn
	dto.ChangeType = services.AccountStoreValue
	return a.Transfer(dto)
}

func (a *accountService) GetEnvelopeAccountByUserId(userId string) *services.AccountDTO {
	domain := accountDomain{}
	return domain.GetEnvelopeByUserId(userId)
}

func (a *accountService) GetAccount(accountNo string) *services.AccountDTO {
	domain := accountDomain{}
	return domain.GetAccount(accountNo)
}
