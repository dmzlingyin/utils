package ioc

var global = New()

func Put(builder any, name string) {
	global.Put(builder, name)
}

func Find(name string) any {
	return global.Find(name)
}

func TryFind(name string) (any, error) {
	return global.TryFind(name)
}

func Call(f any) error {
	return global.Call(f)
}

func Range(f func(name string, ins any) bool) {
	global.Range(f)
}
