package sms

import (
	"context"
	"github.com/dmzlingyin/utils/config"
	"testing"
)

func init() {
	config.SetProfile("../../config/test.json")
}

func TestSes(t *testing.T) {
	ses := New()
	err := ses.Send(context.Background(), &SendEmailArgs{
		Subject:        "测试",
		TemplateData:   `{"code":"1234"}`,
		ToEmailAddress: "dmzlingyin@163.com",
	})
	if err != nil {
		t.Fatal(err)
	}
}
