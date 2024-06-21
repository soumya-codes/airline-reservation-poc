-- name: GetSeatWithNoLock :one
SELECT id, seat_id FROM reservation WHERE trip_id = $1 AND passenger_id IS NULL ORDER BY id LIMIT 1;

-- name: GetSeatWithSharedLock :one
SELECT id, seat_id FROM reservation WHERE trip_id = $1 AND passenger_id IS NULL ORDER BY id LIMIT 1 FOR SHARE;

-- name: GetSeatWithSharedLockSkipped :one
SELECT id, seat_id FROM reservation WHERE trip_id = $1 AND passenger_id IS NULL ORDER BY id LIMIT 1 FOR SHARE SKIP LOCKED;

-- name: GetSeatWithExclusiveLock :one
SELECT id, seat_id FROM reservation WHERE trip_id = $1 AND passenger_id IS NULL ORDER BY id LIMIT 1 FOR UPDATE;

-- name: GetSeatWithExclusiveLockSkipped :one
SELECT id, seat_id FROM reservation WHERE trip_id = $1 AND passenger_id IS NULL ORDER BY id LIMIT 1 FOR UPDATE SKIP LOCKED;

-- name: BookSeat :one
UPDATE reservation SET passenger_id = $1 WHERE id = $2 RETURNING 1;

-- name: GetTripSeats :many
SELECT seat_id, passenger_id FROM reservation WHERE trip_id = $1 ORDER BY seat_id;