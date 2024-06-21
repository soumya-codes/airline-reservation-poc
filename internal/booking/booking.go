package booking

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/soumya-codes/airline-reservation-poc/config"
	bookingseat "github.com/soumya-codes/airline-reservation-poc/internal/booking/seat"
	pgconn "github.com/soumya-codes/airline-reservation-poc/internal/postgres/connection"
	pgpool "github.com/soumya-codes/airline-reservation-poc/internal/postgres/connectionpool"
	pgtx "github.com/soumya-codes/airline-reservation-poc/internal/postgres/transaction"
	"github.com/soumya-codes/airline-reservation-poc/internal/store"
)

const (
	totalRows   = 30
	seatsPerRow = 6
	aisleCol    = 2
)

type reservation struct {
	seatNumber    string
	passengerName string
}

func BookSeats(ctx context.Context, config *config.Config) error {
	// Initiate connection to the database
	pgConfig := config.PostgresConfig
	conn, err := pgconn.NewConnection(pgConfig)
	if err != nil {
		log.Fatal("error connecting to DataBase", err)
	}

	defer func() {
		if err := pgconn.Close(conn); err != nil {
			log.Fatal("error closing DataBase connection", err)
		}
	}()

	q := store.New(conn)

	// Get the list of passengers
	passengers, err := GetPassengers(ctx, q)
	if err != nil {
		return fmt.Errorf("error getting passengers: %w", err)
	}

	// Get the next available tripID
	tripID, err := GetNextAvailableTrip(ctx, q)
	if err != nil {
		return fmt.Errorf("error getting next available tripID: %v", err)
	}

	// Mark the tripID, so its not considered for booking again
	err = MarkTripForBooking(ctx, q, tripID)
	if err != nil {
		return fmt.Errorf("error marking tripID as booked: %w", err)
	}

	// Create a connection pool of size maxConn
	pool, err := pgpool.NewConnectionPool(pgConfig, config.MaxConn)
	if err != nil {
		log.Fatalf("Error creating connection pool: %v", err)
	}

	start := time.Now()
	var g errgroup.Group
	var mu sync.Mutex
	reservations := make([]reservation, 180)
	for _, passenger := range passengers {
		// Book a seat for the passenger
		passenger := passenger // Not necessary for Golang versions >= 1.22
		g.Go(func() error {
			return bookSeatTask(ctx, tripID, passenger, pool, config.LockStrategy, config.TxIsolation, reservations, &mu)
		})
	}

	// Wait for all the goroutines to finish, or an error to occur
	if err := g.Wait(); err != nil {
		return fmt.Errorf("error booking seats: %w", err)
	}

	elapsed := time.Since(start)
	fmt.Printf("Total elapsed time: %s\n", elapsed)
	printBookingDetails(reservations, tripID, elapsed)

	return nil
}

func GetPassengers(ctx context.Context, q *store.Queries) ([]store.Passenger, error) {
	// Get the list of passengers
	passengers, err := q.GetPassengers(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting passengers: %w", err)
	}

	return passengers, nil
}

func GetNextAvailableTrip(ctx context.Context, q *store.Queries) (int32, error) {
	// Get the next trip on schedule
	tripID, err := q.GetNextAvailableTrip(ctx)
	if err != nil {
		return -1, fmt.Errorf("error getting the trip: %d", err)
	}

	return tripID, nil
}

func bookSeatTask(ctx context.Context,
	tripID int32,
	passenger store.Passenger,
	pool *pgpool.ConnectionPool,
	lockStrategy bookingseat.LockStrategy,
	isolationLevel pgtx.IsolationLevel,
	reservations []reservation,
	mu *sync.Mutex) error {
	// Acquire a connection from the pool
	conn := pool.Acquire()
	defer pool.Release(conn)

	// Start a transaction
	tx, err := pgtx.BeginTxWithIsolationLevel(ctx, conn, isolationLevel)
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}

	// Get queries resource to execute queries
	q := store.New(conn).WithTx(tx)

	// Get the next available seat
	seat, err := lockStrategy(ctx, q, tripID)
	if err != nil {
		rollbackErr := tx.Rollback(ctx)
		if rollbackErr != nil {
			return fmt.Errorf("error rolling back transaction: %w, "+
				"error getting next available seat for passenger %s, for trip-id: %d, %w", rollbackErr, passenger.Name, tripID, err)
		}

		return fmt.Errorf("error getting seat for passenger: %s, for trip-id: %d, %w", passenger.Name, tripID, err)
	}

	// Book a seat for the passenger
	_, err = q.BookSeat(ctx, store.BookSeatParams{PassengerID: passenger.Identifier, Identifier: seat.ID})
	if err != nil {
		rollbackErr := tx.Rollback(ctx)
		if rollbackErr != nil {
			return fmt.Errorf("error rolling back transaction: %w, "+
				"error booking seat for passesger %s for the trip-id: %d, %w", rollbackErr, passenger.Name, tripID, err)
		}

		return fmt.Errorf("error booking seat for passenger: %s, for the trip-id: %d, %w", passenger.Name, tripID, err)

	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("error committing transaction for passenger: %s, for the trip-id: %d,: %w", passenger.Name, tripID, err)
	}

	mu.Lock()
	defer mu.Unlock()
	reservations[(seat.ID-1)%180] = reservation{seatNumber: seat.SeatID, passengerName: passenger.Name}
	return nil
}

func MarkTripForBooking(ctx context.Context, q *store.Queries, tripID int32) error {
	// Mark the trip as booked
	_, err := q.MarkTripForBooking(ctx, tripID)
	if err != nil {
		return fmt.Errorf("error marking trip as booked: %w", err)
	}

	return nil
}

func printBookingDetails(reservations []reservation, tripId int32, t time.Duration) {

	for _, reservation := range reservations {
		if reservation.passengerName != "" {
			fmt.Printf("Seat: %s, is assigned to Passenger: %s\n", reservation.seatNumber, reservation.passengerName)
		}
	}

	fmt.Print("\n\n")
	fmt.Printf("Total time taken to book the seats for trip-id: %d is %s\n\n", tripId, t)

	for col := range seatsPerRow {
		for row := range totalRows {
			index := row*seatsPerRow + col
			if reservations[index].passengerName == "" {
				fmt.Print(".")
			} else {
				fmt.Print("x")
			}

			// Print a space after each seat except the last one in the column
			if col == aisleCol && row == totalRows-1 {
				fmt.Print("\n\n") // Print an extra space for the aisle
			} else if row != totalRows-1 {
				fmt.Print("  ")
			}
		}

		// Move to the next column
		fmt.Println()
	}

	fmt.Print("\n\n")
}
