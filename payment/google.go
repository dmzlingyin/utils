package pay

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/awa/go-iap/playstore"
	"github.com/dmzlingyin/utils/log"
	"google.golang.org/api/androidpublisher/v3"
	"time"
)

type GooglePay struct {
	logger  log.Logger
	options map[string]string
}

func newGooglePay(options map[string]string) *GooglePay {
	return &GooglePay{
		options: options,
	}
}

func (p *GooglePay) GetChannel() string {
	return ChannelGoogle
}

func (p *GooglePay) Verify(ctx context.Context, args *VerifyArgs) (*VerifyRes, error) {
	key, err := p.GetKey(ctx)
	if err != nil {
		return nil, err
	}
	// 判断是订阅支付，还是普通支付
	if args.Kind == "subscription" {
		return p.VerifySub(ctx, args.Receipt, key)
	}

	client, err := playstore.New(key)
	if err != nil {
		return nil, err
	}
	res, err := client.VerifyProduct(ctx, p.options[OptionPkgName], args.ProductID, args.Receipt)
	if err != nil {
		return nil, err
	}

	// https://pkg.go.dev/google.golang.org/api/androidpublisher/v3@v0.103.0#ProductPurchase
	// 0. Purchased 1. Canceled 2. Pending
	verifyRes := &VerifyRes{}
	if res.PurchaseState == 0 {
		if res.PurchaseType != nil {
			verifyRes.Sandbox = *res.PurchaseType == 0
		}
		return verifyRes, nil
	}
	return nil, fmt.Errorf("wrong purchase state: %d", res.PurchaseState)
}

func (p *GooglePay) VerifySub(ctx context.Context, token string, key []byte) (*VerifyRes, error) {
	client, err := playstore.New(key)
	if err != nil {
		return nil, err
	}

	res, err := client.VerifySubscriptionV2(ctx, p.options[OptionPkgName], token)
	if err != nil {
		return nil, err
	}
	if res.SubscriptionState != "SUBSCRIPTION_STATE_ACTIVE" && res.SubscriptionState != "SUBSCRIPTION_STATE_CANCELED" {
		return nil, errors.New("wrong subscribe state: " + res.SubscriptionState)
	}
	if len(res.LineItems) <= 0 {
		return nil, errors.New("invalid purchase state")
	}

	st, et := p.getSubTime(res)
	return &VerifyRes{
		Sandbox:    res.TestPurchase != nil,
		OrderID:    res.LatestOrderId,
		ProductID:  res.LineItems[0].ProductId,
		StartTime:  st,
		ExpiryTime: et,
	}, nil
}

func (p *GooglePay) Create(ctx context.Context, args *CreateArgs) (res *CreateResult, err error) {
	return nil, errors.New("not yet implemented")
}

func (p *GooglePay) GetKey(ctx context.Context) ([]byte, error) {
	return base64.StdEncoding.DecodeString(p.options["json_key"])
}

func (p *GooglePay) Query(ctx context.Context, orderID string) (res *QueryResult, err error) {
	return nil, errors.New("not yet implemented")
}

func (p *GooglePay) CreateSub(ctx context.Context, args *CreateSubArgs) (res *CreateSubResult, err error) {
	return nil, errors.New("not yet implemented")
}

func (p *GooglePay) QuerySub(ctx context.Context, args *QuerySubArgs) (*SubDetail, error) {
	return nil, errors.New("not yet implemented")
}

func (p *GooglePay) Capture(ctx context.Context, orderID string, amount int32) (string, error) {
	return "", errors.New("not yet implemented")
}

func (p *GooglePay) getSubTime(sp *androidpublisher.SubscriptionPurchaseV2) (startTime, expiryTime time.Time) {
	var err error
	purchaseItem := sp.LineItems[0]
	startTime, err = time.Parse(time.RFC3339, sp.StartTime)
	if err != nil {
		log.Errorf("parse google's startTime error: %s", err)
	}
	expiryTime, err = time.Parse(time.RFC3339, purchaseItem.ExpiryTime)
	if err != nil {
		log.Errorf("parse google's expiryTime error: %s", err)
	}
	return
}

func (p *GooglePay) CreatePortal(ctx context.Context, args *CreatePortalArgs) (*CreatePortalResult, error) {
	return nil, errors.New("not yet implemented")
}
