package oss

import (
	"io"
	"net/http"
	"testing"
)

func TestOss(t *testing.T) {
	oss, err := NewClient()
	if err != nil {
		t.Fatal(err)
	}

	res, err := http.Get("")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	if err = oss.PutObject("test", io.Reader(res.Body)); err != nil {
		t.Fatal(err)
	}
}
