package postgresconnectionpool

import (
	"sync"

	"github.com/jackc/pgx/v5"
	pgconn "github.com/soumya-codes/airline-reservation-poc/internal/postgres/connection"
)

type ConnectionPool struct {
	mu      *sync.Mutex
	conns   []*pgx.Conn
	maxConn int
	channel chan interface{}
}

func NewConnectionPool(config *pgconn.Config, maxConn int) (*ConnectionPool, error) {
	var mu = sync.Mutex{}
	cPool := &ConnectionPool{
		mu:      &mu,
		conns:   make([]*pgx.Conn, 0, maxConn),
		maxConn: maxConn,
		channel: make(chan interface{}, maxConn),
	}

	for i := 0; i < maxConn; i++ {
		conn, err := pgconn.NewConnection(config)
		if err != nil {
			return nil, err
		}

		cPool.conns = append(cPool.conns, conn)
		cPool.channel <- struct{}{}
	}

	return cPool, nil
}

func (cPool *ConnectionPool) NumberOfConnections() int {
	return len(cPool.conns)
}

func (cpool *ConnectionPool) Acquire() *pgx.Conn {
	<-cpool.channel

	cpool.mu.Lock()
	c := cpool.conns[0]
	cpool.conns = cpool.conns[1:]
	cpool.mu.Unlock()

	return c
}

func (cPool *ConnectionPool) Release(c *pgx.Conn) {
	cPool.mu.Lock()
	cPool.conns = append(cPool.conns, c)
	cPool.channel <- struct{}{}
	cPool.mu.Unlock()
}

func (cPool *ConnectionPool) Close() error {
	for _, conn := range cPool.conns {
		err := pgconn.Close(conn)
		if err != nil {
			return err
		}
	}

	return nil
}
