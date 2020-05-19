package services

// OrderType 订单类型
type OrderType int

const (
	// OrderTypeSending 发布单
	OrderTypeSending OrderType = 1
	// OrderTypeRefund 退款单
	OrderTypeRefund OrderType = 2
)

// PayStatus 支付状态
type PayStatus int

const (
	// PayNothing 未支付
	PayNothing PayStatus = 1
	// Paying 支付中
	Paying PayStatus = 2
	// Payed 已支付
	Payed PayStatus = 3
	//PayFailure 支付失败
	PayFailure PayStatus = 4
	// RefundNothing 未退款
	RefundNothing PayStatus = 61
	// Refunding 退款中
	Refunding PayStatus = 62
	// Refunded 已退款
	Refunded PayStatus = 63
	// RefundFailure 退款失败
	RefundFailure PayStatus = 64
)

// OrderStatus订单状态
type OrderStatus int

const (
	// OrderCreate 创建
	OrderCreate OrderStatus = 1
	// OrderSending 发布
	OrderSending OrderStatus = 2
	// OrderExpired 过期
	OrderExpired OrderStatus = 3
	// OrderDisabled 失效
	OrderDisabled OrderStatus = 4
	// OrderExpiredRefundSuccessful 过期退款成功
	OrderExpiredRefundSuccessful OrderStatus = 5
	// OrderExpiredRefundFalured 过期退款失败
	OrderExpiredRefundFalured OrderStatus = 6
)

// EnvelopeType 红包类型
type EnvelopeType int

const (
	// GeneralEnvelopeType 普通红包
	GeneralEnvelopeType = 1
	// LuckyEnvelopeType 幸运红包
	LuckyEnvelopeType = 2
)

var EnvelopeTypes = map[EnvelopeType]string{
	GeneralEnvelopeType: "普通红包",
	LuckyEnvelopeType:   "碰运气红包",
}
