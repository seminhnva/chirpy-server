-- name: CreateUser :one
INSERT INTO users (
    id,
    created_at,
    updated_at,
    email,
    password
)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: DeleteUsers :exec
DELETE FROM users;

-- name: GetUserInfoByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: UpdateUser :execrows
UPDATE users 
SET email = $1, password = $2
WHERE id = $3;

-- name: UpdateChirpRedById :execrows
UPDATE users
SET is_chirpy_red = true
WHERE id = $1;