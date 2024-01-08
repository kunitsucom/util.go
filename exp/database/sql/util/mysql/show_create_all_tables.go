package mysql

import (
	"context"
	"database/sql"
	"fmt"

	sqlz "github.com/kunitsucom/util.go/database/sql"
	errorz "github.com/kunitsucom/util.go/errors"
)

type sqlQueryerContext = interface {
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
}

type showCreateAllTablesConfig struct {
	database string
}

type ShowCreateAllTablesOption interface {
	apply(cfg *showCreateAllTablesConfig)
}

type showCreateAllTablesOptionDatabase struct{ database string }

func (o *showCreateAllTablesOptionDatabase) apply(config *showCreateAllTablesConfig) {
	config.database = o.database
}

func WithShowCreateAllTablesOptionSchema(database string) ShowCreateAllTablesOption { //nolint:ireturn
	return &showCreateAllTablesOptionDatabase{database: database}
}

func ShowCreateAllTables(ctx context.Context, db sqlQueryerContext, opts ...ShowCreateAllTablesOption) (query string, err error) {
	dbz := sqlz.NewDB(db)

	cfg := new(showCreateAllTablesConfig)
	for _, opt := range opts {
		opt.apply(cfg)
	}

	type TableName struct {
		TableName string `db:"TABLE_NAME"`
	}

	q := "SELECT TABLE_NAME FROM information_schema.tables"
	if cfg.database != "" {
		q += fmt.Sprintf(" WHERE table_schema = `%s`", cfg.database)
	} else {
		q += " WHERE table_schema = database()"
	}

	tableNames := new([]*TableName)
	if err := dbz.QueryContext(ctx, tableNames, q); err != nil {
		return "", errorz.Errorf("dbz.QueryContext: q=%s: %w", q, err)
	}

	type CreateStatement struct {
		TableName       string `db:"Table"`
		CreateStatement string `db:"Create Table"`
	}
	for _, tableName := range *tableNames {
		createTableStmt := new(CreateStatement)
		q := fmt.Sprintf("SHOW CREATE TABLE `%s`", tableName.TableName)
		if err := dbz.QueryContext(ctx, createTableStmt, q); err != nil {
			return "", errorz.Errorf("dbz.QueryContext: q=%s: %w", q, err)
		}
		query += createTableStmt.CreateStatement
	}

	return query, nil
}
