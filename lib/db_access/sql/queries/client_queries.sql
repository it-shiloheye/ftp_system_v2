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
WHERE peer_id = $1
ORDER BY modification_date DESC;

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
WHERE (
     ip_address = $1
     AND 
     peer_id = $2
)
LIMIT 1;

-- name: CreateClient :many
INSERT INTO peers_table(peer_id,ip_address, PEM,peer_role)
VALUES  (default,$1, $2,$3)
RETURNING *;


-- name: GetAllPem :many
SELECT pem FROM peers_table;

-- name: UpdatePeerRole :exec
UPDATE peers_table
SET 
    peer_role = $2
WHERE 
    peer_id = $1;



-- name: GetCountOfStoragePeers :one
SELECT COUNT(*) FROM peers_table 
WHERE peer_role = 'storage';

-- name: CountIfStored :one
SELECT COUNT(*) FROM file_storage
WHERE 
    (
        file_hash = $1
        OR 
        prev_file_hash = $1
    )
AND
    file_state = 'stored';

-- name: MarkFileDeleted :exec
UPDATE file_storage
SET 
    file_state = 'deleted'
WHERE
    peer_id = $1
    AND 
    file_hash = $2
    AND 
    file_state != 'deleted';


-- name: DownloadFiles :many
SELECT
    file_hash,
    file_name,
    file_data,
    modification_date
FROM file_storage
WHERE 
    peer_id = $1
    AND
    file_state != 'deleted'
    AND
    file_data IS NOT NULL;

-- name: UpdateFileLog :exec
INSERT  INTO file_log (
    peer_id,
    file_hash,
    current_file_status,
    old_file_status 
) VALUES 
    ($1, $2, $3, $4);

