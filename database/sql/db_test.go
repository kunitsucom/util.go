package sqlz //nolint:testpackage

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"strings"
	"testing"
)

func TestMustOpen(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		sql.Register(t.Name(), &driverDriverMock{
			OpenFunc: func(name string) (driver.Conn, error) {
				return &driverConnMock{
					PrepareFunc: func(query string) (driver.Stmt, error) {
						return &driverStmtMock{
							CloseFunc: func() error {
								return nil
							},
							ExecFunc: func(args []driver.Value) (driver.Result, error) {
								return &driverResultMock{}, nil
							},
							QueryFunc: func(args []driver.Value) (driver.Rows, error) {
								return &driverRowsMock{}, nil
							},
						}, nil
					},
				}, nil
			},
		})

		ctx := context.Background()
		db := MustOpen(ctx, t.Name(), ":memory:")
		if db == nil {
			t.Fatalf("❌: MustOpen: db == nil")
		}
	})

	t.Run("failure,sqlUnknownDriver", func(t *testing.T) {
		t.Parallel()

		defer func() {
			const expect = "sql: unknown driver"
			if actual := fmt.Sprintf("%v", recover()); !strings.Contains(actual, expect) {
				t.Errorf("❌: recover: expect(%v) != actual(%s)", expect, actual)
			}
		}()

		MustOpen(context.Background(), t.Name(), "")
	})

	t.Run("failure,sqlDriverOpenError", func(t *testing.T) {
		t.Parallel()

		defer func() {
			expect := context.Canceled
			if actual := fmt.Sprintf("%v", recover()); fmt.Sprint(expect) != fmt.Sprint(actual) {
				t.Errorf("❌: recover: expect(%v) != actual(%s)", expect, actual)
			}
		}()

		sql.Register(t.Name(), &driverDriverMock{})

		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		MustOpen(ctx, t.Name(), "")
	})
}
