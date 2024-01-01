package cockroachdb

import (
	"reflect"

	errorz "github.com/kunitsucom/util.go/errors"
	"github.com/kunitsucom/util.go/exp/database/sql/ddl"
)

type DiffCreateTableConfig struct {
	UseAlterTableAddConstraintNotValid bool
}

type DiffCreateTableOption interface {
	apply(c *DiffCreateTableConfig)
}

func DiffCreateTableUseAlterTableAddConstraintNotValid(notValid bool) DiffCreateTableOption { //nolint:ireturn
	return &diffCreateTableConfigUseConstraintNotValid{
		useAlterTableAddConstraintNotValid: notValid,
	}
}

type diffCreateTableConfigUseConstraintNotValid struct {
	useAlterTableAddConstraintNotValid bool
}

func (o *diffCreateTableConfigUseConstraintNotValid) apply(c *DiffCreateTableConfig) {
	c.UseAlterTableAddConstraintNotValid = o.useAlterTableAddConstraintNotValid
}

//nolint:funlen,cyclop
func DiffCreateTable(before, after *CreateTableStmt, opts ...DiffCreateTableOption) (*DDL, error) {
	config := &DiffCreateTableConfig{}

	for _, opt := range opts {
		opt.apply(config)
	}

	result := &DDL{}

	switch {
	case before == nil && after != nil:
		// CREATE TABLE table_name
		result.Stmts = append(result.Stmts, after)
		return result, nil
	case before != nil && after == nil:
		// DROP TABLE table_name;
		result.Stmts = append(result.Stmts, &DropTableStmt{
			Schema: before.Schema,
			Name:   before.Name,
		})
		return result, nil
	case (before == nil && after == nil) || reflect.DeepEqual(before, after) || before.String() == after.String():
		return nil, errorz.Errorf("before: %s, after: %s: %w", before.GetPlainName(), after.GetPlainName(), ddl.ErrNoDifference)
	}

	if before.Name.Name != after.Name.Name {
		// ALTER TABLE table_name RENAME TO new_table_name;
		return nil, errorz.Errorf("ALTER TABLE %s RENAME TO %s: %w", before.Name, after.Name, ddl.ErrNotSupported)
	}

	for _, beforeConstraint := range before.Constraints {
		afterConstraint := findConstraintByName(beforeConstraint.GetName().Name, after.Constraints)
		if afterConstraint == nil {
			switch bc := beforeConstraint.(type) {
			case *IndexConstraint:
				// DROP INDEX index_name;
				result.Stmts = append(result.Stmts, &DropIndexStmt{
					Schema: before.Schema,
					Name:   bc.GetName(),
				})
			default:
				// ALTER TABLE table_name DROP CONSTRAINT constraint_name;
				result.Stmts = append(result.Stmts, &AlterTableStmt{
					Schema: before.Schema,
					Name:   before.Name,
					Action: &DropConstraint{
						Name: beforeConstraint.GetName(),
					},
				})
			}
			continue
		}
	}

	config.diffCreateTableColumn(result, before, after)

	for _, beforeConstraint := range before.Constraints {
		afterConstraint := findConstraintByName(beforeConstraint.GetName().Name, after.Constraints)
		if afterConstraint != nil {
			if beforeConstraint.PlainString() != afterConstraint.PlainString() {
				switch ac := afterConstraint.(type) {
				case *IndexConstraint:
					// DROP INDEX index_name;
					// CREATE INDEX index_name ON table_name (column_name);
					result.Stmts = append(
						result.Stmts,
						&DropIndexStmt{
							Schema: before.Schema,
							Name:   beforeConstraint.GetName(),
						},
						&CreateIndexStmt{
							Unique:    ac.Unique,
							Schema:    after.Schema,
							Name:      ac.GetName(),
							TableName: after.Name,
							Columns:   ac.Columns,
						},
					)
				default:
					// ALTER TABLE table_name DROP CONSTRAINT constraint_name;
					// ALTER TABLE table_name ADD CONSTRAINT constraint_name constraint;
					result.Stmts = append(
						result.Stmts,
						&AlterTableStmt{
							Schema: before.Schema,
							Name:   before.Name,
							Action: &DropConstraint{
								Name: beforeConstraint.GetName(),
							},
						},
						&AlterTableStmt{
							Schema: after.Schema,
							Name:   after.Name,
							Action: &AddConstraint{
								Constraint: afterConstraint,
								NotValid:   config.UseAlterTableAddConstraintNotValid,
							},
						},
					)
				}
			}
			continue
		}
	}

	for _, afterConstraint := range onlyLeftConstraint(after.Constraints, before.Constraints) {
		switch ac := afterConstraint.(type) {
		case *IndexConstraint:
			// CREATE INDEX index_name ON table_name (column_name);
			result.Stmts = append(result.Stmts, &CreateIndexStmt{
				Unique:    ac.Unique,
				Schema:    after.Schema,
				Name:      ac.GetName(),
				TableName: after.Name,
				Columns:   ac.Columns,
			})
		default:
			// ALTER TABLE table_name ADD CONSTRAINT constraint_name constraint;
			result.Stmts = append(result.Stmts, &AlterTableStmt{
				Schema: after.Schema,
				Name:   after.Name,
				Action: &AddConstraint{
					Constraint: afterConstraint,
					NotValid:   config.UseAlterTableAddConstraintNotValid,
				},
			})
		}
	}

	return result, nil
}

//nolint:funlen,cyclop
func (config *DiffCreateTableConfig) diffCreateTableColumn(ddls *DDL, before, after *CreateTableStmt) {
	for _, beforeColumn := range before.Columns {
		afterColumn := findColumnByName(beforeColumn.Name.Name, after.Columns)
		if afterColumn == nil {
			// ALTER TABLE table_name DROP COLUMN column_name;
			ddls.Stmts = append(ddls.Stmts, &AlterTableStmt{
				Schema: before.Schema,
				Name:   before.Name,
				Action: &DropColumn{
					Name: beforeColumn.Name,
				},
			})
			continue
		}

		if beforeColumn.DataType.String() != afterColumn.DataType.String() {
			// ALTER TABLE table_name ALTER COLUMN column_name SET DATA TYPE data_type;
			ddls.Stmts = append(ddls.Stmts, &AlterTableStmt{
				Schema: after.Schema,
				Name:   after.Name,
				Action: &AlterColumn{
					Name:   afterColumn.Name,
					Action: &AlterColumnSetDataType{DataType: afterColumn.DataType},
				},
			})
		}

		switch {
		case beforeColumn.Default != nil && afterColumn.Default == nil:
			// ALTER TABLE table_name ALTER COLUMN column_name DROP DEFAULT;
			ddls.Stmts = append(ddls.Stmts, &AlterTableStmt{
				Schema: after.Schema,
				Name:   after.Name,
				Action: &AlterColumn{
					Name:   afterColumn.Name,
					Action: &AlterColumnDropDefault{},
				},
			})
		case afterColumn.Default != nil && beforeColumn.Default.PlainString() != afterColumn.Default.PlainString():
			// ALTER TABLE table_name ALTER COLUMN column_name SET DEFAULT default_value;
			ddls.Stmts = append(ddls.Stmts, &AlterTableStmt{
				Schema: after.Schema,
				Name:   after.Name,
				Action: &AlterColumn{
					Name:   afterColumn.Name,
					Action: &AlterColumnSetDefault{Default: afterColumn.Default},
				},
			})
		}

		switch {
		case beforeColumn.NotNull && !afterColumn.NotNull:
			// ALTER TABLE table_name ALTER COLUMN column_name DROP NOT NULL;
			ddls.Stmts = append(ddls.Stmts, &AlterTableStmt{
				Schema: after.Schema,
				Name:   after.Name,
				Action: &AlterColumn{
					Name:   afterColumn.Name,
					Action: &AlterColumnDropNotNull{},
				},
			})
		case !beforeColumn.NotNull && afterColumn.NotNull:
			// ALTER TABLE table_name ALTER COLUMN column_name SET NOT NULL;
			ddls.Stmts = append(ddls.Stmts, &AlterTableStmt{
				Schema: after.Schema,
				Name:   after.Name,
				Action: &AlterColumn{
					Name:   afterColumn.Name,
					Action: &AlterColumnSetNotNull{},
				},
			})
		}
	}

	for _, afterColumn := range onlyLeftColumn(after.Columns, before.Columns) {
		// ALTER TABLE table_name ADD COLUMN column_name data_type;
		ddls.Stmts = append(ddls.Stmts, &AlterTableStmt{
			Schema: after.Schema,
			Name:   after.Name,
			Action: &AddColumn{
				Column: afterColumn,
			},
		})
	}
}

func onlyLeftColumn(left, right []*Column) []*Column {
	onlyLeftColumns := make([]*Column, 0)
	for _, leftColumn := range left {
		foundColumnByRight := findColumnByName(leftColumn.Name.Name, right)
		if foundColumnByRight == nil {
			onlyLeftColumns = append(onlyLeftColumns, leftColumn)
		}
	}
	return onlyLeftColumns
}

func findColumnByName(name string, columns []*Column) *Column {
	for _, column := range columns {
		if column.Name.Name == name {
			return column
		}
	}
	return nil
}

func onlyLeftConstraint(left, right Constraints) []Constraint {
	onlyLeftConstraints := make(Constraints, 0)
	for _, leftConstraint := range left {
		foundConstraintByRight := findConstraintByName(leftConstraint.GetName().Name, right)
		if foundConstraintByRight == nil {
			onlyLeftConstraints = onlyLeftConstraints.Append(leftConstraint)
		}
	}
	return onlyLeftConstraints
}

func findConstraintByName(name string, constraints []Constraint) Constraint { //nolint:ireturn
	for _, constraint := range constraints {
		if constraint.GetName().Name == name {
			return constraint
		}
	}
	return nil
}
