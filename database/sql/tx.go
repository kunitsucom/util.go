package sqlz

import (
	"context"
	"database/sql"
)

func MustBeginTx(ctx context.Context, db sqlTxBeginner, opts *sql.TxOptions) *sql.Tx {
	tx, err := db.BeginTx(ctx, opts)
	if err != nil {
		panic(err)
	}
	return tx
}
