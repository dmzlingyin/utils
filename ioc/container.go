package ioc

import (
	"reflect"

	"github.com/dmzlingyin/utils/lazy"
)

type Container struct {
	instances map[string]*instance
}

func New() *Container {
	return &Container{
		instances: make(map[string]*instance),
	}
}

func (c *Container) Put(builder any, name string) {
	if name == "" {
		panic("empty instance name")
	}
	if _, ok := c.instances[name]; ok {
		panic("duplicate instance name: " + name)
	}

	t := reflect.TypeOf(builder)
	v := reflect.ValueOf(builder)
	if t.Kind() != reflect.Func || v.IsNil() || t.NumOut() != 1 {
		panic("builder must be a function with one return value")
	}

	ins := &instance{t: t, v: v, name: name}
	ins.value = &lazy.Value[any]{
		New: func() (any, error) {
			return ins.build(c)
		},
	}
	c.instances[name] = ins
}

func (c *Container) Find() {}

func (c *Container) Call() {}

func (c *Container) Range() {}

type instance struct {
	t     reflect.Type
	v     reflect.Value
	value *lazy.Value[any]
	name  string
}

func (i *instance) build(c *Container) (any, error) {
	return nil, nil
}
