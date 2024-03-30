package translation

import (
	"github.com/cuigh/auxo/config"
	"testing"
)

func init() {
	config.AddFolder("../../config")
}

func TestTranslate(t *testing.T) {
	tr, err := New()
	if err != nil {
		t.Fatal(err)
	}
	res, err := tr.Translate("你好")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(res)
}
