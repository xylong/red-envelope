package envelopes

import (
	"context"
	"github.com/tietang/dbx"
	"path"
	"red-envelope/core/accounts"
	"red-envelope/infra/base"
	"red-envelope/services"
)

// SendOut 发红包
func (d *goodsDomain) SendOut(dto services.RedEnvelopeGoodsDTO) (activity *services.RedEnvelopeActivity, err error) {
	// 创建红包
	d.Create(dto)
	// 创建活动
	activity = new(services.RedEnvelopeActivity)
	domain := base.GetEnvelopeActivityDomain()
	link := base.GetEnvelopeActivityLink()
	activity.Link = path.Join(domain, link, d.EnvelopeNo)
	accountDomain := accounts.NewAccountDomain()

	err = base.Tx(func(runner *dbx.TxRunner) error {
		ctx := base.WithValueContext(context.Background(), runner)
		// 保存红包
		id, err := d.Save(ctx)
		if id < 0 || err != nil {
			return err
		}
		//1. 需要红包中间商的红包资金账户，定义在配置文件中，事先初始化到资金账户表中
		//2. 从红包发送人的资金账户中扣减红包金额 ，把红包金额从红包发送人的资金账户里扣除
		body := services.TradeParticipator{
			AccountNo: dto.AccountNo,
			UserId:    dto.UserId,
			Username:  dto.Username,
		}
		systemAccount := base.GetSystemAccount()
		target := services.TradeParticipator{
			AccountNo: systemAccount.AccountNo,
			UserId:    systemAccount.UserId,
			Username:  systemAccount.Username,
		}
		transfer := services.AccountTransferDTO{
			TradeNo:     d.EnvelopeNo,
			TradeBody:   body,
			TradeTarget: target,
			Amount:      d.Amount,
			ChangeType:  services.EnvelopeOutgoing,
			ChangeFlag:  services.FlagTransferOut,
			Decs:        "红包金额支付",
		}
		status, err := accountDomain.TransferWithContextTx(ctx, transfer)
		if status != services.TransferedStatusSuccess {
			return nil
		}
		//3. 将扣减的红包总金额转入红包中间商的红包资金账户
		transfer = services.AccountTransferDTO{
			TradeNo:     d.EnvelopeNo,
			TradeBody:   target,
			TradeTarget: body,
			Amount:      d.Amount,
			ChangeType:  services.EnvelopeIncoming,
			ChangeFlag:  services.FlagTransferIn,
			Decs:        "红包金额转入",
		}
		status, err = accountDomain.TransferWithContextTx(ctx, transfer)
		if status != services.TransferedStatusSuccess {
			return nil
		}

		return err
	})

	if err != nil {
		return nil, err
	}

	activity.RedEnvelopeGoodsDTO = *d.RedEnvelopeGoods.ToDTO()

	return
}
