package translation

import (
	"testing"
)

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
