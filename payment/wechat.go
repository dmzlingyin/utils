package pay

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/dmzlingyin/utils/config"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/auth/verifiers"
	"github.com/wechatpay-apiv3/wechatpay-go/core/downloader"
	"github.com/wechatpay-apiv3/wechatpay-go/core/notify"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/app"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/h5"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/native"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
	"net/http"
	"time"
)

const (
	WechatPayTypeApp    = "app"
	WechatPayTypeH5     = "h5"
	WechatPayTypeNative = "native"
	WechatPayTypeJsapi  = "jsapi"
)

const (
	WechatPayTradeStateSuccess      = "SUCCESS"    // 支付成功
	WechatPayTradeStateRefund       = "REFUND"     // 转入退款
	WechatPayTradeStateNotPay       = "NOTPAY"     // 未支付
	WechatPayTradeStateClosed       = "CLOSED"     // 已关闭
	WechatPayTradeStateRevoked      = "REVOKED"    // 已撤销(付款码支付)
	WechatPayTradeStateUserPaying   = "USERPAYING" // 用户支付中(付款码支付)
	WechatPayTradeStateUserPayError = "PAYERROR"   // 支付失败
)

type WechatPrepayReq struct {
	OutTradeNo  string // 商户内部订单号
	Amount      int64  // 支付金额(分)
	Description string // 商品描述
	OpenID      string // 用户在普通商户AppID下的唯一标识
	PayType     string // 支付类型: app、h5、jsapi、native
}

type WechatPrepayResp struct {
	PrepayID string `json:"prepay_id,omitempty"` // 预支付ID(app)
	H5Url    string `json:"h5_url,omitempty"`    // 支付跳转链接(h5)
	CodeUrl  string `json:"code_url,omitempty"`  // 用于生成支付二维码，然后提供给用户扫码支付(native)
	PrepayWithRequestPaymentResp
}

type PrepayWithRequestPaymentResp struct {
	AppID     string `json:"app_id,omitempty"`     // 应用ID
	PartnerID string `json:"partner_id,omitempty"` // 直连商户号
	TimeStamp string `json:"timestamp,omitempty"`  // 时间戳
	NonceStr  string `json:"nonce_str,omitempty"`  // 随机字符串
	Package   string `json:"package,omitempty"`    // 暂填写固定值: WXPay
	SignType  string `json:"sign_type,omitempty"`  // 签名类型: RSA
	Sign      string `json:"pay_sign,omitempty"`   // 签名值
}

type Transaction payments.Transaction

type WechatPayConfig struct {
	AppID           string // 应用ID
	MchID           string // 直连商户号
	MchCertSerialNo string // 证书序列号
	MchAPIv3Key     string // api v3秘钥
	PrivateKeyPath  string // 私钥路径
	NotifyURL       string // 回调地址
}

type WechatNotifyResp struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type WechatPay struct {
	cfg *WechatPayConfig
	aas *app.AppApiService       // app支付
	has *h5.H5ApiService         // h5支付
	jss *jsapi.JsapiApiService   // jsapi支付
	nas *native.NativeApiService // native支付(扫码支付)
	pk  *rsa.PrivateKey
	nh  *notify.Handler
}

func NewWechatPay() (*WechatPay, error) {
	cfg := &WechatPayConfig{
		AppID:           config.GetString("pay.wechat.app_id"),
		MchID:           config.GetString("pay.wechat.mch_id"),
		MchCertSerialNo: config.GetString("pay.wechat.mch_cert_serial_no"),
		MchAPIv3Key:     config.GetString("pay.wechat.mch_api_v3_key"),
		PrivateKeyPath:  config.GetString("pay.wechat.private_key_path"),
		NotifyURL:       config.GetString("pay.wechat.notify_url"),
	}

	privateKey, err := base64.StdEncoding.DecodeString(cfg.PrivateKeyPath)
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
		cfg: cfg,
		aas: &app.AppApiService{Client: client},
		has: &h5.H5ApiService{Client: client},
		jss: &jsapi.JsapiApiService{Client: client},
		nas: &native.NativeApiService{Client: client},
		pk:  key,
	}
	s.nh, err = s.newNotifyHandler()
	return s, err
}

// PrePay 商户系统先调用该接口在微信支付服务后台生成预支付交易单，返回正确的预支付交易会话标识后再按Native、JSAPI、APP等不同场景生成交易串调起支付。
func (p *WechatPay) PrePay(ctx context.Context, req *WechatPrepayReq) (*WechatPrepayResp, error) {
	switch req.PayType {
	case WechatPayTypeApp:
		return p.prepayApp(ctx, req)
	case WechatPayTypeH5:
		return p.prepayH5(ctx, req)
	case WechatPayTypeNative:
		return p.prepayNative(ctx, req)
	case WechatPayTypeJsapi:
		return p.prepayJsapi(ctx, req)
	default:
		return nil, fmt.Errorf("invalid paytype, payType:%s", req.PayType)
	}
}

// prepayApp 商户系统先调用该接口在微信支付服务后台生成预支付交易单，返回正确的预支付交易会话标识后再按Native、JSAPI、APP等不同场景生成交易串调起支付
func (p *WechatPay) prepayApp(ctx context.Context, req *WechatPrepayReq) (*WechatPrepayResp, error) {
	prepayRequest := app.PrepayRequest{
		Appid:       core.String(p.cfg.AppID),
		Mchid:       core.String(p.cfg.MchID),
		Description: core.String(req.Description),
		OutTradeNo:  core.String(req.OutTradeNo),
		NotifyUrl:   core.String(p.cfg.NotifyURL),
		Amount: &app.Amount{
			Total: core.Int64(req.Amount),
		},
	}
	resp, result, err := p.aas.PrepayWithRequestPayment(ctx, prepayRequest)
	if err != nil {
		return nil, err
	}
	if result.Response.StatusCode != http.StatusOK {
		return nil, errors.New(result.Response.Status)
	}
	return &WechatPrepayResp{
		PrepayID: *resp.PrepayId,
		PrepayWithRequestPaymentResp: PrepayWithRequestPaymentResp{
			AppID:     p.cfg.AppID,
			PartnerID: *resp.PartnerId,
			TimeStamp: *resp.TimeStamp,
			NonceStr:  *resp.NonceStr,
			Package:   *resp.Package,
			SignType:  "RSA",
			Sign:      *resp.Sign,
		},
	}, nil
}

// prepayH5 拉起微信支付收银台的中间页面，可通过访问该URL来拉起微信客户端，完成支付，h5_url的有效期为5分钟
func (p *WechatPay) prepayH5(ctx context.Context, req *WechatPrepayReq) (*WechatPrepayResp, error) {
	prepayRequest := h5.PrepayRequest{
		Appid:       core.String(p.cfg.AppID),
		Mchid:       core.String(p.cfg.MchID),
		Description: core.String(req.Description),
		OutTradeNo:  core.String(req.OutTradeNo),
		NotifyUrl:   core.String(p.cfg.NotifyURL),
		Amount: &h5.Amount{
			Total: core.Int64(req.Amount),
		},
	}
	resp, result, err := p.has.Prepay(ctx, prepayRequest)
	if err != nil {
		return nil, err
	}
	if result.Response.StatusCode != http.StatusOK {
		return nil, errors.New(result.Response.Status)
	}
	return &WechatPrepayResp{H5Url: *resp.H5Url}, nil
}

// prepayNative 生成支付链接参数code_url，然后将该参数值生成二维码图片展示给用户。用户在使用微信客户端扫描二维码后，可以直接跳转到微信支付页面完成支付操作
func (p *WechatPay) prepayNative(ctx context.Context, req *WechatPrepayReq) (*WechatPrepayResp, error) {
	prepayRequest := native.PrepayRequest{
		Appid:       core.String(p.cfg.AppID),
		Mchid:       core.String(p.cfg.MchID),
		Description: core.String(req.Description),
		OutTradeNo:  core.String(req.OutTradeNo),
		NotifyUrl:   core.String(p.cfg.NotifyURL),
		Amount: &native.Amount{
			Total: core.Int64(req.Amount),
		},
	}
	resp, result, err := p.nas.Prepay(ctx, prepayRequest)
	if err != nil {
		return nil, err
	}
	if result.Response.StatusCode != http.StatusOK {
		return nil, errors.New(result.Response.Status)
	}
	return &WechatPrepayResp{CodeUrl: *resp.CodeUrl}, nil
}

// prepayJsapi 商户系统先调用该接口在微信支付服务后台生成预支付交易单，返回正确的预支付交易会话标识后再按Native、JSAPI、APP等不同场景生成交易串调起支付
func (p *WechatPay) prepayJsapi(ctx context.Context, req *WechatPrepayReq) (*WechatPrepayResp, error) {
	prepayRequest := jsapi.PrepayRequest{
		Appid:       core.String(p.cfg.AppID),
		Mchid:       core.String(p.cfg.MchID),
		Description: core.String(req.Description),
		OutTradeNo:  core.String(req.OutTradeNo),
		NotifyUrl:   core.String(p.cfg.NotifyURL),
		Amount: &jsapi.Amount{
			Total: core.Int64(req.Amount),
		},
		Payer: &jsapi.Payer{
			Openid: core.String(req.OpenID),
		},
	}
	resp, result, err := p.jss.PrepayWithRequestPayment(ctx, prepayRequest)
	if err != nil {
		return nil, err
	}
	if result.Response.StatusCode != http.StatusOK {
		return nil, errors.New(result.Response.Status)
	}
	return &WechatPrepayResp{
		PrepayID: *resp.PrepayId,
		PrepayWithRequestPaymentResp: PrepayWithRequestPaymentResp{
			AppID:     *resp.Appid,
			TimeStamp: *resp.TimeStamp,
			NonceStr:  *resp.NonceStr,
			Package:   *resp.Package,
			SignType:  *resp.SignType,
			Sign:      *resp.PaySign,
		},
	}, nil
}

// QueryOrderByID 根据订单号查询订单信息
func (p *WechatPay) QueryOrderByID(ctx context.Context, id string) (*Transaction, error) {
	resp, result, err := p.aas.QueryOrderById(ctx, app.QueryOrderByIdRequest{
		TransactionId: core.String(id),
		Mchid:         core.String(p.cfg.MchID),
	})
	if err != nil {
		return nil, err
	}
	if result.Response.StatusCode != http.StatusOK {
		return nil, errors.New(result.Response.Status)
	}
	return (*Transaction)(resp), nil
}

// QueryOrderByOutTradeNo 根据商户内部订单号查询订单信息
func (p *WechatPay) QueryOrderByOutTradeNo(ctx context.Context, outTradeNo string) (*Transaction, error) {
	resp, result, err := p.aas.QueryOrderByOutTradeNo(ctx, app.QueryOrderByOutTradeNoRequest{
		OutTradeNo: core.String(outTradeNo),
		Mchid:      core.String(p.cfg.MchID),
	})
	if err != nil {
		return nil, err
	}
	if result.Response.StatusCode != http.StatusOK {
		return nil, errors.New(result.Response.Status)
	}
	return (*Transaction)(resp), nil
}

// HandleNotify 处理回调通知
func (p *WechatPay) HandleNotify(ctx context.Context, req *http.Request, handler func(t *Transaction) error) error {
	// 解析请求
	res := &Transaction{}
	if _, err := p.nh.ParseNotifyRequest(ctx, req, res); err != nil {
		return err
	}
	if err := handler(res); err != nil {
		return err
	}
	return nil
}

func (p *WechatPay) newNotifyHandler() (*notify.Handler, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	// 注册下载器
	err := downloader.MgrInstance().RegisterDownloaderWithPrivateKey(
		ctx, p.pk, p.cfg.MchCertSerialNo, p.cfg.MchID, p.cfg.MchAPIv3Key,
	)
	if err != nil {
		return nil, errors.New("failed to register downloader")
	}

	cm := downloader.MgrInstance().GetCertificateVisitor(p.cfg.MchID)
	return notify.NewRSANotifyHandler(p.cfg.MchAPIv3Key, verifiers.NewSHA256WithRSAVerifier(cm))
}
