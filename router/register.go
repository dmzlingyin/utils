package router

import (
	"fmt"
	"github.com/dmzlingyin/utils/ioc"
	"github.com/dmzlingyin/utils/log"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
	"strings"
)

// Register 注册路由
func Register(g *gin.RouterGroup, prefix, name string) {
	ins, err := ioc.TryFind(prefix + name)
	if err != nil {
		panic(err)
	}
	v := reflect.Indirect(reflect.ValueOf(ins))
	if v.Kind() != reflect.Struct {
		panic("invalid handler type: " + name)
	}
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		h, ok := f.Interface().(gin.HandlerFunc)
		if !ok {
			continue
		}
		field := t.Field(i)
		if h == nil {
			panic(fmt.Sprintf("handler %s.%s isn't initialized", t.Name(), field.Name))
		}
		path := field.Tag.Get("path")
		relativePath := fmt.Sprintf("/%s%s", name, path)
		method := field.Tag.Get("method")
		registerMethod(g, method, relativePath, h)
	}
}

func registerMethod(g *gin.RouterGroup, method, rPath string, h gin.HandlerFunc) {
	switch strings.ToUpper(method) {
	case http.MethodGet:
		g.GET(rPath, h)
	case http.MethodPost:
		g.POST(rPath, h)
	case http.MethodPut:
		g.PUT(rPath, h)
	case http.MethodDelete:
		g.DELETE(rPath, h)
	case http.MethodPatch:
		g.PATCH(rPath, h)
	default:
		panic("unsupported method: " + method)
	}
	log.Infof("handler registered: [%s], path: %s", method, rPath)
}
