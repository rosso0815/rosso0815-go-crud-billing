
-- name: UserkvList :many
SELECT
  *
FROM
    userkv
WHERE
    user_id = $1;

-- name: UserkvDelete :exec
DELETE FROM userkv WHERE userkv_id = @userkvid;

-- name: UserkvCreate :one
INSERT INTO Userkv
	(user_id, key, value)
VALUES
	(@user_id, @key, @value)
RETURNING userkv_id;

-- name: UserkvUpdate :exec
UPDATE userkv
SET
	user_id = @user_id,
	key = @key,
	value = @value
WHERE userkv_id = @userkv_id;

-- EOF
