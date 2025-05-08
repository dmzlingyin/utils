package payment

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/awa/go-iap/playstore"
	"github.com/dmzlingyin/utils/config"
	"github.com/dmzlingyin/utils/log"
	"google.golang.org/api/androidpublisher/v3"
	"os"
	"strings"
	"time"
)

const (
	SubStatusReNew     int32 = 0
	SubStatusCancelled int32 = 1
	SubStatusNone      int32 = 2
	SubStatusTest      int32 = 3
)

type VerifyGooglePayArgs struct {
	Subscription  bool
	PackageName   string
	PurchaseToken string
	ProductID     string
}

type VerifyGooglePayRes struct {
	Sandbox       bool
	TransactionID string    // 交易ID
	StartTime     time.Time // 订阅开始时间
	ExpiryTime    time.Time // 订阅到期时间
}

type GooglePayNotification struct {
	SubStatus             int32     `map:"sub_status"`
	UUID                  string    `map:"uuid"`
	TransactionID         string    `map:"tran_id"`
	OriginalTransactionID string    `map:"org_tran_id"`
	ProductID             string    `map:"product_id"`
	StartTime             time.Time `map:"start"`
	ExpiryTime            time.Time `map:"expiry"`
	Sandbox               bool      `map:"sandbox"`
}

type GooglePay struct {
	client *playstore.Client
}

func NewGooglePay() (*GooglePay, error) {
	key, err := os.ReadFile(config.GetString("pay.google.key_path"))
	if err != nil {
		return nil, err
	}
	client, err := playstore.New(key)
	if err != nil {
		return nil, err
	}
	return &GooglePay{
		client: client,
	}, nil
}

func (g *GooglePay) Verify(ctx context.Context, args *VerifyGooglePayArgs) (*VerifyGooglePayRes, error) {
	if args.Subscription {
		return g.verifySub(ctx, args)
	}

	res, err := g.client.VerifyProduct(ctx, args.PackageName, args.ProductID, args.PurchaseToken)
	if err != nil {
		return nil, err
	}
	// https://pkg.go.dev/google.golang.org/api/androidpublisher/v3@v0.103.0#ProductPurchase
	// 0. Purchased 1. Canceled 2. Pending
	if res.PurchaseState != 0 {
		return nil, fmt.Errorf("wrong purchase state: %d", res.PurchaseState)
	}
	verifyRes := &VerifyGooglePayRes{TransactionID: res.OrderId}
	if res.PurchaseType != nil {
		verifyRes.Sandbox = *res.PurchaseType == 0
	}
	return verifyRes, nil
}

func (g *GooglePay) verifySub(ctx context.Context, args *VerifyGooglePayArgs) (*VerifyGooglePayRes, error) {
	sp, err := g.client.VerifySubscriptionV2(ctx, args.PackageName, args.PurchaseToken)
	if err != nil {
		return nil, err
	}
	if sp.SubscriptionState != "SUBSCRIPTION_STATE_ACTIVE" && sp.SubscriptionState != "SUBSCRIPTION_STATE_CANCELED" {
		log.Infof("google subscription purchases: %+v", sp)
		return nil, errors.New("invalid subscription state")
	}
	if len(sp.LineItems) <= 0 {
		log.Infof("google subscription purchases: %+v", sp)
		return nil, errors.New("invalid purchase state")
	}

	st, et := g.getSubTime(sp)
	return &VerifyGooglePayRes{
		Sandbox:       sp.TestPurchase != nil,
		TransactionID: sp.LatestOrderId,
		StartTime:     st,
		ExpiryTime:    et,
	}, nil
}

// GooglePub 谷歌开发者实时通知请求体
type GooglePub struct {
	Message struct {
		Attributes struct {
			Key string `json:"key"`
		} `json:"attributes"`
		Data      string `json:"data"`
		MessageId string `json:"messageId"`
	} `json:"message"`
	Subscription string `json:"subscription"`
}

type DeveloperNotification struct {
	Version                    string                                `json:"version"`
	PackageName                string                                `json:"packageName"`
	EventTimeMillis            string                                `json:"eventTimeMillis"`
	OneTimeProductNotification *playstore.OneTimeProductNotification `json:"oneTimeProductNotification"`
	SubscriptionNotification   *playstore.SubscriptionNotification   `json:"subscriptionNotification"`
	TestNotification           *playstore.TestNotification           `json:"testNotification"`
}

func (g *GooglePay) ParseNotify(ctx context.Context, body []byte) (*GooglePayNotification, error) {
	var gp GooglePub
	if err := json.Unmarshal(body, &gp); err != nil {
		return nil, err
	}

	baseData := gp.Message.Data
	decoded, err := base64.StdEncoding.DecodeString(baseData)
	if err != nil {
		return nil, err
	}
	// 将JSON字节解析为结构体
	developerNotification := DeveloperNotification{}
	if err = json.Unmarshal(decoded, &developerNotification); err != nil {
		return nil, err
	}
	res := &GooglePayNotification{}
	if developerNotification.TestNotification != nil {
		res.SubStatus = SubStatusTest
		return res, nil
	}
	// 向google获取订单状态
	subNotification := developerNotification.SubscriptionNotification
	sp, err := g.client.VerifySubscriptionV2(ctx, developerNotification.PackageName, subNotification.PurchaseToken)
	if err != nil {
		return nil, err
	}

	// 续订
	if subNotification.NotificationType == playstore.SubscriptionNotificationTypeRenewed {
		if sp.AcknowledgementState != "ACKNOWLEDGEMENT_STATE_ACKNOWLEDGED" {
			return nil, errors.New("invalid acknowledgement state")
		}
		if len(sp.LineItems) <= 0 {
			return nil, errors.New("invalid purchase state")
		}
		res.SubStatus = SubStatusReNew
	} else if subNotification.NotificationType == playstore.SubscriptionNotificationTypeCanceled {
		res.SubStatus = SubStatusCancelled
	} else {
		res.SubStatus = SubStatusNone
	}

	st, et := g.getSubTime(sp)
	res.StartTime = st
	res.ExpiryTime = et
	res.TransactionID = sp.LatestOrderId
	res.UUID = gp.Message.MessageId
	if strings.Contains(res.TransactionID, "..") {
		res.OriginalTransactionID = strings.Split(res.TransactionID, "..")[0]
	} else {
		res.OriginalTransactionID = res.TransactionID
	}
	res.ProductID = subNotification.SubscriptionID
	return res, nil
}

func (g *GooglePay) getSubTime(sp *androidpublisher.SubscriptionPurchaseV2) (startTime, expiryTime time.Time) {
	var err error
	purchaseItem := sp.LineItems[0]
	// 此处将订阅的起始时间赋值为了最后一次付款时间(付款时间是paypal的定义),影响不大
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
