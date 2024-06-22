-- name: GetFiles :many
SELECT (
    file_name,
    file_path,
    file_type,
    modification_date,
    file_state,
    file_hash,
    prev_file_hash
) FROM file_storage 
WHERE peer_id = $1;

-- name: InsertFile :one
INSERT INTO file_storage (
    peer_id,
    file_name,
    file_path,
    file_type,
    modification_date,
    file_state,
    file_data
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7
) RETURNING (
    id,
    file_hash
);


-- name: GetFileData :one
SELECT (
    peer_id,
    file_state,
    file_data
) FROM file_storage 
WHERE file_hash = $1;