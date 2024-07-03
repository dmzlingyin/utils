package pay

import (
	"context"
	"errors"
	"github.com/awa/go-iap/appstore"
	"github.com/dmzlingyin/utils/config"
)

type ApplePay struct {
	password string
}

func NewApplePay() *ApplePay {
	return &ApplePay{
		password: config.GetString("pay.apple.password"),
	}
}

func (p *ApplePay) GetChannel() string {
	return ChannelApple
}

func (p *ApplePay) Verify(ctx context.Context, args *VerifyArgs) (*VerifyRes, error) {
	req := appstore.IAPRequest{
		ReceiptData: args.Receipt,
		Password:    p.password,
	}
	res := &appstore.IAPResponse{}
	err := appstore.New().Verify(ctx, req, res)
	if err != nil {
		return nil, err
	}

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
