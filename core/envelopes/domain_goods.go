package envelopes

import (
	"context"
	"github.com/segmentio/ksuid"
	"github.com/shopspring/decimal"
	"github.com/tietang/dbx"
	"red-envelope/infra/base"
	"red-envelope/services"
	"time"
)

type goodsDomain struct {
	RedEnvelopeGoods
}

// createEnvelopeNo 生成红包编号
func (good *goodsDomain) createEnvelopeNo() {
	good.EnvelopeNo = ksuid.New().Next().String()
}

func (good *goodsDomain) Create(dto services.RedEnvelopeGoodsDTO) {
	good.RedEnvelopeGoods.FromDTO(&dto)
	good.RemainQuantity = dto.Quantity
	good.Username.Valid = true
	good.Blessing.Valid = true

	// 普通红包
	if good.EnvelopeType == services.GeneralEnvelopeType {
		good.Amount = dto.AmountOne.Mul(decimal.NewFromFloat(float64(dto.Quantity)))
	}
	// 幸运红包
	if good.EnvelopeType == services.LuckyEnvelopeType {
		good.Amount = decimal.NewFromFloat(0)
	}
	good.RemainAmount = good.Amount
	good.ExpiredAt = time.Now().Add(time.Hour * 24)
	good.Status = services.OrderCreate
	good.createEnvelopeNo()
}

// Save 保存
func (good *goodsDomain) Save(ctx context.Context) (id int64, err error) {
	err = base.ExecuteContext(ctx, func(runner *dbx.TxRunner) error {
		dao := RedEnvelopeGoodsDao{runner}
		id, err = dao.Insert(good.RedEnvelopeGoods)
		return err
	})
	return
}

// CreateAndSave 创建并保存
func (good *goodsDomain) CreateAndSave(ctx context.Context, dto services.RedEnvelopeGoodsDTO) (id int64, err error) {
	good.Create(dto)
	return good.Save(ctx)
}
