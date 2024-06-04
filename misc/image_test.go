package misc

import (
	"os"
	"testing"
)

func TestGetImageInfo(t *testing.T) {
	r, err := os.Open("./test.jpg")
	if err != nil {
		t.Fatal(err)
	}
	img, err := GetImageInfo(r)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(img.Width, img.Height)
}
