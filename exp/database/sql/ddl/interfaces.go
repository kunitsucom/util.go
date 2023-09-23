package ddlz

type DDL[T any] interface {
	// String returns a pretty-printed DDL with 4 width space indent.
	String() string
	// PrettyPrint returns a pretty-printed DDL.
	PrettyPrint(indent string) string
	// Diff returns a diff of two DDLs.
	Diff(ddl DDL[T]) (DDL[T], error)
}
