package sms

import (
	"github.com/dmzlingyin/utils/config"
	"github.com/dmzlingyin/utils/misc"
	"testing"
)

func init() {
	config.SetProfile("../../config/test.json")
}

func TestSend(t *testing.T) {
	sms := New()
	err := sms.Send("xxx", misc.RandStr(6))
	if err != nil {
		t.Fatal(err)
	}
}
