-- name: GetBook :one

SELECT * FROM books WHERE id = sqlc.arg(id) LIMIT 1;

-- name: ListBooks :many

SELECT * FROM books ORDER BY published_at;

-- name: CreateBook :execresult

INSERT INTO
    books (title, published_at)
VALUES (
        sqlc.arg(title),
        sqlc.arg(published_at)
    );

-- name: UpdateBook :execresult

UPDATE books
SET
    title = COALESCE(sqlc.narg(title), title),
    version = version + 1
WHERE
    id = sqlc.arg(id)
    AND version = sqlc.arg(version);

-- name: DeleteBook :exec

DELETE FROM books WHERE id = sqlc.arg(id);