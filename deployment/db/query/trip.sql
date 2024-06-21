-- name: MarkTripForBooking :one
UPDATE trip SET completed = TRUE WHERE id = $1 RETURNING 1;

-- name: GetNextAvailableTrip :one
SELECT id FROM trip WHERE completed = FALSE ORDER BY schedule LIMIT 1 FOR UPDATE SKIP LOCKED;