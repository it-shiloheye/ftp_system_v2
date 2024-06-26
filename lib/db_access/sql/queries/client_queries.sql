-- name: GetFiles :many
SELECT 
    file_name,
    file_path,
    file_type,
    modification_date,
    file_state,
    file_hash,
    prev_file_hash
 FROM file_storage 
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
    $1, $2, $3, $4, $5, $6, $7
) RETURNING 
    id,
    file_hash
;


-- name: GetFileData :many
SELECT 
    peer_id,
    file_state,
    file_data
FROM file_storage 
WHERE file_hash = $1
LIMIT 1;

-- name: ConnectClient :many
SELECT * FROM peers_table
WHERE ip_address = $1
LIMIT 1;

-- name: CreateClient :many
INSERT INTO peers_table(peer_id,ip_address, PEM )
VALUES  (default,$1, $2)
RETURNING (peer_id);


-- name: GetAllPem :many
SELECT (pem) FROM peers_table;