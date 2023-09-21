package parser

type DMLStatementType int

const (
	DMLStatementTypeUnknown DMLStatementType = iota
	DMLStatementTypeSelect
	DMLStatementTypeInsert
	DMLStatementTypeUpdate
	DMLStatementTypeDelete
)

type DML struct {
	StatementType DMLStatementType
}
