package postgresconnection

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
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
