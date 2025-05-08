package payment

import (
	"time"
)

// 各个平台的Option字段
const (
	// google
	OptionPkgName = "package_name"
	// paypal
	OptionClientId = "client_id"
	OptionSecretId = "secret_id"
	OptionSandbox  = "sandbox"
	// WeChat
	OptionAppId     = "app_id"
	OptionMchId     = "mch_id"
	OptionNotifyURL = "notify_url"
	// douyin
	OptionSecret = "app_secret"
	OptionSalt   = "salt"
	OptionToken  = "token"
)

type UpdateStatusArgs struct {
	BizID      string `json:"bizId"` // 业务系统ID
	CustomerID string `json:"customerId"`
	PaymentID  string `json:"paymentId"` // 第三方平台ID
	Status     int32  `json:"status"`
	Money      int32  `json:"money,omitempty"` // 支付金额
	OutOrderId string `json:"outOrderId"`      // 商户订单号
}

type VerifyArgs struct {
	Kind      string // payment/subscription (普通支付/订阅支付)
	PayID     string
	Receipt   string // 对应 iOS 的 ReceiptData，对应 Android 的 PurchaseToken
	ProductID string
	Money     int32 // 金额，单位分
}

type VerifyRes struct {
	Sandbox    bool      // 是否为沙盒环境
	OrderID    string    // 订单ID
	ProductID  string    // 订阅产品ID
	StartTime  time.Time // 订阅开始时间
	ExpiryTime time.Time // 订阅到期时间
}

type CreateArgs struct {
	CustomerID  string // 对应wechat的openid
	Money       int32  // 金额，单位分
	Description string // 描述信息
	OrderID     string // 订单ID
	PayType     string // 支付方式 jsapi/native
	ReturnURL   string // 支付完成的回调地址
	CancelURL   string // stripe取消支付的重定向URL
	PriceID     string // stripe管理后台配置的价格ID
}

type CreateResult struct {
	OrderID    string          // 订单id
	OrderToken string          // 订单token
	CodeURL    string          // 二维码，可以根据此生成二维码，让用户扫码支付
	OutOrderID string          // 商户内部订单
	WechatRes  *WechatJsapiRes // 微信Jsapi额外返回信息
}

type CreateSubArgs struct {
	PlanID              string    // 计划ID(兼容stripe的priceID)
	ReturnURL           string    // 用户订阅后的重定向URL
	CancelURL           string    // stripe取消支付的重定向URL
	StartTime           time.Time // 订阅开始时间
	BizID               string    // 业务侧ID
	CustomerID          string    // stripe侧用户ID
	BizUserID           string    // 用于创建stripe customerID, bizType-userID
	AllowPromotionCodes bool      // 是否开启 stripe 促销码
}

type CreateSubResult struct {
	SubID      string // 订阅ID
	PayURL     string // 支付URL
	CustomerID string // stripe侧用户ID
	SessionID  string // stripe侧首次订阅支付ID
}

type QuerySubArgs struct {
	SubID     string // 订阅ID
	SessionID string // stripe的支付会话ID
}

type SubDetail struct {
	PlanID          string    // 计划ID
	SubID           string    // 订阅ID
	CyclesCompleted int32     // 已完成的周期数
	Status          string    // 订阅状态
	LastPaymentTime time.Time // 上次支付时间
	NextBillingTime time.Time // 下次付款时间
}

type WechatJsapiRes struct {
	IsPay     bool
	TimeStamp string
	NonceStr  string
	Package   string
	SignType  string
	PaySign   string
}

type QueryResult struct {
	Money   int32
	Status  string
	OrderId string
}

type CreatePortalArgs struct {
	CustomerID string
	ReturnURL  string
}

type CreatePortalResult struct {
	URL string // stripe用户门户链接
}
