package sqlz //nolint:testpackage

import (
	"context"
	"database/sql"
	"testing"
)

func TestMustBeginTx(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		db := &sqlDBMock{
			BeginTxFunc: func(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
				return &sql.Tx{}, nil
			},
		}
		tx := MustBeginTx(context.Background(), db, &sql.TxOptions{})
		if tx == nil {
			t.Fatalf("❌: MustBeginTx: tx == nil")
		}
	})
	t.Run("failure", func(t *testing.T) {
		t.Parallel()
		db := &sqlDBMock{
			BeginTxFunc: func(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
				return nil, sql.ErrConnDone
			},
		}
		defer func() {
			if err := recover(); err == nil {
				t.Errorf("❌: recover: err == nil")
			}
		}()
		_ = MustBeginTx(context.Background(), db, &sql.TxOptions{})
	})
}
