package pay

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/dmzlingyin/utils/cast"
	"github.com/dmzlingyin/utils/lazy"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/auth/verifiers"
	"github.com/wechatpay-apiv3/wechatpay-go/core/downloader"
	"github.com/wechatpay-apiv3/wechatpay-go/core/notify"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/native"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
	"net/http"
	"time"
)

type Transaction = payments.Transaction

type Config struct {
	AppID           string `option:"app_id"`
	MchID           string `option:"mch_id"`
	MchCertSerialNo string `option:"mch_cert_serial_no"`
	MchAPIv3Key     string `option:"mch_api_v3_key"`
	PrivateKey      string `option:"private_key"`
	NotifyURL       string `option:"notify_url"`
	PayType         string `option:"pay_type"` // 支付类型，比如native（网站支付），或者jsapi（公众号内支付）
}

type WechatPayResponse struct {
	TimeStamp string `json:"timeStamp"`
	NonceStr  string `json:"nonceStr"`
	Package   string `json:"package"`
	SignType  string `json:"signType"`
	PaySign   string `json:"paySign"`
	IsPay     bool   `json:"isPay"`
}

type WechatPay struct {
	cfg     *Config
	jss     *jsapi.JsapiApiService
	nas     *native.NativeApiService
	pk      *rsa.PrivateKey
	nh      *lazy.Value[*notify.Handler]
	options map[string]string
}

func newWechatPay(options map[string]string) (*WechatPay, error) {
	cfg := &Config{
		AppID:           options[OptionAppId],
		MchID:           options[OptionMchId],
		MchCertSerialNo: options[OptionMchCertSerialNo],
		MchAPIv3Key:     options[OptionMchAPIv3Key],
		PrivateKey:      options[OptionPrivateKey],
		NotifyURL:       options[OptionNotifyURL],
	}

	privateKey, err := base64.StdEncoding.DecodeString(cfg.PrivateKey)
	// 加载私钥
	key, err := utils.LoadPrivateKey(string(privateKey))
	if err != nil {
		return nil, err
	}

	// 创建客户端
	opts := []core.ClientOption{
		option.WithWechatPayAutoAuthCipher(cfg.MchID, cfg.MchCertSerialNo, key, cfg.MchAPIv3Key),
	}

	client, err := core.NewClient(context.Background(), opts...)
	if err != nil {
		return nil, err
	}

	s := &WechatPay{
		cfg:     cfg,
		jss:     &jsapi.JsapiApiService{Client: client},
		nas:     &native.NativeApiService{Client: client},
		pk:      key,
		options: options,
	}
	s.nh = &lazy.Value[*notify.Handler]{New: s.newNotifyHandler}
	return s, err
}

func (p *WechatPay) GetChannel() string {
	return ChannelWechat
}

func (p *WechatPay) HandleNotify(ctx context.Context, req *http.Request, handler func(t *Transaction) (args *UpdateStatusArgs)) (args *UpdateStatusArgs, err error) {
	// 获取通知处理器
	nh, err := p.nh.Get()
	if err != nil {
		return nil, fmt.Errorf("failed to create notify handler: %s", err)
	}

	// 解析请求
	rc := new(Transaction)
	if _, err := nh.ParseNotifyRequest(ctx, req, rc); err != nil {
		return nil, fmt.Errorf("failed to parse request: %s", err)
	}

	return handler(rc), err
}

func (p *WechatPay) Verify(ctx context.Context, args *VerifyArgs) (*VerifyRes, error) {
	if args.Money <= 0 {
		return nil, nil
	}

	res, err := p.Query(ctx, args.PayID)
	if err != nil {
		return nil, err
	}

	//s.logger.Infof("wechat verify: money: %d coupon: %d api res: %d", money, coupon, *res.Amount.Total)
	if res.Status != "SUCCESS" || args.Money != res.Money {
		return nil, errors.New("wechat pay verified failed")
	}
	return &VerifyRes{}, nil
}

func (p *WechatPay) Query(ctx context.Context, payId string) (*QueryResult, error) {
	res := new(QueryResult)
	// core是回调通知处理库
	req := jsapi.QueryOrderByOutTradeNoRequest{
		OutTradeNo: core.String(payId),
		Mchid:      core.String(p.cfg.MchID),
	}
	orderRes, _, err := p.jss.QueryOrderByOutTradeNo(ctx, req)
	if err != nil {
		return nil, err
	}
	res.Status = *orderRes.TradeState
	res.Money = int32(*orderRes.Amount.Total)
	if orderRes.TransactionId != nil {
		res.OrderId = *orderRes.TransactionId
	}
	return res, err
}

func (p *WechatPay) Create(ctx context.Context, req *CreateArgs) (*CreateResult, error) {
	switch req.PayType {
	case "native":
		return p.createNative(ctx, req)
	case "jsapi":
		return p.createJsapi(ctx, req)
	default:
		return nil, fmt.Errorf("invalid paytype, payType:%s", req.PayType)
	}
}

func (p *WechatPay) createNative(ctx context.Context, req *CreateArgs) (*CreateResult, error) {
	prepayRequest := native.PrepayRequest{
		Appid:       core.String(p.cfg.AppID),
		Mchid:       core.String(p.cfg.MchID),
		Description: core.String(req.Description),
		OutTradeNo:  core.String(req.OrderID),
		NotifyUrl:   core.String(p.cfg.NotifyURL),
		Amount: &native.Amount{
			Currency: core.String("CNY"),
			Total:    core.Int64(int64(req.Money)),
		},
	}
	result, _, err := p.nas.Prepay(ctx, prepayRequest)
	if err != nil {
		return nil, err
	}
	//s.logger.Debugf("code url: %s", *result.CodeUrl)
	return &CreateResult{
		CodeURL: *result.CodeUrl,
	}, nil
}

func (p *WechatPay) createJsapi(ctx context.Context, req *CreateArgs) (*CreateResult, error) {
	prepayRequest := jsapi.PrepayRequest{
		Appid:       core.String(p.cfg.AppID),
		Mchid:       core.String(p.cfg.MchID),
		Description: core.String(req.Description),
		OutTradeNo:  core.String(req.OrderID),
		NotifyUrl:   core.String(p.cfg.NotifyURL),
		Amount: &jsapi.Amount{
			Currency: core.String("CNY"),
			Total:    core.Int64(int64(req.Money)),
		},
		Payer: &jsapi.Payer{
			Openid: core.String(req.CustomerID),
		},
	}
	rep, _, err := p.jss.Prepay(ctx, prepayRequest)
	if err != nil {
		return nil, err
	}
	//s.logger.Debugf("code url: %s", *rep.PrepayId)
	return &CreateResult{
		OrderID: *rep.PrepayId,
	}, nil
}

func (p *WechatPay) newNotifyHandler() (h *notify.Handler, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	// 注册下载器
	err = downloader.MgrInstance().RegisterDownloaderWithPrivateKey(
		ctx, p.pk, p.cfg.MchCertSerialNo, p.cfg.MchID, p.cfg.MchAPIv3Key,
	)
	if err != nil {
		return nil, errors.New("failed to register downloader")
	}

	cm := downloader.MgrInstance().GetCertificateVisitor(p.cfg.MchID)
	return notify.NewRSANotifyHandler(p.cfg.MchAPIv3Key, verifiers.NewSHA256WithRSAVerifier(cm))
}

func (p *WechatPay) GenerateResponse(paymentId string) (*WechatPayResponse, error) {
	res := &WechatPayResponse{}
	// 时间戳(秒)
	timeUnix := time.Now().Unix()
	res.TimeStamp = cast.ToString(timeUnix)
	// 随机字符串
	nonceStr, err := utils.GenerateNonce()
	if err != nil {
		return nil, nil
	}
	res.NonceStr = nonceStr

	// 订单详情扩展字符串
	res.Package = "prepay_id=" + paymentId
	// 签名方式
	res.SignType = "RSA"

	str := p.cfg.AppID + "\n" + res.TimeStamp + "\n" + res.NonceStr + "\n" + res.Package + "\n"
	sign, err := utils.SignSHA256WithRSA(str, p.pk)
	if err != nil {
		return nil, err
	}
	res.PaySign = sign

	return res, nil
}

func (p *WechatPay) CreateSub(ctx context.Context, args *CreateSubArgs) (res *CreateSubResult, err error) {
	return nil, errors.New("not yet implemented")
}

func (p *WechatPay) QuerySub(ctx context.Context, args *QuerySubArgs) (*SubDetail, error) {
	return nil, errors.New("not yet implemented")
}

func (p *WechatPay) Capture(ctx context.Context, orderID string, amount int32) (string, error) {
	return "", errors.New("not yet implemented")
}

func (p *WechatPay) CreatePortal(ctx context.Context, args *CreatePortalArgs) (*CreatePortalResult, error) {
	return nil, errors.New("not yet implemented")
}
