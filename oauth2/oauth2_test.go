package oauth2

import (
	"context"
	"testing"
)

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
