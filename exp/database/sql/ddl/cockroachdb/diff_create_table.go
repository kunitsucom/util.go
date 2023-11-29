package cockroachdb

import (
	"reflect"

	errorz "github.com/kunitsucom/util.go/errors"
	"github.com/kunitsucom/util.go/exp/database/sql/ddl"
	"github.com/kunitsucom/util.go/exp/diff/simplediff"
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
		after.Comment = simplediff.Diff(before.String(), after.String()).String()
		result.Stmts = append(result.Stmts, after)
		return result, nil
	case before != nil && after == nil:
		// DROP TABLE table_name;
		result.Stmts = append(result.Stmts, &DropTableStmt{
			Comment: simplediff.Diff(before.String(), after.String()).String(),
			Name:    before.Name,
		})
		return result, nil
	case (before == nil && after == nil) || reflect.DeepEqual(before, after) || before.String() == after.String():
		return nil, errorz.Errorf("before: %s, after: %s: %w", before.GetNameForDiff(), after.GetNameForDiff(), ddl.ErrNoDifference)
	}

	if before.Name.StringForDiff() != after.Name.StringForDiff() {
		// ALTER TABLE table_name RENAME TO new_table_name;
		rename := &RenameTable{
			NewName: after.Name,
		}
		if rename.NewName.Schema == nil {
			rename.NewName.Schema = before.Name.Schema
		}
		result.Stmts = append(result.Stmts, &AlterTableStmt{
			Comment: simplediff.Diff(before.String(), after.String()).String(),
			Name:    before.Name,
			Action:  rename,
		})
	}

	for _, beforeConstraint := range before.Constraints {
		afterConstraint := findConstraintByName(beforeConstraint.GetName().Name, after.Constraints)
		if afterConstraint == nil {
			switch bc := beforeConstraint.(type) { //diff:ignore-line-postgres-cockroach
			case *IndexConstraint: //diff:ignore-line-postgres-cockroach
				// DROP INDEX index_name; //diff:ignore-line-postgres-cockroach
				result.Stmts = append(result.Stmts, &DropIndexStmt{ //diff:ignore-line-postgres-cockroach
					Comment: simplediff.Diff(before.String(), "").String(), //diff:ignore-line-postgres-cockroach
					Name: &ObjectName{ //diff:ignore-line-postgres-cockroach
						Schema: before.Name.Schema, //diff:ignore-line-postgres-cockroach
						Name:   bc.GetName(),       //diff:ignore-line-postgres-cockroach
					}, //diff:ignore-line-postgres-cockroach
				}) //diff:ignore-line-postgres-cockroach
			default: //diff:ignore-line-postgres-cockroach
				// ALTER TABLE table_name DROP CONSTRAINT constraint_name;
				result.Stmts = append(result.Stmts, &AlterTableStmt{
					Comment: simplediff.Diff(before.String(), after.String()).String(),
					Name:    after.Name, // ALTER TABLE RENAME TO で変更された後の可能性があるため after.Name を使用する
					Action: &DropConstraint{
						Name: beforeConstraint.GetName(),
					},
				})
			} //diff:ignore-line-postgres-cockroach
			continue
		}
	}

	config.diffCreateTableColumn(result, before, after)

	for _, beforeConstraint := range before.Constraints {
		afterConstraint := findConstraintByName(beforeConstraint.GetName().Name, after.Constraints)
		if afterConstraint != nil {
			if beforeConstraint.StringForDiff() != afterConstraint.StringForDiff() {
				switch ac := afterConstraint.(type) { //diff:ignore-line-postgres-cockroach
				case *IndexConstraint: //diff:ignore-line-postgres-cockroach
					// DROP INDEX index_name;                               //diff:ignore-line-postgres-cockroach
					// CREATE INDEX index_name ON table_name (column_name); //diff:ignore-line-postgres-cockroach
					result.Stmts = append( //diff:ignore-line-postgres-cockroach
						result.Stmts, //diff:ignore-line-postgres-cockroach
						&DropIndexStmt{ //diff:ignore-line-postgres-cockroach
							Comment: simplediff.Diff(before.String(), after.String()).String(), //diff:ignore-line-postgres-cockroach
							Name: &ObjectName{ //diff:ignore-line-postgres-cockroach
								Schema: before.Name.Schema,         //diff:ignore-line-postgres-cockroach
								Name:   beforeConstraint.GetName(), //diff:ignore-line-postgres-cockroach
							}, //diff:ignore-line-postgres-cockroach
						}, //diff:ignore-line-postgres-cockroach
						&CreateIndexStmt{ //diff:ignore-line-postgres-cockroach
							Comment: simplediff.Diff(before.String(), after.String()).String(), //diff:ignore-line-postgres-cockroach
							Unique:  ac.Unique,                                                 //diff:ignore-line-postgres-cockroach
							Name: &ObjectName{ //diff:ignore-line-postgres-cockroach
								Schema: after.Name.Schema, //diff:ignore-line-postgres-cockroach
								Name:   ac.GetName(),      //diff:ignore-line-postgres-cockroach
							}, //diff:ignore-line-postgres-cockroach
							TableName: after.Name, //diff:ignore-line-postgres-cockroach
							Columns:   ac.Columns, //diff:ignore-line-postgres-cockroach
						}, //diff:ignore-line-postgres-cockroach
					) //diff:ignore-line-postgres-cockroach
				default: //diff:ignore-line-postgres-cockroach
					// ALTER TABLE table_name DROP CONSTRAINT constraint_name;
					// ALTER TABLE table_name ADD CONSTRAINT constraint_name constraint;
					result.Stmts = append(
						result.Stmts,
						&AlterTableStmt{
							Comment: simplediff.Diff(before.String(), after.String()).String(),
							Name:    after.Name, // ALTER TABLE RENAME TO で変更された後の可能性があるため after.Name を使用する
							Action: &DropConstraint{
								Name: beforeConstraint.GetName(),
							},
						},
						&AlterTableStmt{
							Comment: simplediff.Diff(before.String(), after.String()).String(),
							Name:    after.Name,
							Action: &AddConstraint{
								Constraint: afterConstraint,
								NotValid:   config.UseAlterTableAddConstraintNotValid,
							},
						},
					)
				} //diff:ignore-line-postgres-cockroach
			}
			continue
		}
	}

	for _, afterConstraint := range onlyLeftConstraint(after.Constraints, before.Constraints) {
		switch ac := afterConstraint.(type) { //diff:ignore-line-postgres-cockroach
		case *IndexConstraint: //diff:ignore-line-postgres-cockroach
			// CREATE INDEX index_name ON table_name (column_name); //diff:ignore-line-postgres-cockroach
			result.Stmts = append(result.Stmts, &CreateIndexStmt{ //diff:ignore-line-postgres-cockroach
				Comment: simplediff.Diff(before.String(), after.String()).String(), //diff:ignore-line-postgres-cockroach
				Unique:  ac.Unique,                                                 //diff:ignore-line-postgres-cockroach
				Name: &ObjectName{ //diff:ignore-line-postgres-cockroach
					Schema: after.Name.Schema, //diff:ignore-line-postgres-cockroach
					Name:   ac.GetName(),      //diff:ignore-line-postgres-cockroach
				}, //diff:ignore-line-postgres-cockroach
				TableName: after.Name, //diff:ignore-line-postgres-cockroach
				Columns:   ac.Columns, //diff:ignore-line-postgres-cockroach
			}) //diff:ignore-line-postgres-cockroach
		default: //diff:ignore-line-postgres-cockroach
			// ALTER TABLE table_name ADD CONSTRAINT constraint_name constraint;
			result.Stmts = append(result.Stmts, &AlterTableStmt{
				Comment: simplediff.Diff(before.String(), after.String()).String(),
				Name:    after.Name,
				Action: &AddConstraint{
					Constraint: afterConstraint,
					NotValid:   config.UseAlterTableAddConstraintNotValid,
				},
			})
		} //diff:ignore-line-postgres-cockroach
	}

	if len(result.Stmts) == 0 {
		return nil, errorz.Errorf("before: %s, after: %s: %w", before.GetNameForDiff(), after.GetNameForDiff(), ddl.ErrNoDifference)
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
				Comment: simplediff.Diff(before.String(), after.String()).String(),
				Name:    after.Name, // ALTER TABLE RENAME TO で変更された後の可能性があるため after.Name を使用する
				Action: &DropColumn{
					Name: beforeColumn.Name,
				},
			})
			continue
		}

		if beforeColumn.DataType.StringForDiff() != afterColumn.DataType.StringForDiff() {
			// ALTER TABLE table_name ALTER COLUMN column_name SET DATA TYPE data_type;
			ddls.Stmts = append(ddls.Stmts, &AlterTableStmt{
				Comment: simplediff.Diff(before.String(), after.String()).String(),
				Name:    after.Name,
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
				Comment: simplediff.Diff(before.String(), after.String()).String(),
				Name:    after.Name,
				Action: &AlterColumn{
					Name:   afterColumn.Name,
					Action: &AlterColumnDropDefault{},
				},
			})
		case afterColumn.Default != nil && beforeColumn.Default.StringForDiff() != afterColumn.Default.StringForDiff():
			// ALTER TABLE table_name ALTER COLUMN column_name SET DEFAULT default_value;
			ddls.Stmts = append(ddls.Stmts, &AlterTableStmt{
				Comment: simplediff.Diff(before.String(), after.String()).String(),
				Name:    after.Name,
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
				Comment: simplediff.Diff(before.String(), after.String()).String(),
				Name:    after.Name,
				Action: &AlterColumn{
					Name:   afterColumn.Name,
					Action: &AlterColumnDropNotNull{},
				},
			})
		case !beforeColumn.NotNull && afterColumn.NotNull:
			// ALTER TABLE table_name ALTER COLUMN column_name SET NOT NULL;
			ddls.Stmts = append(ddls.Stmts, &AlterTableStmt{
				Comment: simplediff.Diff(before.String(), after.String()).String(),
				Name:    after.Name,
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
			Comment: simplediff.Diff(before.String(), after.String()).String(),
			Name:    after.Name,
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
