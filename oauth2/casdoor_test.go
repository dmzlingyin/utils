package oauth2

import (
	"testing"
)

func TestCasdoor(t *testing.T) {
	c := NewCasdoor()
	c.Authorize("f055a6a87e666ba83d67")
}
