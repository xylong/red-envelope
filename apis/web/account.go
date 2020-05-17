package web

import (
	"github.com/kataras/iris"
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
	create(groupRouter)
}

func create(groupRouter iris.Party) {
	groupRouter.Post("/create", func(ctx iris.Context) {
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
	})
}
