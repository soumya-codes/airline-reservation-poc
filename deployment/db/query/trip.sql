-- name: MarkTripForBooking :one
UPDATE trip SET booked = TRUE WHERE id = $1 RETURNING 1;

-- name: GetNextAvailableTrip :one
SELECT id FROM trip WHERE booked = FALSE ORDER BY schedule LIMIT 1 FOR UPDATE SKIP LOCKED;