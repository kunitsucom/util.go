package sqlz //nolint:testpackage

import (
	"context"
	"database/sql"
	"database/sql/driver"
)

type driverDriverMock struct {
	OpenFunc func(name string) (driver.Conn, error)
}

func (m *driverDriverMock) Open(name string) (driver.Conn, error) {
	return m.OpenFunc(name)
}

var _ driver.Driver = (*driverDriverMock)(nil)

type driverConnMock struct {
	PrepareFunc func(query string) (driver.Stmt, error)
	CloseFunc   func() error
	BeginFunc   func() (driver.Tx, error)
}

func (m *driverConnMock) Prepare(query string) (driver.Stmt, error) {
	return m.PrepareFunc(query)
}

func (m *driverConnMock) Close() error {
	return m.CloseFunc()
}

func (m *driverConnMock) Begin() (driver.Tx, error) {
	return m.BeginFunc()
}

var _ driver.Conn = (*driverConnMock)(nil)

type driverStmtMock struct {
	CloseFunc    func() error
	NumInputFunc func() int
	ExecFunc     func(args []driver.Value) (driver.Result, error)
	QueryFunc    func(args []driver.Value) (driver.Rows, error)
}

func (m *driverStmtMock) Close() error {
	return m.CloseFunc()
}

func (m *driverStmtMock) NumInput() int {
	return m.NumInputFunc()
}

func (m *driverStmtMock) Exec(args []driver.Value) (driver.Result, error) {
	return m.ExecFunc(args)
}

func (m *driverStmtMock) Query(args []driver.Value) (driver.Rows, error) {
	return m.QueryFunc(args)
}

var _ driver.Stmt = (*driverStmtMock)(nil)

type driverResultMock struct {
	LastInsertIdFunc func() (int64, error) //nolint:stylecheck // NOTE: sql.Result has LastInsertId method
	RowsAffectedFunc func() (int64, error)
}

func (m *driverResultMock) LastInsertId() (int64, error) {
	return m.LastInsertIdFunc()
}

func (m *driverResultMock) RowsAffected() (int64, error) {
	return m.RowsAffectedFunc()
}

var _ driver.Result = (*driverResultMock)(nil)

type driverRowsMock struct {
	CloseFunc   func() error
	ColumnsFunc func() []string
	NextFunc    func(dest []driver.Value) error
}

func (m *driverRowsMock) Close() error {
	return m.CloseFunc()
}

func (m *driverRowsMock) Columns() []string {
	return m.ColumnsFunc()
}

func (m *driverRowsMock) Next(dest []driver.Value) error {
	return m.NextFunc(dest)
}

var _ driver.Rows = (*driverRowsMock)(nil)

type sqlDBMock struct {
	Rows  *sql.Rows
	Error error

	BeginTxFunc func(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

var (
	_ sqlQueryerContext = (*sqlDBMock)(nil)
	_ sqlTxBeginner     = (*sqlDBMock)(nil)
)

func (m *sqlDBMock) QueryContext(_ context.Context, _ string, _ ...interface{}) (*sql.Rows, error) {
	return m.Rows, m.Error
}

func (m *sqlDBMock) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return m.BeginTxFunc(ctx, opts)
}

type sqlRowsMock struct {
	CloseFunc   func() error
	ColumnsFunc func() ([]string, error)
	NextFunc    func() bool
	ScanFunc    func(dest ...interface{}) error
	ErrFunc     func() error
}

var _ sqlRows = (*sqlRowsMock)(nil)

func (m *sqlRowsMock) Close() error {
	if m.CloseFunc == nil {
		return nil
	}
	return m.CloseFunc()
}

func (m *sqlRowsMock) Columns() ([]string, error) {
	return m.ColumnsFunc()
}

func (m *sqlRowsMock) Next() bool {
	return m.NextFunc()
}

func (m *sqlRowsMock) Scan(dest ...interface{}) error {
	return m.ScanFunc(dest...)
}

func (m *sqlRowsMock) Err() error {
	if m.ErrFunc == nil {
		return nil
	}
	return m.ErrFunc()
}
