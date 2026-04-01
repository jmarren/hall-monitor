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


-- name: GetPostById :one
SELECT *
FROM posts
WHERE id = @id;


-- name: DeletePostById :exec
DELETE FROM posts
WHERE id = @id;


-- name: GetMostRecentUserPostById :one
SELECT id
FROM posts
WHERE user_id = @user_id
ORDER BY created_at DESC
LIMIT 1;


-- name: GetMostRecentUserPostByUserName :one
SELECT posts.id
FROM posts
JOIN users
	ON users.id = posts.user_id
WHERE users.name = @user_name
ORDER BY posts.created_at DESC
LIMIT 1;

-- name: GetPostAuthor :one
SELECT users.id
FROM users
JOIN posts
	ON posts.user_id = users.id
WHERE posts.id = @post_id;


-- name: GetPostsByUserId :many
SELECT *
FROM posts
WHERE posts.user_id = @user_id;


-- name: DeletePostsByUserId :exec
DELETE FROM posts
WHERE posts.user_id = @user_id;

