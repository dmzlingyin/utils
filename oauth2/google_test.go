package oauth2

import (
	"context"
	"golang.org/x/oauth2"
	"io"
	"testing"
)

func TestToken(t *testing.T) {
	g := NewGoogle()
	token := &oauth2.Token{
		AccessToken: "ya29.a0AWY7CknnWB6XGfzVLTik0ZJFsPheO30VM9JtlmjzLWx1CNq72syGT9VD8W7WkqADGo1cnqNuEI4cQZF2cOB3es3WdyWFMa5uDtibjFnAOyUBVuXKt2UZadHXm7xAO4UnfhZH0-L0uO41ysvxc8H7w_xgglHX9gaCgYKAXsSARMSFQG1tDrphTaPAcqWvLVYYXrf54Q-wA0165",
		TokenType:   "Bearer",
	}
	res, err := g.cfg.Client(context.TODO(), token).Get(GoogleUserURL)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf(string(body))
}
