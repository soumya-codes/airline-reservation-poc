package main

import (
	"context"
	"log"
	"runtime"
	"time"

	"github.com/soumya-codes/airline-reservation-poc/config"
	"github.com/soumya-codes/airline-reservation-poc/internal/booking"
	"github.com/soumya-codes/airline-reservation-poc/internal/booking/seat"
	pgtx "github.com/soumya-codes/airline-reservation-poc/internal/postgres/transaction"
)

func main() {
	// Set the maximum number of CPUs that can be executing simultaneously.
	runtime.GOMAXPROCS(8)

	//cfg := config.DefaultConfig()

	cfg := config.NewConfig(
		config.WithTxIsolation(pgtx.ReadCommitted),
		config.WithLockStrategy(seat.GetSeatWithSharedLock),
		config.WithMaxConn(5),
		config.WithMaxRetries(3),
	)

	ctx, cancelFunc := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancelFunc()

	err := booking.BookSeats(ctx, cfg)
	if err != nil {
		log.Fatalf("Error running booking process: %v", err)
	}
}
