package cos

import (
	"net/http"
	"testing"
)

func TestCosPutObject(t *testing.T) {
	c, err := NewCosClient()
	if err != nil {
		t.Fatal(err)
	}

	res, err := http.Get("http://xxx.jpg")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	if err = c.PutObject("test/xxx.jpg", res.Body); err != nil {
		t.Fatal(err)
	}
}

func TestCosPutFromFile(t *testing.T) {
	c, err := NewCosClient()
	if err != nil {
		t.Fatal(err)
	}
	if err = c.PutFromFile("./cos.go", "test/cos.go"); err != nil {
		t.Fatal(err)
	}
}
