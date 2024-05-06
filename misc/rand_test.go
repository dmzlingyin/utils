package misc

import "testing"

func TestRandStr(t *testing.T) {
	for i := 0; i < 10; i++ {
		println(RandStr(6))
	}
	for i := 0; i < 10; i++ {
		println(RandStr(6, true))
	}
}
