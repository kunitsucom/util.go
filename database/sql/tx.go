package sqlz

import (
	"context"
	"database/sql"
	"fmt"
)

func MustBeginTx(ctx context.Context, db sqlTxBeginner, opts *sql.TxOptions) *sql.Tx {
	tx, err := db.BeginTx(ctx, opts)
	if err != nil {
		err = fmt.Errorf("BeginTx: %w", err)
		panic(err)
	}
	return tx
}
