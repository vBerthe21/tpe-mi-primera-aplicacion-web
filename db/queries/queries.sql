-- name: CreateMovie :one
INSERT INTO movie(name, author, duration, score, review, description)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, name, author, duration, score, review, description;

-- name: GetMovie :one
SELECT id, name, author, duration, score, review, description
FROM movie
WHERE id = $1;

-- name: ListMovie :many
SELECT id, name, author, duration, score, review, description
FROM movie m
ORDER BY name;

-- name: UpdateMovie :exec
UPDATE movie
SET name = $2, author = $3, duration = $4, score = $5, review = $6, description = $7
WHERE id = $1;

-- name: DeleteMovie :exec
DELETE FROM movie
WHERE id = $1;

-- name: AddMovieGenre :one
INSERT INTO movie_genre(movie_id,genre_id)
VALUES( $1, $2)
RETURNING movie_id,genre_id;

-- name: RemoveMovieGenre :exec
DELETE FROM movie_genre
WHERE movie_id = $1 AND genre_id = $2;

-- name: CreateGenre :one
INSERT INTO genre(name_genre)
VALUES ($1)
RETURNING id, name_genre;

-- name: GetGenre :one
SELECT id, name_genre
FROM genre
WHERE id = $1;

-- name: ListGenre :many
SELECT id, name_genre
FROM genre
ORDER BY name_genre;

-- name: UpdateGenre :exec
UPDATE genre
SET name_genre = $2
WHERE id = $1;

-- name: DeleteGenre :exec
DELETE FROM genre
WHERE id = $1;