package ioc

import (
	"errors"
	"reflect"

	"github.com/dmzlingyin/utils/lazy"
)

var TypeError = reflect.TypeOf((*error)(nil)).Elem()

type Container struct {
	instances map[string]*instance
	types     map[reflect.Type]*instance
}

func New() *Container {
	return &Container{
		instances: make(map[string]*instance),
		types:     make(map[reflect.Type]*instance),
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
	if t.Kind() != reflect.Func || t.NumOut() != 1 || v.IsNil() {
		panic("builder must be a function with one return value")
	}
	if ins := c.types[t.Out(0)]; ins != nil {
		panic("duplicate instance type: " + t.String())
	}

	ins := &instance{t: t, v: v, name: name}
	ins.value = &lazy.Value[any]{
		New: func() (any, error) {
			return ins.build(c)
		},
	}
	c.instances[name] = ins
	c.types[t.Out(0)] = ins
}

func (c *Container) Find(name string) any {
	v, _ := c.TryFind(name)
	return v
}

func (c *Container) TryFind(name string) (any, error) {
	if v, ok := c.instances[name]; ok {
		return v.new()
	}
	return nil, errors.New("not regeistered instance: " + name)
}

func (c *Container) Call(f any) error {
	t := reflect.TypeOf(f)
	v := reflect.ValueOf(f)
	if t.Kind() != reflect.Func || v.IsNil() {
		return errors.New("parameter f is not a valid function")
	}

	res, err := c.call(t, v)
	if err != nil {
		return err
	}
	if len(res) > 0 {
		if res[0].Type() == TypeError && !res[0].IsNil() {
			return res[0].Interface().(error)
		}
	}
	return nil
}

func (c *Container) call(t reflect.Type, v reflect.Value) ([]reflect.Value, error) {
	args := make([]reflect.Value, t.NumIn())
	for i := range args {
		svc, err := c.get(t.In(i))
		if err != nil {
			return nil, err
		}
		args[i] = reflect.ValueOf(svc)
	}
	return v.Call(args), nil
}

func (c *Container) get(t reflect.Type) (any, error) {
	if s := c.types[t]; s != nil {
		return s.new()
	}
	return nil, errors.New("cannot get service: " + t.Name())
}

func (c *Container) Range(f func(name string, ins any) bool) {
	for n, i := range c.instances {
		res, err := i.build(c)
		if err != nil {
			panic(err)
		}
		if !f(n, res) {
			return
		}
	}
}

type instance struct {
	t     reflect.Type
	v     reflect.Value
	value *lazy.Value[any]
	name  string
}

func (ins *instance) new() (any, error) {
	return ins.value.Get()
}

func (ins *instance) build(c *Container) (any, error) {
	res, err := c.call(ins.t, ins.v)
	if err != nil {
		return nil, err
	}
	return res[0].Interface(), nil
}
