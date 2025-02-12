package net

import (
	"net/http"
	"testing"
)

func TestRequest(t *testing.T) {
	type Data struct {
		Fact   string `json:"fact"`
		Length int    `json:"length"`
	}
	res, err := Request[Data](http.MethodGet, "https://catfact.ninja/fact", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(res)
}
