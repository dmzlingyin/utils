package pay

import (
	"github.com/dmzlingyin/utils/config"
	"github.com/smartwalle/alipay/v3"
)

type AliPayReq struct {
	OutTradeNo string `json:"out_trade_no"` // 业务侧订单号
	Amount     string `json:"amount"`       // 订单金额(元)
	Subject    string `json:"subject"`      // 订单标题
	NotifyURL  string `json:"notify_url"`   // 支付宝异步通知地址
}

type Alipay struct {
	client *alipay.Client
}

func NewAlipay() (*Alipay, error) {
	appID := config.GetString("pay.alipay.app_id")
	privateKey := config.GetString("pay.alipay.private_key")
	isProduction := config.GetBool("pay.alipay.is_production")
	client, err := alipay.New(appID, privateKey, isProduction)
	if err != nil {
		return nil, err
	}
	return &Alipay{
		client: client,
	}, nil
}

func (p *Alipay) Pay(req *AliPayReq) (string, error) {
	return p.client.TradeAppPay(alipay.TradeAppPay{Trade: alipay.Trade{
		NotifyURL:   req.NotifyURL,
		Subject:     req.Subject,
		OutTradeNo:  req.OutTradeNo,
		TotalAmount: req.Amount,
	}})
}
