-- name: GetUserById :one
SELECT *
FROM users
WHERE id = @id;


-- name: GetUserByName :one
SELECT *
FROM users
WHERE name = @name;


-- name: DeleteUserById :exec
DELETE FROM users
WHERE id = @id;


-- name: DeleteUserByName :exec
DELETE FROM users
WHERE name = @name;
