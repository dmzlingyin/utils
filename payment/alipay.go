package pay

import (
	"context"
	"github.com/dmzlingyin/utils/config"
	"github.com/smartwalle/alipay/v3"
	"net/url"
)

type AliPayReq struct {
	OutTradeNo string `json:"out_trade_no"` // 业务侧订单号
	Amount     string `json:"amount"`       // 订单金额(元)
	Subject    string `json:"subject"`      // 订单标题
	NotifyURL  string `json:"notify_url"`   // 支付宝异步通知地址
	ReturnURL  string `json:"return_url"`
}

type (
	QueryRes           alipay.TradeQueryRsp
	AlipayNotification alipay.Notification
)

type Alipay struct {
	client *alipay.Client
}

func NewAlipay() (*Alipay, error) {
	appID := config.GetString("pay.alipay.app_id")
	privateKey := config.GetString("pay.alipay.private_key")
	publicKey := config.GetString("pay.alipay.public_key")
	isProduction := config.GetBool("pay.alipay.is_production")
	client, err := alipay.New(appID, privateKey, isProduction)
	if err != nil {
		return nil, err
	}
	if err = client.LoadAliPayPublicKey(publicKey); err != nil {
		return nil, err
	}
	return &Alipay{
		client: client,
	}, nil
}

func (p *Alipay) Pay(req *AliPayReq) (string, error) {
	return p.client.TradeAppPay(alipay.TradeAppPay{Trade: alipay.Trade{
		NotifyURL:   req.NotifyURL,
		ReturnURL:   req.ReturnURL,
		Subject:     req.Subject,
		OutTradeNo:  req.OutTradeNo,
		TotalAmount: req.Amount,
	}})
}

func (p *Alipay) Query(ctx context.Context, outTradeNo string) (*QueryRes, error) {
	res, err := p.client.TradeQuery(ctx, alipay.TradeQuery{OutTradeNo: outTradeNo})
	if err != nil {
		return nil, err
	}
	return (*QueryRes)(res), nil
}

func (p *Alipay) ParseNotify(value url.Values) (*AlipayNotification, error) {
	res, err := p.client.DecodeNotification(value)
	if err != nil {
		return nil, err
	}
	return (*AlipayNotification)(res), nil
}
