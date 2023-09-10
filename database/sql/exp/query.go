package sqlz

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	sqlz "github.com/kunitsucom/util.go/database/sql"
)

const defaultStructTag = "db"

//nolint:gochecknoglobals
var (
	columnsCache sync.Map
)

type NewQueryBuilderOption interface {
	apply(*queryBuilder)
}

type newQueryBuilderOptionStructTag string

func (o newQueryBuilderOptionStructTag) apply(qb *queryBuilder) {
	qb.structTag = string(o)
}

func WithNewQueryBuilderOptionStructTag(tag string) NewQueryBuilderOption { //nolint:ireturn
	return newQueryBuilderOptionStructTag(tag)
}

func NewQueryBuilder(opts ...NewQueryBuilderOption) QueryBuilder { //nolint:ireturn
	qb := &queryBuilder{
		structTag: defaultStructTag,
	}

	for _, opt := range opts {
		opt.apply(qb)
	}

	return qb
}

type QueryBuilder interface {
	TableName(tableStruct interface{}) string
	Columns(tableStruct interface{}) []string

	private()
}

type queryBuilder struct {
	structTag string
}

func (q *queryBuilder) private() {}

func (q *queryBuilder) TableName(tableStruct interface{}) string {
	if s, ok := tableStruct.(interface{ TableName() string }); ok {
		return s.TableName()
	}

	structType := extractStruct(tableStruct)
	return structType.Name()
}

func (q *queryBuilder) Columns(tableStruct interface{}) []string {
	return columns(tableStruct, q.structTag)
}

//nolint:revive
func columns(tableStruct interface{}, structTag string) []string {
	if s, ok := tableStruct.(interface{ Columns() []string }); ok {
		return s.Columns()
	}

	structType := extractStruct(tableStruct)
	if columns, ok := columnsCache.Load(structType); ok {
		return columns.([]string) //nolint:forcetypeassert
	}

	var columns []string
	for i := 0; i < structType.NumField(); i++ {
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

func ResetGlobalColumnsCache() {
	columnsCache.Range(func(key, value interface{}) bool {
		columnsCache.Delete(key)
		return true
	})
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
