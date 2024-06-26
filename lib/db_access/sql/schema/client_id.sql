

CREATE TYPE file_status_type AS ENUM ('new','deleted','updated');
CREATE TYPE peer_role_type AS ENUM (
    'client', -- push changes to database
    'storage', -- create a local copy of all files
    'server'    -- c
);

CREATE TABLE IF NOT EXISTS peers_table (
    id serial PRIMARY KEY,
    peer_id uuid DEFAULT gen_random_uuid() UNIQUE,
    ip_address TEXT NOT NULL UNIQUE,
    peer_role peer_role_type,
    PEM BYTEA
);


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

CREATE INDEX IF NOT EXISTS file_hash_index
ON file_storage(peer_id, file_hash);

CREATE UNIQUE INDEX IF NOT EXISTS file_name_index
ON file_storage(peer_id, file_path, file_type);