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
	ioc.Put(NewA, "test.put")
	// ioc.Put(NewA, "test.put")
	// ioc.Put(NewA, "test.put1")
	// ioc.Put(NewAA, "test.put2")
	ioc.Put(NewAAA, "test.put3")
}

func TestTryFind(t *testing.T) {
	ioc.TryFind("")
}
