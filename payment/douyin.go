package pay

import (
	"bytes"
	"context"
	"crypto/md5"
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"
)

type DouyinConfig struct {
	AppID     string
	MchID     string
	NotifyURL string
	Secret    string
	Salt      string
	Token     string
}

type PaymentInfo struct {
	TotalFee    int    `json:"total_fee"`
	OrderStatus string `json:"order_status"`
	PayTime     string `json:"pay_time"`
	Way         int    `json:"way"`
	ChannelNo   string `json:"channel_no"`
	SellerUid   string `json:"seller_uid"`
	ItemId      string `json:"item_id"`
	CpsInfo     string `json:"cps_info"`
}

// CreateReq 具体用法，详见: https://developer.open-douyin.com/docs/resource/zh-CN/mini-app/develop/server/ecpay/pay-list/pay
type CreateReq struct {
	AppID       string `json:"app_id"`       // appid
	OutOrderNo  string `json:"out_order_no"` // 商户内部订单号
	TotalAmount int32  `json:"total_amount"` // 总费用
	Subject     string `json:"subject"`      // 商品描述
	Body        string `json:"body"`         // 商品详情
	ValidTime   int32  `json:"valid_time"`   // 订单过期时间(最大两天)
	Sign        string `json:"sign"`         // 签名
	NotifyUrl   string `json:"notify_url"`   // 回调地址
}

type CreateResp struct {
	ErrNo   int    `json:"err_no"`
	ErrTips string `json:"err_tips"`
	Data    struct {
		OrderID    string `json:"order_id"`
		OrderToken string `json:"order_token"`
	} `json:"data"`
}

type QueryResp struct {
	ErrNo       int         `json:"err_no"`
	ErrTips     string      `json:"err_tips"`
	OutOrderNo  string      `json:"out_order_no"`
	OrderId     string      `json:"order_id"`
	PaymentInfo PaymentInfo `json:"payment_info"`
}

// NotifyResp 具体用法，详见：https://developer.open-douyin.com/docs/resource/zh-CN/mini-app/develop/server/ecpay/pay-list/callback
type NotifyResp struct {
	Timestamp    string `json:"timestamp"`     // 时间戳
	Nonce        string `json:"nonce"`         // 随机字符串
	Msg          string `json:"msg"`           // 订单信息的json字符串，对应下面的NotifyMsg结构体
	Type         string `json:"type"`          // 回调类型标记，支付成功固定值为"payment"
	MsgSignature string `json:"msg_signature"` // 签名
}

type NotifyMsg struct {
	AppID          string `json:"appid"`        // appid
	CpOrderNo      string `json:"cp_orderno"`   // 商户内部订单号
	CpExtra        string `json:"cp_extra"`     // 预下单传入字段
	Way            string `json:"way"`          // 支付渠道标识: 1-微信支付 2-支付宝支付 10-抖音支付
	ChannelNo      string `json:"channel_no"`   // 支付渠道侧单号
	PaymentOrderNo string `json:""`             // 微信或支付宝单号
	TotalAmount    int32  `json:"total_amount"` // 支付金额(分)
	Status         string `json:"status"`       // 固定SUCCESS
	ItemID         string `json:"item_id"`
	SellerUid      string `json:"seller_uid"` // 卖家商户号
	PaidAt         int32  `json:"paid_at"`    // 支付时间
	OrderID        string `json:"order_id"`   // 抖音侧单号
}

type DouyinPay struct {
	cfg     *DouyinConfig
	options map[string]string
}

func newDouyinPay(options map[string]string) *DouyinPay {
	cfg := &DouyinConfig{
		AppID:     options[OptionAppId],
		MchID:     options[OptionMchId],
		Secret:    options[OptionSecret],
		Salt:      options[OptionSalt],
		NotifyURL: options[OptionNotifyURL],
		Token:     options[OptionToken],
	}
	return &DouyinPay{
		cfg:     cfg,
		options: options,
	}
}

func (p *DouyinPay) GetChannel() string {
	return ChannelDouyin
}

func (p *DouyinPay) Create(ctx context.Context, args *CreateArgs) (*CreateResult, error) {
	paramsMap := map[string]any{
		"oon":     args.OrderID,
		"amount":  args.Money,
		"subject": args.Description,
		"body":    args.Description,
		"vt":      900,
		"nt":      p.cfg.NotifyURL,
	}
	sign := p.RequestSign(paramsMap)
	req := CreateReq{
		AppID:       p.cfg.AppID,
		OutOrderNo:  args.OrderID,
		TotalAmount: args.Money,
		Subject:     args.Description,
		Body:        args.Description,
		ValidTime:   900,
		Sign:        sign,
		NotifyUrl:   p.cfg.NotifyURL,
	}

	url := "https://developer.toutiao.com/api/apps/ecpay/v1/create_order"
	breq, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(breq))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var cresp CreateResp
	if err = json.NewDecoder(resp.Body).Decode(&cresp); err != nil {
		return nil, err
	}

	if cresp.ErrNo != 0 {
		return nil, errors.New(cresp.ErrTips)
	}

	return &CreateResult{
		OrderID:    cresp.Data.OrderID,
		OrderToken: cresp.Data.OrderToken,
	}, nil
}

func (p *DouyinPay) Verify(ctx context.Context, args *VerifyArgs) (*VerifyRes, error) {
	if args.Money <= 0 {
		return nil, nil
	}
	res, err := p.Query(ctx, args.PayID)
	if err != nil {
		return nil, err
	}
	if res.Money != args.Money || res.Status != "SUCCESS" {
		return nil, errors.New("douyin verify failed")
	}
	return &VerifyRes{}, nil
}

func (p *DouyinPay) Query(ctx context.Context, payID string) (res *QueryResult, err error) {
	m := map[string]any{"oon": payID}
	sign := p.RequestSign(m)
	var req = struct {
		AppID      string `json:"app_id"`
		OutOrderNo string `json:"out_order_no"`
		Sign       string `json:"sign"`
	}{
		p.cfg.AppID,
		payID,
		sign,
	}

	url := "https://developer.toutiao.com/api/apps/ecpay/v1/query_order"
	breq, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(breq))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var queryResp QueryResp
	if err = json.NewDecoder(resp.Body).Decode(&queryResp); err != nil {
		return nil, err
	}
	if queryResp.ErrNo != 0 {
		return nil, errors.New(queryResp.ErrTips)
	}
	return &QueryResult{
		Money:   int32(queryResp.PaymentInfo.TotalFee),
		Status:  queryResp.PaymentInfo.OrderStatus,
		OrderId: queryResp.OrderId,
	}, nil
}

// RequestSign 担保支付请求签名算法.
// 参数："paramsMap" 所有的请求参数
func (p *DouyinPay) RequestSign(paramsMap map[string]interface{}) string {
	var paramsArr []string
	for _, v := range paramsMap {
		value := strings.TrimSpace(fmt.Sprintf("%v", v))
		if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") && len(value) > 1 {
			value = value[1 : len(value)-1]
		}
		value = strings.TrimSpace(value)
		if value == "" || value == "null" {
			continue
		}
		paramsArr = append(paramsArr, value)
	}

	paramsArr = append(paramsArr, p.cfg.Salt)
	sort.Strings(paramsArr)
	return fmt.Sprintf("%x", md5.Sum([]byte(strings.Join(paramsArr, "&"))))
}

// HandleNotify 负责处理抖音的回调
func (p *DouyinPay) HandleNotify(req *http.Request, handler func(msg *NotifyMsg) (args *UpdateStatusArgs)) (args *UpdateStatusArgs, err error) {
	var notifyResp NotifyResp
	if err = json.NewDecoder(req.Body).Decode(&notifyResp); err != nil {
		return nil, err
	}
	// 验签
	if !p.checkSign(&notifyResp) {
		return nil, errors.New("callback sign error")
	}

	var notifyMsg NotifyMsg
	err = json.Unmarshal([]byte(notifyResp.Msg), &notifyMsg)
	if err != nil {
		return nil, err
	}

	return handler(&notifyMsg), err
}

// checkSing 用于验证回调签名
func (p *DouyinPay) checkSign(notifyResp *NotifyResp) bool {
	s := make([]string, 0)
	s = append(s, notifyResp.Timestamp)
	s = append(s, notifyResp.Nonce)
	s = append(s, notifyResp.Msg)
	s = append(s, p.cfg.Token)
	sort.Strings(s)

	h := sha1.New()
	h.Write([]byte(strings.Join(s, "")))
	return notifyResp.MsgSignature == fmt.Sprintf("%x", h.Sum(nil))
}

func (p *DouyinPay) CreateSub(ctx context.Context, args *CreateSubArgs) (res *CreateSubResult, err error) {
	return nil, errors.New("not yet implemented")
}

func (p *DouyinPay) QuerySub(ctx context.Context, args *QuerySubArgs) (*SubDetail, error) {
	return nil, errors.New("not yet implemented")
}

func (p *DouyinPay) Capture(ctx context.Context, orderID string, amount int32) (string, error) {
	return "", errors.New("not yet implemented")
}

func (p *DouyinPay) CreatePortal(ctx context.Context, args *CreatePortalArgs) (*CreatePortalResult, error) {
	return nil, errors.New("not yet implemented")
}
