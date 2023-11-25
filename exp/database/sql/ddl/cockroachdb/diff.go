package cockroachdb

import (
	"reflect"

	errorz "github.com/kunitsucom/util.go/errors"
	"github.com/kunitsucom/util.go/exp/database/sql/ddl"
	"github.com/kunitsucom/util.go/must"
)

//nolint:funlen,cyclop,gocognit
func Diff(before, after *DDL) (*DDL, error) {
	result := &DDL{}

	switch {
	case before == nil && after != nil:
		result.Stmts = append(result.Stmts, after.Stmts...)
		return result, nil
	case before != nil && after == nil:
		for _, stmt := range before.Stmts {
			switch s := stmt.(type) {
			case *CreateTableStmt:
				result.Stmts = append(result.Stmts, &DropTableStmt{
					Name: s.Name,
				})
			case *CreateIndexStmt:
				result.Stmts = append(result.Stmts, &DropIndexStmt{
					Name: s.Name,
				})
			default:
				return nil, errorz.Errorf("%s: %T: %w", s.GetPlainName(), s, ddl.ErrNotSupported)
			}
		}
		return result, nil
	case (before == nil && after == nil) || reflect.DeepEqual(before, after) || before.String() == after.String():
		return nil, ddl.ErrNoDifference
	}

	// DROP TABLE table_name;
	for _, stmt := range onlyLeftStmt(before, after) {
		switch beforeStmt := stmt.(type) {
		case *CreateTableStmt:
			result.Stmts = append(result.Stmts, &DropTableStmt{
				Name: beforeStmt.Name,
			})
		case *CreateIndexStmt:
			result.Stmts = append(result.Stmts, &DropIndexStmt{
				Name: beforeStmt.Name,
			})
		default:
			return nil, errorz.Errorf("%s: %T: %w", beforeStmt.GetPlainName(), beforeStmt, ddl.ErrNotSupported)
		}
	}

	// CREATE TABLE table_name
	for _, stmt := range onlyLeftStmt(after, before) {
		switch afterStmt := stmt.(type) {
		case *CreateTableStmt:
			result.Stmts = append(result.Stmts, afterStmt)
		case *CreateIndexStmt:
			result.Stmts = append(result.Stmts, afterStmt)
		default:
			return nil, errorz.Errorf("%s: %T: %w", afterStmt.GetPlainName(), afterStmt, ddl.ErrNotSupported)
		}
	}

	// ALTER TABLE table_name ...
	// DROP INDEX index_name; CREATE INDEX index_name ...
	for _, beforeStmt := range before.Stmts {
		switch beforeStmt := beforeStmt.(type) { //nolint:gocritic
		case *CreateTableStmt:
			if afterStmt := findStmtByTypeAndName(beforeStmt, after.Stmts); afterStmt != nil {
				afterStmt := afterStmt.(*CreateTableStmt) //nolint:forcetypeassert
				// MEMO: in this case, DiffCreateTable does not return error
				alterStmt := must.One(DiffCreateTable(beforeStmt, afterStmt))
				result.Stmts = append(result.Stmts, alterStmt.Stmts...)
			}
		case *CreateIndexStmt:
			if afterStmt := findStmtByTypeAndName(beforeStmt, after.Stmts); afterStmt != nil {
				afterStmt := afterStmt.(*CreateIndexStmt) //nolint:forcetypeassert
				if beforeStmt.PlainString() != afterStmt.PlainString() {
					result.Stmts = append(result.Stmts,
						&DropIndexStmt{
							Schema: beforeStmt.Schema,
							Name:   beforeStmt.Name,
						},
						afterStmt,
					)
				}
			}
		}
	}

	return result, nil
}

func onlyLeftStmt(left, right *DDL) []Stmt {
	result := make([]Stmt, 0)

	for _, stmt := range left.Stmts {
		if findStmtByTypeAndName(stmt, right.Stmts) == nil {
			result = append(result, stmt)
		}
	}

	return result
}

func findStmtByTypeAndName(stmt Stmt, stmts []Stmt) Stmt { //nolint:ireturn
	for _, s := range stmts {
		if reflect.TypeOf(stmt) == reflect.TypeOf(s) && stmt.GetPlainName() == s.GetPlainName() {
			return s
		}
	}
	return nil
}
