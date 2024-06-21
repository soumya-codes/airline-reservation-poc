package postgrestransaction

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type IsolationLevel string

const (
	//ReadUncommitted = "READ UNCOMMITTED"
	ReadCommitted  = "READ COMMITTED"
	RepeatableRead = "REPEATABLE READ"
	Serializable   = "SERIALIZABLE"
)

func BeginTxWithIsolationLevel(ctx context.Context, conn *pgx.Conn, isolationLevel IsolationLevel) (pgx.Tx, error) {
	tx, err := conn.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %w", err)
	}

	if _, err := tx.Exec(ctx, fmt.Sprintf("SET TRANSACTION ISOLATION LEVEL %s", isolationLevel)); err != nil {
		rErr := tx.Rollback(ctx)
		return nil, fmt.Errorf("error setting isolation level: %w", rErr)
	}

	return tx, nil
}
