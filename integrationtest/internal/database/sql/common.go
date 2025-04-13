package sqltest

import (
	"fmt"
	"reflect"
	"strings"
)

type TestUser struct {
	ID         int     `db:"id"`
	Name       string  `db:"name"`
	Profile    *string `db:"profile"`
	NullString *string `db:"null_string"`
}

func (u *TestUser) GoString() string {
	typ := reflect.TypeOf(*u)
	val := reflect.ValueOf(*u)
	elems := make([]string, typ.NumField())
	for i := range typ.NumField() {
		elems[i] = fmt.Sprintf("%s:%#v", typ.Field(i).Name, val.Field(i))
	}
	return fmt.Sprintf("&%s{%s}", typ, strings.Join(elems, ", "))
}
