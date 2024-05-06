package misc

import "testing"

func TestSha256(t *testing.T) {
	str := "hello world"
	actual := Sha256(str)
	if actual != "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9" {
		t.Error("sha256 wrong")
	}
}
