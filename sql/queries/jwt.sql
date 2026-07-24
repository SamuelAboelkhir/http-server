-- name: SaveRefreshToken :exec
INSERT INTO refresh_tokens(
created_at, updated_at, token, user_id, expires_at
)
VALUES(
  NOW(),
  NOW(),
  $1, 
  $2, 
  $3);

-- name: GetRefreshToken :one
SELECT * FROM refresh_tokens
WHERE token = $1;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET revoked_at = $1, updated_at = $2
WHERE token = $3;
