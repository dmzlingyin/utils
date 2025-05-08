package payment

import (
	"context"
	"errors"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/client"
	"time"
)

type StripePay struct {
	client    *client.API
	cancelURL string
}

func newStripePay(options map[string]string) (*StripePay, error) {
	pay := &StripePay{
		client:    client.New(options["key"], nil),
		cancelURL: options["cancel_url"],
	}
	return pay, nil
}

func (p *StripePay) Verify(ctx context.Context, args *VerifyArgs) (*VerifyRes, error) {
	return &VerifyRes{}, nil
}

func (p *StripePay) Create(ctx context.Context, args *CreateArgs) (res *CreateResult, err error) {
	cancelURL := args.CancelURL
	if cancelURL == "" {
		cancelURL = p.cancelURL
	}
	params := &stripe.CheckoutSessionParams{
		SuccessURL:   stripe.String(args.ReturnURL + "?session_id={CHECKOUT_SESSION_ID}"),
		CancelURL:    stripe.String(cancelURL),
		Mode:         stripe.String(string(stripe.CheckoutSessionModePayment)),
		AutomaticTax: &stripe.CheckoutSessionAutomaticTaxParams{Enabled: stripe.Bool(true)},
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Quantity: stripe.Int64(1),
				Price:    stripe.String(args.PriceID),
			},
		},
		PaymentIntentData: &stripe.CheckoutSessionPaymentIntentDataParams{
			CaptureMethod: stripe.String("manual"),
		},
	}
	if args.CustomerID != "" {
		params.Customer = stripe.String(args.CustomerID)
		params.CustomerUpdate = &stripe.CheckoutSessionCustomerUpdateParams{
			Address: stripe.String("auto"),
		}
	}

	s, err := p.client.CheckoutSessions.New(params)
	if err != nil {
		return nil, err
	}

	res = &CreateResult{
		OrderID: s.ID,
		CodeURL: s.URL,
	}
	return
}

func (p *StripePay) Capture(ctx context.Context, sessionID string, amount int32) (string, error) {
	s, err := p.client.CheckoutSessions.Get(sessionID, nil)
	if err != nil {
		return "", err
	}

	// 校验session状态
	if s.Status != stripe.CheckoutSessionStatusComplete || s.AmountSubtotal != int64(amount) {
		return "", errors.New("invalid sessionID")
	}
	// 捕获session
	pi, err := p.client.PaymentIntents.Capture(s.PaymentIntent.ID, nil)
	if err != nil {
		return "", err
	}

	if string(pi.Status) == "succeeded" {
		return "COMPLETED", nil
	}
	return "", errors.New("invalid payment status")
}

func (p *StripePay) Query(ctx context.Context, orderID string) (res *QueryResult, err error) {
	res = &QueryResult{}
	s, err := p.client.CheckoutSessions.Get(orderID, nil)
	if err != nil {
		return res, err
	}

	res.Status = string(s.PaymentStatus)
	if res.Status == "paid" {
		res.Status = "SUCCESS"
	}
	res.Money = int32(s.AmountSubtotal)
	res.OrderId = orderID
	if s.PaymentIntent != nil {
		res.OrderId = s.PaymentIntent.ID
	}

	return
}

func (p *StripePay) CreateSub(ctx context.Context, args *CreateSubArgs) (*CreateSubResult, error) {
	cancelURL := args.CancelURL
	if cancelURL == "" {
		cancelURL = p.cancelURL
	}
	params := &stripe.CheckoutSessionParams{
		Customer:            stripe.String(args.CustomerID),
		ClientReferenceID:   stripe.String(args.BizID),
		SuccessURL:          stripe.String(args.ReturnURL + "?session_id={CHECKOUT_SESSION_ID}&biz_id=" + args.BizID),
		CancelURL:           stripe.String(cancelURL),
		Mode:                stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		AutomaticTax:        &stripe.CheckoutSessionAutomaticTaxParams{Enabled: stripe.Bool(true)},
		AllowPromotionCodes: stripe.Bool(args.AllowPromotionCodes),
		CustomerUpdate: &stripe.CheckoutSessionCustomerUpdateParams{
			Address: stripe.String("auto"),
		},
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(args.PlanID),
				Quantity: stripe.Int64(1),
			},
		},
	}

	res := &CreateSubResult{}
	if args.CustomerID == "" {
		cusID, err := p.createCustomer(args.BizUserID)
		if err == nil {
			params.Customer = stripe.String(cusID)
			res.CustomerID = cusID
		}
	}

	s, err := p.client.CheckoutSessions.New(params)
	if err != nil {
		return nil, err
	}

	res.PayURL = s.URL
	res.SessionID = s.ID
	return res, nil
}

func (p *StripePay) createCustomer(bizUserID string) (string, error) {
	params := &stripe.CustomerParams{
		Params: stripe.Params{IdempotencyKey: stripe.String(bizUserID)},
	}
	c, err := p.client.Customers.New(params)
	if err != nil {
		return "", err
	}
	return c.ID, nil
}

func (p *StripePay) QuerySub(_ context.Context, args *QuerySubArgs) (*SubDetail, error) {
	if args.SessionID != "" {
		subID, err := p.getSubID(args.SessionID)
		if err != nil {
			return nil, err
		}
		args.SubID = subID
	}

	s, err := p.client.Subscriptions.Get(args.SubID, nil)
	if err != nil {
		return nil, err
	}

	return &SubDetail{
		SubID:           args.SubID,
		Status:          string(s.Status),
		LastPaymentTime: time.Unix(s.CurrentPeriodStart, 0),
		NextBillingTime: time.Unix(s.CurrentPeriodEnd, 0),
	}, nil
}

func (p *StripePay) getSubID(sessionID string) (string, error) {
	s, err := p.client.CheckoutSessions.Get(sessionID, nil)
	if err != nil {
		return "", err
	}

	if s == nil || s.Subscription == nil {
		return "", errors.New("invalid sessionID")
	}

	return s.Subscription.ID, nil
}

func (p *StripePay) CreatePortal(_ context.Context, args *CreatePortalArgs) (*CreatePortalResult, error) {
	params := &stripe.BillingPortalSessionParams{
		Customer:  stripe.String(args.CustomerID),
		ReturnURL: stripe.String(args.ReturnURL),
	}
	result, err := p.client.BillingPortalSessions.New(params)
	if err != nil {
		return nil, err
	}
	return &CreatePortalResult{URL: result.URL}, nil
}
