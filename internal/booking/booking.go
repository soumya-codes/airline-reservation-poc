package booking

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	pgconn2 "github.com/jackc/pgx/v5/pgconn"
	"github.com/sirupsen/logrus"
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

type booking struct {
	passengerName string
	seatNumber    int32
	seatId        string
}

type bookingStatus struct {
	booking
	err error
}

func BookSeats(ctx context.Context, config *config.Config) error {
	// Initiate connection to the database
	pgConfig := config.PostgresConfig
	conn, err := pgconn.NewConnection(pgConfig)
	if err != nil {
		return fmt.Errorf("error connecting to database to start the bookings: %w", err)
	}

	defer func() {
		if err := pgconn.Close(conn); err != nil {
			logrus.WithError(err).Error("error closing database connection")
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
		return fmt.Errorf("error getting next available tripID: %s", err.Error())
	}

	// Mark the tripID, so it's not considered for booking again
	err = MarkTripForBooking(ctx, q, tripID)
	if err != nil {
		return fmt.Errorf("error marking tripID as booked: %w", err)
	}

	// Create a connection pool of size maxConn
	pool, err := pgpool.NewConnectionPool(pgConfig, config.MaxConn)
	if err != nil {
		return fmt.Errorf("error creating connection pool: %s", err.Error())
	}

	start := time.Now()

	bks := make(chan bookingStatus, len(passengers))
	var wg sync.WaitGroup

	for _, passenger := range passengers {
		wg.Add(1)
		// Book a seat for the passenger
		passenger := passenger // Not necessary for Golang versions >= 1.22
		go func() {
			defer wg.Done()
			bookSeatTask(
				ctx,
				tripID,
				passenger,
				pool,
				config.LockStrategy,
				config.TxIsolation,
				bks,
				config.MaxRetries)
		}()
	}

	// Close the channel once all the goroutines are done
	go func() {
		wg.Wait()
		close(bks)
	}()

	// Get the results from the channels and store them in the bookings and reservation slice
	bookings := make([]string, 0)
	reservations := make([]booking, int32(len(passengers)))
	// Loop through the bookings channel to get booking info/error, the loop ends when bks channel is closed
	for bk := range bks {
		if bk.err != nil {
			bookings = append(bookings, fmt.Sprintf("ERROR: couldn't book seat: %s", bk.err.Error()))
		} else {
			bookings = append(bookings, fmt.Sprintf("Seat: %s is booked for passenger: %s", bk.seatId, bk.passengerName))
			reservations[(bk.seatNumber-1)%int32(len(passengers))] = bk.booking
		}
	}

	elapsed := time.Since(start)
	defer printBookingAndReservationDetails(reservations, bookings, tripID, elapsed)

	return nil
}

// GetPassengers retrieves the list of passengers from the database.
func GetPassengers(ctx context.Context, q *store.Queries) ([]store.Passenger, error) {
	// Get the list of passengers
	passengers, err := q.GetPassengers(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting passengers: %w", err)
	}

	return passengers, nil
}

// GetNextAvailableTrip retrieves the next available trip ID from the database.
func GetNextAvailableTrip(ctx context.Context, q *store.Queries) (int32, error) {
	// Get the next trip on schedule
	tripID, err := q.GetNextAvailableTrip(ctx)
	if err != nil {
		return -1, fmt.Errorf("error getting the trip: %d", err)
	}

	return tripID, nil
}

// MarkTripForBooking marks a trip as booked in the database.
func MarkTripForBooking(ctx context.Context, q *store.Queries, tripID int32) error {
	// Mark the trip as booked
	_, err := q.MarkTripForBooking(ctx, tripID)
	if err != nil {
		return fmt.Errorf("error marking trip for booking: %w", err)
	}

	return nil
}

// bookSeatTask handles the booking of a seat for a passenger.
func bookSeatTask(ctx context.Context,
	tripID int32,
	passenger store.Passenger,
	pool *pgpool.ConnectionPool,
	seatLockStrategy bookingseat.LockStrategy,
	isolationLevel pgtx.IsolationLevel,
	bs chan<- bookingStatus,
	maxRetries int,
) {
	// Acquire a connection from the pool
	conn := pool.Acquire()
	defer pool.Release(conn)

	for retry := 1; retry <= maxRetries; retry++ {
		// Start a transaction
		tx, err := pgtx.BeginTxWithIsolationLevel(ctx, conn, isolationLevel)
		if err != nil {
			retriesLeft := handleRetries(
				retry,
				maxRetries,
				booking{passengerName: passenger.Name},
				fmt.Errorf("retry %d/%d failed: error starting transaction: %w",
					retry,
					maxRetries,
					err,
				),
				bs,
			)
			if retriesLeft {
				continue
			}

			return
		}

		// Get queries instance to execute requests in the transaction
		q := store.New(conn).WithTx(tx)

		// Get the next available seat
		seat, err := seatLockStrategy(ctx, q, tripID)
		if err != nil {
			txErrMsg := fmt.Sprintf("retry %d/%d failed: error getting next available seat for passenger %s",
				retry,
				maxRetries,
				passenger.Name,
			)

			err = handleTransactionError(ctx, tx, txErrMsg, err)
			retriesLeft := handleRetries(
				retry,
				maxRetries,
				booking{passengerName: passenger.Name},
				err,
				bs)
			if retriesLeft {
				continue
			}

			return
		}

		// Book a seat for the passenger
		_, err = q.BookSeat(ctx, store.BookSeatParams{PassengerID: passenger.Identifier, Identifier: seat.ID})
		if err != nil {
			txErrMsg := fmt.Sprintf("retry %d/%d failed: error booking seat %s for passenger %s",
				retry,
				maxRetries,
				seat.SeatID,
				passenger.Name,
			)

			err = handleTransactionError(ctx, tx, txErrMsg, err)
			retriesLeft := handleRetries(
				retry,
				maxRetries,
				booking{passengerName: passenger.Name, seatId: seat.SeatID, seatNumber: seat.ID},
				err,
				bs)
			if retriesLeft {
				continue
			}

			return
		}

		err = tx.Commit(ctx)
		if err != nil {
			retriesLeft := handleRetries(
				retry,
				maxRetries,
				booking{
					passengerName: passenger.Name,
					seatId:        seat.SeatID,
					seatNumber:    seat.ID,
				},
				fmt.Errorf("retry %d/%d failed: error committing transaction for passenger: %s, %w",
					retry,
					maxRetries,
					passenger.Name,
					err,
				),
				bs,
			)
			if retriesLeft {
				continue
			}

			return
		}

		bs <- bookingStatus{
			err: nil,
			booking: booking{
				passengerName: passenger.Name,
				seatId:        seat.SeatID,
				seatNumber:    seat.ID,
			},
		}

		return
	}
}

// handleTransactionError handles transaction rollback and sends the error to the result channel.
func handleTransactionError(ctx context.Context, tx pgx.Tx, msg string, err error) error {
	errMsg := fmt.Errorf("%s: %w", msg, err)
	if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
		errMsg = fmt.Errorf("%s: %w, error rolling back transaction: %w", msg, err, rollbackErr)
	}

	return errMsg
}

// handleRetries checks if the max retries are exhausted and sends error to the result channel.
func handleRetries(retry int, maxRetries int, bk booking, err error, bs chan<- bookingStatus) bool {
	bs <- bookingStatus{
		err:     err,
		booking: bk,
	}

	if pgErr, ok := err.(*pgconn2.PgError); ok && pgErr.Code == "40P01" {
		// Deadlock detected
		// This value is tightly coupled with query execution time and "deadlock_timeout" value set in postgresql.conf
		// TODO: Implement a better way to handle deadlocks
		// TODO: Try to make this configurable by adding a field deadlock_retry_time_interval in the Config struct
		time.Sleep(90 * time.Millisecond)
	}

	// Do we really need 2 separate sleep times for deadlock and other errors?
	// TODO: Implement a better way to handle deadlocks
	// TODO: Try to make this configurable by adding a field retry_time_interval in the Config struct
	time.Sleep(30 * time.Millisecond)

	if allRetriesExhausted(retry, maxRetries) {
		return false
	}

	return true
}

func allRetriesExhausted(retry int, maxRetries int) bool {
	return maxRetries > 0 && retry >= maxRetries
}

// print booking process(successful and failed tx) details, including the final reservation details.
func printBookingAndReservationDetails(reservations []booking, bookings []string, tripID int32, elapsedTime time.Duration) {
	logrus.Infof("Total time taken to book the seats for trip-id: %d is %v", tripID, elapsedTime)

	fmt.Print("\n\n")

	// Print the booking details, this contains details of the successful, overlapping and failed bookings/transactions.
	printBookingDetails(bookings)

	fmt.Print("\n\n")

	// Print the final seat reservation details.
	printReservationDetails(reservations)

	fmt.Print("\n\n\n\n")
}

func printBookingDetails(bookings []string) {
	logrus.Info("Booking details:")
	for _, booking := range bookings {
		if strings.Contains(booking, "ERROR") {
			logrus.Error(booking)
		} else {
			logrus.Info(booking)
		}
	}
}

func printReservationDetails(reservations []booking) {
	logrus.Info("Final seat reservation details:")
	for _, reservation := range reservations {
		if reservation.passengerName != "" {
			logrus.Infof("Seat: %s, is assigned to Passenger: %s", reservation.seatId, reservation.passengerName)
		}
	}

	for col := 0; col < seatsPerRow; col++ {
		for row := 0; row < totalRows; row++ {
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
}
