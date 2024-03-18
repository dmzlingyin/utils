package ioc

var global = New()

func Put(builder any, name string) {
	global.Put(builder, name)
}

func Find() {
	global.Find()
}

func Call() {
	global.Call()
}

func Range() {
	global.Range()
}
