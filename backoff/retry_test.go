package backoff

import (
	"errors"
	"testing"
)

func TestRetry(t *testing.T) {
	err := Retry(3, testFunc)
	if err == nil {
		t.Fatal("error")
	}
	err = Retry(-1, testFunc)
	if err == nil {
		t.Fatal("error")
	}
}

func testFunc() error {
	return errors.New("error")
}
