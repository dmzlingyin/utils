package config

import "testing"

func TestConfig(t *testing.T) {
	SetProfile("test")

	addr := Get("app.addr")
	if !addr.Exists() {
		t.Fatal("no addr filed")
	}
	t.Log(addr.String())

	port := Get("app.port")
	if !port.Exists() {
		t.Fatal("no port field")
	}
	t.Log(port.Int())
}
