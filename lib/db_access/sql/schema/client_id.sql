

CREATE TYPE file_status_type AS ENUM ('new','deleted','updated','stored');

-- peers can only set/update the status of their owned filehash
-- other peers can respond to file_status of other clients 
-- eg requested, to-delete (pull )
CREATE TYPE file_tracker_status AS ENUM(
    -- file doesn't exist and is being requested
    -- eg. peer b wants file a from peer a (and not in file_storage)
    'requested', 
    -- file has been uploaded
    -- eg. peer a uploads file a to file_storage
    'uploaded',
    -- file has been stored
    -- eg. storage peer c has a hard disk copy of file a
    'stored', 
    -- file marked as deleted 
    -- eg. storage peer c has a hard disk copy of file a so peer a deletes it from disk
    'deleted'
);
CREATE TYPE peer_role_type AS ENUM (
    'client', -- push changes to database
    'storage', -- create a local copy of all files
    'server'    -- respend to rpc requests
);

CREATE TABLE IF NOT EXISTS peers_table (
    id serial PRIMARY KEY,
    peer_id uuid DEFAULT gen_random_uuid() UNIQUE,
    ip_address TEXT NOT NULL,
    peer_role peer_role_type,
    peer_name TEXT UNIQUE,
    PEM BYTEA
);



CREATE UNIQUE INDEX peer_index 
ON peers_table(peer_id, ip_address);


CREATE TABLE IF NOT EXISTS file_storage (
    id serial PRIMARY KEY,
    peer_id uuid references peers_table(peer_id) NOT NULL,
    file_name TEXT NOT NULL,
    file_path TEXT NOT NULL,
    file_type VARCHAR(7) NOT NULL,
    file_hash VARCHAR(256) GENERATED ALWAYS AS (encode(sha256(file_data::bytea), 'hex')) STORED UNIQUE,
    prev_file_hash VARCHAR(256) references file_storage(file_hash),
    creation timestamptz default NOW(),
    modification_date TIMESTAMP NOT NULL,
    file_state file_status_type,
    file_data BYTEA
);


CREATE UNIQUE INDEX IF NOT EXISTS file_hash_index
ON file_storage(peer_id, file_hash);

CREATE UNIQUE INDEX IF NOT EXISTS file_name_index
ON file_storage(peer_id, file_path);


CREATE TABLE IF NOT EXISTS file_log ( 
    id serial PRIMARY KEY,
    peer_id uuid references peers_table(peer_id) NOT NULL,
    file_hash varchar(256) references file_storage(file_hash) NOT NULL,
    current_file_status file_status_type,
    old_file_status file_status_type,
    delta_time timestamptz default NOW()
);