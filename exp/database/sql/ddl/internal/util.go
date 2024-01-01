package internal

import (
	"fmt"
	"reflect"
	"strings"
)

func GoString(v interface{}) string {
	typ := reflect.TypeOf(v)
	if typ.Kind() != reflect.Struct {
		panic(fmt.Errorf("kind=%s expected=%s", typ.Kind(), reflect.Struct)) //nolint:goerr113
	}

	val := reflect.ValueOf(v)
	elems := make([]string, typ.NumField())
	for i := 0; typ.NumField() > i; i++ {
		elems[i] = fmt.Sprintf("%s:%#v", typ.Field(i).Name, val.Field(i))
	}
	return fmt.Sprintf("&%s{%s}", typ, strings.Join(elems, ", "))
}
