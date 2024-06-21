package main

import (
	"context"
	"log"
	"runtime"

	"github.com/soumya-codes/airline-reservation-poc/config"
	"github.com/soumya-codes/airline-reservation-poc/internal/booking"
	//"github.com/soumya-codes/airline-reservation-poc/internal/booking/seat"
	//pgtx "github.com/soumya-codes/airline-reservation-poc/internal/postgres/transaction"
)

func main() {
	// Set the maximum number of CPUs that can be executing simultaneously.
	runtime.GOMAXPROCS(4)

	/*
		cfg := config.NewConfig(
			config.WithMaxConn(180),
			config.WithLockStrategy(seat.GetSeatWithNoLock),
			config.WithTxIsolation(pgtx.ReadCommitted),
		)
	*/

	cfg := config.DefaultConfig()
	ctx, cancelFunc := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancelFunc()

	err := booking.BookSeats(ctx, cfg)
	if err != nil {
		log.Fatalf("Error running booking process: %v", err)
	}
}
