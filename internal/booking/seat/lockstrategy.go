package seat

import (
	"context"

	"github.com/soumya-codes/airline-reservation-poc/internal/store"
)

type Seat struct {
	ID     int32
	SeatID string
}

type LockStrategy func(ctx context.Context, q *store.Queries, tripID int32) (*Seat, error)

func GetSeatWithNoLock(ctx context.Context, q *store.Queries, tripID int32) (*Seat, error) {
	s, err := q.GetSeatWithNoLock(ctx, tripID)
	if err != nil {
		return nil, err
	}

	return &Seat{
		ID:     s.Identifier,
		SeatID: s.SeatID,
	}, nil
}

func GetSeatWithSharedLock(ctx context.Context, q *store.Queries, tripID int32) (*Seat, error) {
	s, err := q.GetSeatWithSharedLock(ctx, tripID)
	if err != nil {
		return nil, err
	}

	return &Seat{
		ID:     s.Identifier,
		SeatID: s.SeatID,
	}, nil
}

func GetSeatWithSharedLockSkipped(ctx context.Context, q *store.Queries, tripID int32) (*Seat, error) {
	s, err := q.GetSeatWithSharedLockSkipped(ctx, tripID)
	if err != nil {
		return nil, err
	}

	return &Seat{
		ID:     s.Identifier,
		SeatID: s.SeatID,
	}, nil
}

func GetSeatWithExclusiveLock(ctx context.Context, q *store.Queries, tripID int32) (*Seat, error) {
	s, err := q.GetSeatWithExclusiveLock(ctx, tripID)
	if err != nil {
		return nil, err
	}

	return &Seat{
		ID:     s.Identifier,
		SeatID: s.SeatID,
	}, nil
}

func GetSeatWithExclusiveLockSkipped(ctx context.Context, q *store.Queries, tripID int32) (*Seat, error) {
	s, err := q.GetSeatWithExclusiveLockSkipped(ctx, tripID)
	if err != nil {
		return nil, err
	}

	return &Seat{
		ID:     s.Identifier,
		SeatID: s.SeatID,
	}, nil
}
