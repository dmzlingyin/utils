package payment

import (
	"context"
	"encoding/json"
	"github.com/awa/go-iap/appstore"
	"github.com/awa/go-iap/appstore/api"
	"github.com/dmzlingyin/utils/config"
	"os"
	"time"
)

type VerifyApplePayArgs struct {
	TransactionID string
}

type VerifyApplePayRes struct {
	Sandbox       bool      // 是否为沙盒环境
	TransactionID string    // 交易ID
	ProductID     string    // 产品ID
	StartTime     time.Time // 订阅开始时间
	ExpiryTime    time.Time // 订阅到期时间
}

type ApplePayNotification struct {
	Type                  string    `map:"type"`
	UUID                  string    `map:"uuid"`
	TransactionID         string    `map:"tran_id"`
	OriginalTransactionID string    `map:"org_tran_id"`
	ProductID             string    `map:"product_id"`
	StartTime             time.Time `map:"start"`
	ExpiryTime            time.Time `map:"expiry"`
	Sandbox               bool      `map:"sandbox"`
}

type ApplePay struct {
	apiClient      *api.StoreClient
	appstoreClient *appstore.Client
}

func NewApplePay() (*ApplePay, error) {
	key, err := os.ReadFile(config.GetString("pay.apple.key_path"))
	if err != nil {
		return nil, err
	}
	cfg := &api.StoreConfig{
		KeyContent: key,
		KeyID:      config.GetString("pay.apple.key_id"),
		BundleID:   config.GetString("pay.apple.bundle_id"),
		Issuer:     config.GetString("pay.apple.issuer"),
		Sandbox:    config.GetBool("pay.apple.sandbox"),
	}
	return &ApplePay{
		apiClient:      api.NewStoreClient(cfg),
		appstoreClient: appstore.New(),
	}, nil
}

func (a *ApplePay) Verify(ctx context.Context, args *VerifyApplePayArgs) (*VerifyApplePayRes, error) {
	rsp, err := a.apiClient.GetTransactionInfo(ctx, args.TransactionID)
	if err != nil {
		return nil, err
	}
	transaction, err := a.apiClient.ParseSignedTransaction(rsp.SignedTransactionInfo)
	if err != nil {
		return nil, err
	}
	// 包装结果
	return &VerifyApplePayRes{
		Sandbox:       transaction.Environment == api.Sandbox,
		TransactionID: transaction.TransactionID,
		ProductID:     transaction.ProductID,
		StartTime:     time.UnixMilli(transaction.PurchaseDate),
		ExpiryTime:    time.UnixMilli(transaction.ExpiresDate),
	}, nil
}

func (a *ApplePay) ParseNotify(ctx context.Context, body []byte) (*ApplePayNotification, error) {
	var signedPayload appstore.SubscriptionNotificationV2SignedPayload
	if err := json.Unmarshal(body, &signedPayload); err != nil {
		return nil, err
	}

	np := appstore.SubscriptionNotificationV2DecodedPayload{}
	if err := a.appstoreClient.ParseNotificationV2WithClaim(signedPayload.SignedPayload, &np); err != nil {
		return nil, err
	}
	tp := appstore.JWSTransactionDecodedPayload{}
	if err := a.appstoreClient.ParseNotificationV2WithClaim(string(np.Data.SignedTransactionInfo), &tp); err != nil {
		return nil, err
	}

	return &ApplePayNotification{
		Type:                  string(np.NotificationType),
		UUID:                  np.NotificationUUID,
		TransactionID:         tp.TransactionId,
		OriginalTransactionID: tp.OriginalTransactionId,
		ProductID:             tp.ProductId,
		StartTime:             time.UnixMilli(tp.PurchaseDate),
		ExpiryTime:            time.UnixMilli(tp.ExpiresDate),
		Sandbox:               tp.Environment == appstore.Sandbox,
	}, nil
}
