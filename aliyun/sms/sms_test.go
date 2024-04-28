package sms

import (
	"github.com/dmzlingyin/utils/config"
	"testing"
)

func TestSend(t *testing.T) {
	config.SetProfile("../../config/test.json")
	sms := New()
	err := sms.Send("xxx", misc.RandStr(6))
	if err != nil {
		t.Fatal(err)
	}
}
