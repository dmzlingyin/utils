package oauth2

import (
	"context"
	"testing"
)

func TestDiscord(t *testing.T) {
	code := "SiA6z6qkOtEiGJHG6lOOz42S7Roxe4"
	d := NewDiscord()
	token, user, err := d.Authorize(context.Background(), code)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("token: %+v, user: %+v", token, user)
}
