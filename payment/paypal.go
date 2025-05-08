package payment

import (
	"context"
	"errors"
	"fmt"
	"github.com/dmzlingyin/utils/cast"
	"github.com/plutov/paypal/v4"
	"strconv"
	"time"
)

const (
	PaypalOrderStatusCreated   = "CREATED"   // The order was created with the specified context.
	PaypalOrderStatusCompleted = "COMPLETED" // The payment was authorized or the authorized payment was captured for the order.
	PaypalOrderStatusApproved  = "APPROVED"  // The customer approved the payment through the PayPal wallet or another form of guest or unbranded payment. For example, a card, bank account, or so on.
)

type PaypalPay struct {
	options map[string]string
	apiBase string
}

func newPaypalPay(options map[string]string) (*PaypalPay, error) {
	pay := &PaypalPay{
		options: options,
	}

	// 转换一下sandBox，赋值apiBase
	if sandbox, err := strconv.ParseBool(options[OptionSandbox]); err == nil {
		if sandbox {
			pay.apiBase = paypal.APIBaseSandBox
		} else {
			pay.apiBase = paypal.APIBaseLive
		}
	} else {
		return pay, err
	}
	return pay, nil
}

func (p *PaypalPay) Verify(ctx context.Context, args *VerifyArgs) (*VerifyRes, error) {
	order, err := p.OrderGet(ctx, args.PayID)
	if err != nil {
		return nil, err
	}
	// 验证订单状态
	if order.Status != PaypalOrderStatusCompleted {
		return nil, errors.New("invalid order status: " + order.Status)
	}
	unit := order.PurchaseUnits[0]

	// 验证订单产品
	if unit.ReferenceID != args.ProductID {
		return nil, errors.New("invalid product id: " + unit.ReferenceID)
	}
	if int32(cast.ToFloat32(unit.Amount.Value)*100) != args.Money {
		return nil, errors.New("invalid amount")
	}
	return &VerifyRes{}, nil
}

func (p *PaypalPay) Create(ctx context.Context, args *CreateArgs) (res *CreateResult, err error) {
	client, err := p.getClient(ctx)
	if err != nil {
		return nil, err
	}

	unit := paypal.PurchaseUnitRequest{
		Description: args.Description,
		CustomID:    args.CustomerID,
		Amount: &paypal.PurchaseUnitAmount{
			Currency: "USD",
			Value:    fmt.Sprintf("%.2f", float32(args.Money)/100),
		},
	}

	appCtx := &paypal.ApplicationContext{
		ReturnURL: args.ReturnURL,
	}
	order, err := client.CreateOrder(ctx, paypal.OrderIntentCapture, []paypal.PurchaseUnitRequest{unit}, nil, appCtx)
	if err != nil {
		return nil, err
	}

	var payURL string
	for _, link := range order.Links {
		if link.Rel == "approve" {
			payURL = link.Href
			break
		}
	}

	return &CreateResult{
		OrderID: order.ID,
		CodeURL: payURL,
	}, nil
}

func (p *PaypalPay) Capture(ctx context.Context, orderID string, amount int32) (string, error) {
	client, err := p.getClient(ctx)
	if err != nil {
		return "", err
	}

	// 查验订单状态
	order, err := p.OrderGet(ctx, orderID)
	if err != nil {
		return "", err
	}
	value := int32(cast.ToFloat32(order.PurchaseUnits[0].Amount.Value) * 100)
	if order.Status != paypal.OrderStatusApproved || value != amount {
		return "", errors.New("invalid order detail")
	}

	// 捕获订单
	captureRes, err := client.CaptureOrder(ctx, orderID, paypal.CaptureOrderRequest{})
	if err != nil {
		return "", err
	}
	return captureRes.Status, nil
}

func (p *PaypalPay) OrderGet(ctx context.Context, orderID string) (res *paypal.Order, err error) {
	c, err := p.getClient(ctx)
	if err != nil {
		return nil, err
	}

	return c.GetOrder(ctx, orderID)
}

func (p *PaypalPay) Query(ctx context.Context, orderID string) (res *QueryResult, err error) {
	order, err := p.OrderGet(ctx, orderID)
	unit := order.PurchaseUnits[0]

	return &QueryResult{
		Status:  order.Status,
		Money:   int32(cast.ToFloat32(unit.Amount.Value) * 100),
		OrderId: order.ID,
	}, err
}

func (p *PaypalPay) CreateSub(ctx context.Context, args *CreateSubArgs) (*CreateSubResult, error) {
	client, err := p.getClient(ctx)
	if err != nil {
		return nil, err
	}

	nsb := paypal.SubscriptionBase{
		PlanID:      args.PlanID,
		AutoRenewal: true, // 开启自动续费
		ApplicationContext: &paypal.ApplicationContext{
			ReturnURL: args.ReturnURL + "?biz_id=" + args.BizID,
		},
	}
	if args.StartTime.Local().After(time.Now()) {
		st := paypal.JSONTime(args.StartTime)
		nsb.StartTime = &st
	}

	resp, err := client.CreateSubscription(ctx, nsb)
	if err != nil {
		return nil, err
	}
	return &CreateSubResult{
		SubID:  resp.ID,
		PayURL: resp.Links[0].Href,
	}, nil
}

func (p *PaypalPay) QuerySub(ctx context.Context, args *QuerySubArgs) (*SubDetail, error) {
	client, err := p.getClient(ctx)
	if err != nil {
		return nil, err
	}
	res, err := client.GetSubscriptionDetails(ctx, args.SubID)
	if err != nil {
		return nil, err
	}

	ces := res.BillingInfo.CycleExecutions
	var cyclesCompleted int
	for _, ce := range ces {
		if ce.TenureType == "REGULAR" {
			cyclesCompleted = ce.CyclesCompleted
			break
		}
	}
	return &SubDetail{
		PlanID:          res.PlanID,
		CyclesCompleted: int32(cyclesCompleted),
		Status:          string(res.SubscriptionStatus),
		SubID:           args.SubID,
		LastPaymentTime: res.BillingInfo.LastPayment.Time,
		NextBillingTime: res.BillingInfo.NextBillingTime,
	}, nil
}

func (p *PaypalPay) getClient(ctx context.Context) (*paypal.Client, error) {
	client, err := paypal.NewClient(p.options[OptionClientId], p.options[OptionSecretId], p.apiBase)
	if err != nil {
		return nil, err
	}
	_, err = client.GetAccessToken(ctx)
	return client, err
}

func (p *PaypalPay) CreatePortal(ctx context.Context, args *CreatePortalArgs) (*CreatePortalResult, error) {
	return nil, errors.New("not yet implemented")
}
