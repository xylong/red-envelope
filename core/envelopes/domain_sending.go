package envelopes

import (
	"fmt"
	"path"
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
	// 保存红包
	// 支付
}
