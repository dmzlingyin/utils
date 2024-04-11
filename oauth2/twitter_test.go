package oauth2

import (
	"context"
	"testing"
)

func TestTwitter(t *testing.T) {
	code := "X25FY0Q4WkVUX0FOZFFTZjlnS1VxN2JRUkVabFRDTlE4X2JRNldueGpvZHJrOjE2OTgzNzk3MjAyNDk6MToxOmFjOjE"
	d := NewTwitter()
	token, user, err := d.Authorize(context.Background(), code)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("token: %+v, user: %+v", token, user)
}
