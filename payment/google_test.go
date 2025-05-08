package payment

import (
	"context"
	"testing"
)

func TestGooglePayParseNotify(t *testing.T) {
	g, err := NewGooglePay()
	if err != nil {
		t.Fatal(err)
	}
	body := []byte("{\"message\":{\"data\":\"eyJ2ZXJzaW9uIjoiMS4wIiwicGFja2FnZU5hbWUiOiJjb20uYWl0dWJvLmltYWdlLmNyZWF0b3IiLCJldmVudFRpbWVNaWxsaXMiOiIxNzQ2NjgzODg1NDYwIiwic3Vic2NyaXB0aW9uTm90aWZpY2F0aW9uIjp7InZlcnNpb24iOiIxLjAiLCJub3RpZmljYXRpb25UeXBlIjozLCJwdXJjaGFzZVRva2VuIjoiY2tsaWhkYWNlY25sbWFhY2NhY2hnbWljLkFPLUoxT3lZYXlic0VFbll6X0Y3Z1hrT3B4Tk9GY1ltUng0OXphNXBVdFRmMzRfVHk5X20xakd4eXpLUmE5UHZtREJWeVhUa28xc3pRWUc3aW5IUkl2Z0h5MGdKd3hQSGNXcHBDOU1lZ25maWNfZUhxaHY5WmVjIiwic3Vic2NyaXB0aW9uSWQiOiJzdWJzY3JpYmVfYW5udWFsIn19\",\"messageId\":\"14697227525191478\",\"message_id\":\"14697227525191478\",\"publishTime\":\"2025-05-08T05:58:05.617Z\",\"publish_time\":\"2025-05-08T05:58:05.617Z\"},\"subscription\":\"projects/aitubo/subscriptions/subs\"}")
	res, err := g.ParseNotify(context.Background(), body)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(res)
}
