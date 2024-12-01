package oauth2

import (
	"context"
	"github.com/dmzlingyin/utils/config"
	"testing"
)

func init() {
	config.SetProfile("../config/test.json")
}

func TestNew(t *testing.T) {
	p := New(TypeGoogle)
	token, user, err := p.Authorize(context.Background(), &AuthArgs{
		Type: TypeCasdoor,
		Code: "xxxx",
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(token, user)
}

func TestWechat(t *testing.T) {
	p := NewWechat()
	token, user, err := p.Authorize(context.Background(), &AuthArgs{
		Type: TypeWechat,
		Code: "021fbKkl26JnCe44Pnll2I8kTb1fbKkF",
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(token, user)
}
