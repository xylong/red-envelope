package base

import (
	"fmt"
	"github.com/go-playground/locales/zh"
	"github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh2 "github.com/go-playground/validator/v10/translations/zh"
	"github.com/sirupsen/logrus"
	"red-envelope/infra"
)

var (
	validate   *validator.Validate
	translator ut.Translator
)

func Validate() *validator.Validate {
	return validate
}

func Translate() ut.Translator {
	return translator
}

type ValidatorStarter struct {
	infra.BaseStarter
}

func (v *ValidatorStarter) Init(ctx infra.StarterContext) {
	validate = validator.New()
	cn := zh.New()
	uni := ut.New(cn, cn)
	var ok bool
	translator, ok = uni.GetTranslator("zh")
	if ok {
		err := zh2.RegisterDefaultTranslations(validate, trans)
		if err != nil {
			logrus.Error(err)
		}
	} else {
		logrus.Error("not found translator: zh")
	}
}
