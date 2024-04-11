package oauth2

import (
	"context"
	"github.com/cuigh/auxo/config"
	"testing"
)

func init() {
	config.AddFolder("../../config")
}

func TestApple(t *testing.T) {
	apple := NewApple()
	token, user, err := apple.Authorize(context.TODO(), "c0e76b77bc69845af92bdece1cedf1618.0.rrxz.UR4PjkoHYHiyXHz0BNYv1A")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", token)
	t.Logf("%+v", user)
}

func TestDecodeToken(t *testing.T) {
	a := NewApple()
	token := "eyJraWQiOiJXNldjT0tCIiwiYWxnIjoiUlMyNTYifQ.eyJpc3MiOiJodHRwczovL2FwcGxlaWQuYXBwbGUuY29tIiwiYXVkIjoiYWkuYWl0dWJvLndlYi5jcmVhdG9yIiwiZXhwIjoxNjc3NzQ5MjQ1LCJpYXQiOjE2Nzc2NjI4NDUsInN1YiI6IjAwMDE3OS4yNmJkYzYyZjUyNWE0YmI0YjNmNTg0MGFlNmY0OTY1Yy4wNzE5IiwiYXRfaGFzaCI6IklGX1hqLU1JdFY4WnVEb0xIcVJGM1EiLCJlbWFpbCI6ImRtemxpbmd5aW5AMTYzLmNvbSIsImVtYWlsX3ZlcmlmaWVkIjoidHJ1ZSIsImF1dGhfdGltZSI6MTY3NzY2Mjc4NCwibm9uY2Vfc3VwcG9ydGVkIjp0cnVlfQ.nidjkPvdrNVHoV324ENGCVGWe5gq5zeSOxhqpo8gLrgVQtSKu0iaquHrHJoqfbO-z7XvCa8nBYFHPCFmeKc8d_C0pUv4g8T8XAVaWTBRRl5zIKm5v1WpbN0YIbhHddxluAx0Vi-IZj9OICJuuRXhpa4krP-zP1bSSLp8M3RLO-FD4_52KP8MzwWEqlCiaydW_btUlkGrhOedTyBH7nqoo69ob4J8tge6Yi-bmmIEmAY9fXeIuGZQBRvmO3ZHYsdKmYN5v8gZlgrGryGmP-aCZ3NXaZEwrNIQ4fL_Ny3Qui83RM2za0ZxX4dgFxPTk1qUk8pjjXRC19vMsVbs45KEMw"
	a.decoder.Decode(token)
}
