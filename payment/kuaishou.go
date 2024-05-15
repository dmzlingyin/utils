package pay

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

// DURATION access token的有效时间，快手的有效期最大为48小时
// 业务代码设置为24小时，预留足够的时间
const DURATION = 24

// GoodsType 商品类别代码, 不同类别手续费不同
// 当前快手处于试运行阶段，所有手续费为2%
// 虚拟/服务 -> 虚拟卡/会员/游戏 -> 娱乐会员
const GoodsType = 3314

// KSCreateReq 详情:https://mp.kuaishou.com/docs/develop/server/epay/interfaceDefinitionWithoutChannel.html
type KSCreateReq struct {
	OutOrderNo  string `json:"out_order_no"` // 商户内部订单号,长度[6,32]
	OpenID      string `json:"open_id"`
	TotalAmount int32  `json:"total_amount"` // 支付金额，单位分
	Subject     string `json:"subject"`      // 商品描述,长度[1,128]
	Detail      string `json:"detail"`       // 商品详情,长度[1,1024]
	Type        int32  `json:"type"`         // 商品类型
	ExpireTime  int32  `json:"expire_time"`  // 订单过期时间,范围[300,172800]s
	Sign        string `json:"sign"`         // 核心字段签名
	NotifyURL   string `json:"notify_url"`   // 回调地址
}

type KuaishouPay struct {
	appID     string
	appSecret string
	notifyURL string
	at        string // access token
	preTime   time.Time
	options   map[string]string
}

func newKuaishouPay(options map[string]string) (*KuaishouPay, error) {
	ks := &KuaishouPay{
		appID:     options[OptionAppId],
		appSecret: options[OptionSecret],
		notifyURL: options[OptionNotifyURL],
	}
	if err := ks.refreshAT(); err != nil {
		return nil, err
	}
	return ks, nil
}

func (p *KuaishouPay) GetChannel() string {
	return ChannelKuaishou
}

func (p *KuaishouPay) Verify(ctx context.Context, args *VerifyArgs) (*VerifyRes, error) {
	if err := p.refreshAT(); err != nil {
		return nil, err
	}
	res, err := p.Query(ctx, args.PayID)
	if err != nil {
		return nil, err
	}

	if res.Money != args.Money || res.Status != "SUCCESS" {
		return nil, errors.New("kuaishou verify failed: payAmount or payStatus check failed")
	}
	return &VerifyRes{}, nil
}

func (p *KuaishouPay) Query(ctx context.Context, payID string) (res *QueryResult, err error) {
	base := "https://open.kuaishou.com/openapi/mp/developer/epay/query_order"
	queryUrl := fmt.Sprintf("%s?app_id=%s&access_token=%s", base, p.appID, p.at)
	sign := p.SignVerify(payID)
	var req = struct {
		OutOrderNo string `json:"out_order_no"`
		Sign       string `json:"sign"`
	}{
		payID,
		sign,
	}
	breq, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(queryUrl, "application/json", bytes.NewBuffer(breq))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var orderData = struct {
		Result      int32  `json:"result"`
		ErrorMsg    string `json:"error_msg"`
		PaymentInfo struct {
			TotalAmount     int32  `json:"total_amount"`     // 支付金额
			PayStatus       string `json:"pay_status"`       // 支付状态
			PayTime         int64  `json:"pay_time"`         // 支付时间，毫秒时间戳
			PayChannel      string `json:"pay_channel"`      // 支付渠道
			OutOrderNo      string `json:"out_order_no"`     // 商户订单号
			KsOrderNo       string `json:"ks_order_no"`      // 快手侧订单号
			ExtraInfo       string `json:"extra_info"`       // 订单信息来源(直播场景/短视频场景)
			EnablePromotion bool   `json:"enable_promotion"` // 是否参与分销
			PromotionAmount int    `json:"promotion_amount"` // 预计分销金额
			OpenID          string `json:"open_id"`
		} `json:"payment_info"`
	}{}
	if err = json.NewDecoder(resp.Body).Decode(&orderData); err != nil {
		return nil, err
	}
	if orderData.Result != 1 {
		return nil, errors.New(orderData.ErrorMsg)
	}
	return &QueryResult{
		Money:   orderData.PaymentInfo.TotalAmount,
		Status:  orderData.PaymentInfo.PayStatus,
		OrderId: orderData.PaymentInfo.KsOrderNo,
	}, err
}

func (p *KuaishouPay) refreshAT() error {
	if p.at == "" || p.IsExpired() {
		at, err := GetAccessToken(p.appID, p.appSecret)
		if err != nil {
			return err
		}
		p.at = at
		p.preTime = time.Now()
	}
	return nil
}

func (p *KuaishouPay) IsExpired() bool {
	return time.Now().Sub(p.preTime).Hours() > DURATION
}

func GetAccessToken(appID, appSecret string) (string, error) {
	addr := "https://open.kuaishou.com/oauth2/access_token"
	pd := url.Values{}
	pd.Add("app_id", appID)
	pd.Add("app_secret", appSecret)
	pd.Add("grant_type", "client_credentials") // 固定值

	resp, err := http.Post(addr, "application/x-www-form-urlencoded", strings.NewReader(pd.Encode()))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var res = struct {
		Result      int32  `json:"result"`
		AccessToken string `json:"access_token"`
		ExpiresIn   int32  `json:"expires_in"`
		TokenType   string `json:"bearer"`
	}{}
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}
	if res.Result != 1 {
		return "", errors.New("failed to get access token")
	}
	return res.AccessToken, nil
}

func (p *KuaishouPay) SignVerify(oon string) string {
	str := fmt.Sprintf("app_id=%s&out_order_no=%s", p.appID, oon) + p.appSecret
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}

func (p *KuaishouPay) Create(ctx context.Context, args *CreateArgs) (orderRes *CreateResult, err error) {
	if err := p.refreshAT(); err != nil {
		return nil, err
	}

	base := "https://open.kuaishou.com/openapi/mp/developer/epay/create_order_with_channel"
	url := fmt.Sprintf("%s?app_id=%s&access_token=%s", base, p.appID, p.at)
	sign := p.Sign(args, GoodsType, 900, p.notifyURL)
	req := KSCreateReq{
		OutOrderNo:  args.OrderID,
		OpenID:      args.CustomerID,
		TotalAmount: args.Money,
		Subject:     args.Description,
		Detail:      args.Description,
		Type:        GoodsType,
		ExpireTime:  900,
		Sign:        sign,
		NotifyURL:   p.notifyURL,
	}
	breq, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(breq))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var res = struct {
		Result    int    `json:"result"`
		ErrorMsg  string `json:"error_msg"`
		OrderInfo struct {
			OrderNo        string `json:"order_no"`
			OrderInfoToken string `json:"order_info_token"`
		} `json:"order_info"`
	}{}
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}
	if res.Result != 1 {
		return nil, errors.New(res.ErrorMsg)
	}
	return &CreateResult{
		OrderID:    res.OrderInfo.OrderNo,
		OrderToken: res.OrderInfo.OrderInfoToken,
	}, nil
}

func (p *KuaishouPay) Sign(args *CreateArgs, goodsType, expireTime int32, notifyURL string) string {
	signParam := make(map[string]any)
	signParam["app_id"] = p.appID
	signParam["open_id"] = args.CustomerID
	signParam["out_order_no"] = args.OrderID
	signParam["total_amount"] = args.Money
	signParam["subject"] = args.Description
	signParam["detail"] = args.Description
	signParam["expire_time"] = expireTime
	signParam["notify_url"] = notifyURL
	signParam["type"] = goodsType

	params := []string{"app_id", "open_id", "out_order_no", "total_amount", "subject", "detail", "type", "expire_time", "notify_url"}
	sort.Strings(params)

	var str []string
	for _, param := range params {
		str = append(str, fmt.Sprintf("%s=%v", param, signParam[param]))
	}
	signStr := strings.Join(str, "&") + p.appSecret
	return fmt.Sprintf("%x", md5.Sum([]byte(signStr)))
}

// HandleNotify 负责处理快手的回调
func (p *KuaishouPay) HandleNotify(req *http.Request, handler func(orderId, status string, amount int) (args *UpdateStatusArgs)) (args *UpdateStatusArgs, message string, err error) {
	var r = struct {
		Data struct {
			Channel         string `json:"channel"`
			OutOrderNo      string `json:"out_order_no"`
			Attach          string `json:"attach"`
			Status          string `json:"status"`
			KsOrderNo       string `json:"ks_order_no"`
			OrderAmount     int    `json:"order_amount"`
			TradeNo         string `json:"trade_no"`
			ExtraInfo       string `json:"extra_info"`
			EnablePromotion bool   `json:"enable_promotion"`
			PromotionAmount int    `json:"promotion_amount"`
		} `json:"data"`
		BizType   string `json:"biz_type"`
		MessageID string `json:"message_id"`
		AppID     string `json:"app_id"`
		Timestamp int64  `json:"timestamp"`
	}{}

	var buf bytes.Buffer
	io.Copy(&buf, req.Body)
	if err = json.Unmarshal(buf.Bytes(), &r); err != nil {
		return nil, "", err
	}
	// 验签
	if !p.checkSign(string(buf.Bytes()), req.Header.Get("kwaisign")) {
		return nil, "", errors.New("callback sign error")
	}

	return handler(r.Data.OutOrderNo, r.Data.Status, r.Data.OrderAmount), r.MessageID, nil
}

func (p *KuaishouPay) checkSign(msg, sign string) bool {
	expected := fmt.Sprintf("%x", md5.Sum([]byte(msg+p.appSecret)))
	return expected == sign
}

func (p *KuaishouPay) CreateSub(ctx context.Context, args *CreateSubArgs) (res *CreateSubResult, err error) {
	return nil, errors.New("not yet implemented")
}

func (p *KuaishouPay) QuerySub(ctx context.Context, args *QuerySubArgs) (*SubDetail, error) {
	return nil, errors.New("not yet implemented")
}

func (p *KuaishouPay) Capture(ctx context.Context, orderID string, amount int32) (string, error) {
	return "", errors.New("not yet implemented")
}

func (p *KuaishouPay) CreatePortal(ctx context.Context, args *CreatePortalArgs) (*CreatePortalResult, error) {
	return nil, errors.New("not yet implemented")
}
