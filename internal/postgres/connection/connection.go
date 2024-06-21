package postgresconnection

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/soumya-codes/airline-reservation-poc/internal/postgres/transaction"
)

// Config represents the configuration for connecting to the Postgres database.
type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
}

// NewConnection returns a new connection instance to connect to the Postgres database.
func NewConnection(config *Config) (*pgx.Conn, error) {
	connConfig, err := pgx.ParseConfig(fmt.Sprintf("postgres://%s:%s@%s:%d/%s", config.Username, config.Password, config.Host, config.Port, config.Database))
	if err != nil {
		return nil, err
	}

	conn, err := pgx.ConnectConfig(context.Background(), connConfig)
	if err != nil {

		return nil, err
	}

	return conn, nil
}

// Close closes the connection to the Postgres database.
func Close(conn *pgx.Conn) error {
	return conn.Close(context.Background())
}

func BeginTx(ctx context.Context, conn *pgx.Conn, iLevel postgrestransaction.IsolationLevel) (pgx.Tx, error) {
	txOptions := pgx.TxOptions{}

	switch iLevel {
	case postgrestransaction.ReadCommitted:
		txOptions.IsoLevel = pgx.ReadCommitted
	case postgrestransaction.RepeatableRead:
		txOptions.IsoLevel = pgx.RepeatableRead
	case postgrestransaction.Serializable:
		txOptions.IsoLevel = pgx.Serializable
	// No need to test ReadUncommitted as ReadUncommitted is same as ReadCommitted in postgres.
	//case transaction.ReadUncommitted:
	//	txOptions.IsoLevel = pgx.ReadUncommitted
	default:
		return nil, fmt.Errorf("unknown isolation level: %v", iLevel)
	}

	return conn.BeginTx(ctx, txOptions)
}
