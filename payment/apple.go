package pay

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/awa/go-iap/appstore"
	"net/http"
	"net/url"
	"strings"
)

type ApplePay struct {
	options map[string]string
}

func newApplePay(options map[string]string) *ApplePay {
	return &ApplePay{
		options: options,
	}
}

func (p *ApplePay) GetChannel() string {
	return ChannelApple
}

func (p *ApplePay) Verify(ctx context.Context, args *VerifyArgs) (*VerifyRes, error) {
	if args.Kind == "subscription" {
		return p.verifySub(args.Receipt)
	}

	req := appstore.IAPRequest{
		ReceiptData: args.Receipt,
		Password:    p.options[OptionPassword],
	}
	res := &appstore.IAPResponse{}
	err := appstore.New().Verify(ctx, req, res)
	if err != nil {
		return nil, err
	}
	//p.logger.Infof("verify response: %+v", res)

	// 验证订单状态
	if res.Status != 0 {
		return nil, errors.New("invalid pay status")
	}

	// 验证产品ID
	ok := false
	for _, v := range res.Receipt.InApp {
		if args.ProductID != v.ProductID {
			continue
		}
		if args.PayID == v.TransactionID {
			ok = true
			break
		}
	}
	if !ok {
		err = errors.New("can not get receipt attributes, invalid pay id: " + args.PayID)
	}
	return &VerifyRes{Sandbox: res.Environment == "Sandbox"}, err
}

func (p *ApplePay) verifySub(receipt string) (*VerifyRes, error) {

	addr := "https://api.musicringtoneapp.com/api/v1/verify/" + p.options["package_name"]
	pd := url.Values{}
	pd.Add("receipt", receipt)

	resp, err := http.Post(addr, "application/x-www-form-urlencoded", strings.NewReader(pd.Encode()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	type subMeta struct {
		ProductID  string `json:"product_id"`
		SandBoxEnv bool   `json:"sandbox_env"`
		Subscribed bool   `json:"subscribed"`
	}
	type Resp struct {
		Code int32   `json:"code"`
		Msg  string  `json:"msg"`
		Data subMeta `json:"data"`
	}
	var subInfo Resp
	if err = json.NewDecoder(resp.Body).Decode(&subInfo); err != nil {
		return nil, err
	}

	if subInfo.Code != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("%s", subInfo.Msg))
	}
	if !subInfo.Data.Subscribed {
		return nil, errors.New("invalid subscription status")
	}

	return &VerifyRes{Sandbox: subInfo.Data.SandBoxEnv}, nil
}

func (p *ApplePay) Create(ctx context.Context, args *CreateArgs) (res *CreateResult, err error) {
	return nil, errors.New("not yet implemented")
}

func (p *ApplePay) Query(ctx context.Context, orderID string) (res *QueryResult, err error) {
	return nil, errors.New("not yet implemented")
}

func (p *ApplePay) CreateSub(ctx context.Context, args *CreateSubArgs) (res *CreateSubResult, err error) {
	return nil, errors.New("not yet implemented")
}

func (p *ApplePay) QuerySub(ctx context.Context, args *QuerySubArgs) (*SubDetail, error) {
	return nil, errors.New("not yet implemented")
}

func (p *ApplePay) Capture(ctx context.Context, orderID string, amount int32) (string, error) {
	return "", errors.New("not yet implemented")
}

func (p *ApplePay) CreatePortal(ctx context.Context, args *CreatePortalArgs) (*CreatePortalResult, error) {
	return nil, errors.New("not yet implemented")
}
