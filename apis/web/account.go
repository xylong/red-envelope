package web

import (
	"github.com/kataras/iris"
	"github.com/sirupsen/logrus"
	"red-envelope/infra"
	"red-envelope/infra/base"
	"red-envelope/services"
)

func init() {
	infra.RegisterApi(new(AccountApi))
}

type AccountApi struct {
}

func (a *AccountApi) Init() {
	groupRouter := base.Iris().Party("/v1/account")
	groupRouter.Post("create", createHandler)
	groupRouter.Post("/transfer", transferHandler)
	groupRouter.Get("/envelope/get", getEnvelopeAccountHandler)
	groupRouter.Get("/get", getAccountHandler)
}

// createHandler 账户创建的接口: /v1/account/create
func createHandler(ctx iris.Context) {
	account := services.AccountCreatedDTO{}
	err := ctx.ReadJSON(&account)
	r := base.Res{
		Code: base.ResCodeOk,
	}
	if err != nil {
		r.Code = base.ResCodeRequestParamsError
		r.Message = err.Error()
		ctx.JSON(&r)
		return
	}
	service := services.GetAccountService()
	dto, err := service.CreateAccount(account)
	if err != nil {
		r.Code = base.ResCodeInnerServerError
		r.Message = err.Error()
	}
	r.Data = dto
	ctx.JSON(&r)
}

// transferHandler 转账 :/v1/account/transfer
func transferHandler(ctx iris.Context) {
	//获取请求参数，
	account := services.AccountTransferDTO{}
	err := ctx.ReadJSON(&account)
	r := base.Res{
		Code: base.ResCodeOk,
	}
	if err != nil {
		r.Code = base.ResCodeRequestParamsError
		r.Message = err.Error()
		ctx.JSON(&r)
		logrus.Error(err)
		return
	}
	//执行转账逻辑
	service := services.GetAccountService()
	status, err := service.Transfer(account)
	if err != nil {
		r.Code = base.ResCodeInnerServerError
		r.Message = err.Error()
		logrus.Error(err)
	}
	if status != services.TransferedStatusSuccess {
		r.Code = base.ResCodeBizError
		r.Message = err.Error()
	}
	r.Data = status
	ctx.JSON(&r)
}

// getEnvelopeAccountHandler 查询红包账户：/v1/account/envelope/get
func getEnvelopeAccountHandler(ctx iris.Context) {
	userId := ctx.URLParam("userId")
	r := base.Res{
		Code: base.ResCodeOk,
	}
	if userId == "" {
		r.Code = base.ResCodeRequestParamsError
		r.Message = "用户ID不能为空"
		ctx.JSON(&r)
		return
	}
	service := services.GetAccountService()
	account := service.GetEnvelopeAccountByUserId(userId)
	r.Data = account
	ctx.JSON(&r)
}

// getAccountHandler 查询账户信息：/v1/account/get
func getAccountHandler(ctx iris.Context) {
	accountNo := ctx.URLParam("accountNo")
	r := base.Res{
		Code: base.ResCodeOk,
	}
	if accountNo == "" {
		r.Code = base.ResCodeRequestParamsError
		r.Message = "账户编号不能为空"
		ctx.JSON(&r)
		return
	}
	service := services.GetAccountService()
	account := service.GetAccount(accountNo)
	r.Data = account
	ctx.JSON(&r)
}
