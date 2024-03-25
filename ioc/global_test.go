package ioc_test

import (
	"fmt"
	"testing"

	"github.com/dmzlingyin/utils/ioc"
)

type A struct {
	b *B
}

func NewA(b *B) *A {
	return &A{
		b: b,
	}
}

func NewAA(b *B) *A {
	return &A{
		b: b,
	}
}

func NewAAA() *A {
	return &A{}
}

type B struct{}

func NewB() *B {
	return &B{}
}

func (b *B) show() {
	fmt.Println("b show")
}

func TestPut(t *testing.T) {
	ioc.Put(NewA, "test.a")
	// ioc.Put(NewA, "test.a")
	// ioc.Put(NewA, "test.a1")
	// ioc.Put(NewAA, "test.a2")
	ioc.Put(NewAAA, "test.a3")
}

func TestFind(t *testing.T) {
	ioc.Put(NewB, "test.b")
	ins := ioc.Find("test.b")
	if ins == nil {
		t.Fatal("nil instance")
	}
	if v, ok := ins.(*B); ok {
		v.show()
	}
}

func TestTryFind(t *testing.T) {
	ioc.Put(NewA, "test.a")
	ioc.Put(NewB, "test.b")
	ins, err := ioc.TryFind("test.a")
	if err != nil {
		t.Fatal(err)
	}
	if v, ok := ins.(*A); ok {
		v.b.show()
	}
}
