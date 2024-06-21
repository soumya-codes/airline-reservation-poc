package config

import (
	"context"
	"time"

	"github.com/soumya-codes/airline-reservation-poc/internal/booking/seat"
	pgconn "github.com/soumya-codes/airline-reservation-poc/internal/postgres/connection"
	pgtx "github.com/soumya-codes/airline-reservation-poc/internal/postgres/transaction"
	"github.com/soumya-codes/airline-reservation-poc/internal/store"
)

const (
	defaultMaxConn     = 10
	defaultTimeout     = 6 * time.Second
	defaultTxIsolation = pgtx.ReadCommitted
)

type Config struct {
	PostgresConfig *pgconn.Config
	MaxConn        int
	Timeout        time.Duration
	LockStrategy   func(ctx context.Context, q *store.Queries, tripID int32) (*seat.Seat, error)
	TxIsolation    pgtx.IsolationLevel
}

func DefaultConfig() *Config {
	return &Config{
		PostgresConfig: &pgconn.Config{
			Host:     "localhost",
			Port:     5432,
			Username: "postgres",
			Password: "postgres",
			Database: "airline_reservation_db",
		},
		MaxConn:      defaultMaxConn,
		Timeout:      defaultTimeout,
		LockStrategy: seat.GetSeatWithExclusiveLock,
		TxIsolation:  defaultTxIsolation,
	}
}

type Option func(*Config)

func WithMaxConn(maxConn int) Option {
	return func(c *Config) {
		c.MaxConn = maxConn
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.Timeout = timeout
	}
}

func WithLockStrategy(strategy func(ctx context.Context, q *store.Queries, tripID int32) (*seat.Seat, error)) Option {
	return func(c *Config) {
		c.LockStrategy = strategy
	}
}

func WithTxIsolation(isolation pgtx.IsolationLevel) Option {
	return func(c *Config) {
		c.TxIsolation = isolation
	}
}

func NewConfig(opts ...Option) *Config {
	cfg := DefaultConfig()
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}
