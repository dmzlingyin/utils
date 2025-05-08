package payment

import (
	"github.com/dmzlingyin/utils/config"
	"testing"
)

func init() {
	config.SetProfile("../config/test.json")
}

func TestAlipayPay(t *testing.T) {
	pay, err := NewAlipay()
	if err != nil {
		t.Fatal(err)
	}
	orderStr, err := pay.Pay(&AliPayReq{
		OutTradeNo: "12345678",
		Amount:     "1.00",
		Subject:    "测试",
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(orderStr)
}
