package sqlz

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	sqlz "github.com/kunitsucom/util.go/database/sql"
	syncz "github.com/kunitsucom/util.go/sync"
)

//nolint:gochecknoglobals
var (
	columnsCache = syncz.NewMap[[]string](context.Background())
)

func TableName(tableStruct interface{}) string {
	if s, ok := tableStruct.(interface{ TableName() string }); ok {
		return s.TableName()
	}

	structType := extractStruct(tableStruct)
	return structType.Name()
}

func Columns(tableStruct interface{}, structTag string) []string {
	if s, ok := tableStruct.(interface{ Columns() []string }); ok {
		return s.Columns()
	}

	structType := extractStruct(tableStruct)
	if columns, ok := columnsCache.Load(structType); ok {
		return columns //nolint:forcetypeassert
	}

	columns := make([]string, 0, structType.NumField())
	for i := range structType.NumField() {
		structField := structType.Field(i)
		if structField.Anonymous {
			continue
		}
		tagRaw := structField.Tag.Get(structTag)
		tag := strings.Split(tagRaw, ",")[0]
		switch tag {
		case "-", "":
			continue
		}
		columns = append(columns, tag)
	}

	columnsCache.Store(structType, columns)
	return columns
}

func extractStruct(tableStruct interface{}) reflect.Type {
	structType := reflect.TypeOf(tableStruct)
	if structType.Kind() == reflect.Pointer {
		structType = structType.Elem()
	}

	if structType.Kind() != reflect.Struct {
		panic(fmt.Errorf("kind=%s expected=%s: %w", structType.Kind(), reflect.Struct, sqlz.ErrDataTypeNotSupported))
	}

	return structType
}
