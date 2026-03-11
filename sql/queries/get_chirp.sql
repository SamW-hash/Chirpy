-- name: GetChirp :one
SELECT *
FROM chirps
WHERE id = $1;

-- name: GetChirpsByUser :many
SELECT *
FROM chirps
WHERE user_id = $1;