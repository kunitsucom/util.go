package postgres

import (
	"reflect"

	errorz "github.com/kunitsucom/util.go/errors"
	"github.com/kunitsucom/util.go/exp/database/sql/ddl"
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

	// CREATE TABLE table_name
	for _, stmt := range onlyLeftStmt(after, before) {
		switch s := stmt.(type) {
		case *CreateTableStmt:
			result.Stmts = append(result.Stmts, s)
		case *CreateIndexStmt:
			result.Stmts = append(result.Stmts, s)
		default:
			return nil, errorz.Errorf("%s: %T: %w", s.GetPlainName(), s, ddl.ErrNotSupported)
		}
	}

	// ALTER TABLE table_name ...
	// DROP INDEX index_name; CREATE INDEX index_name ...
	for _, beforeStmt := range before.Stmts {
		switch beforeStmt := beforeStmt.(type) { //nolint:gocritic
		case *CreateTableStmt:
			if afterStmt := findStmtByTypeAndName(beforeStmt, after.Stmts); afterStmt != nil {
				afterStmt := afterStmt.(*CreateTableStmt) //nolint:forcetypeassert
				alterStmt, err := DiffCreateTable(beforeStmt, afterStmt)
				if err != nil {
					return nil, errorz.Errorf("DiffCreateTable: %w", err)
				}
				result.Stmts = append(result.Stmts, alterStmt.Stmts...)
			}
		case *CreateIndexStmt:
			if afterStmt := findStmtByTypeAndName(beforeStmt, after.Stmts); afterStmt != nil {
				afterStmt := afterStmt.(*CreateIndexStmt) //nolint:forcetypeassert
				if beforeStmt.PlainString() != afterStmt.PlainString() {
					result.Stmts = append(result.Stmts,
						&DropIndexStmt{Name: beforeStmt.Name},
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
