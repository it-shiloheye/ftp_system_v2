CREATE TABLE IF NOT EXISTS peers_table (
    id serial PRIMARY KEY,
    peer_id uuid DEFAULT gen_random_uuid(),
    ip_address TEXT NOT NULL
);

CREATE TYPE file_status_type AS ENUM ('new','shared','deleted');

CREATE TABLE IF NOT EXISTS file_storage (
    id serial PRIMARY KEY,
    peer_id uuid references peers_table(peer_id) NOT NULL,
    file_name TEXT NOT NULL,
    file_path TEXT NOT NULL,
    file_type VARCHAR(7) NOT NULL,
    file_hash VARCHAR(256) GENERATED ALWAYS AS (encode(sha256(file_data::bytea), 'hex')) STORED,
    prev_file_hash VARCHAR(256) references file_storage(file_hash),
    creation timestamptz default NOW(),
    modification_date TIMESTAMP NOT NULL,
    file_state file_status_type,
    file_data BINARY
);

CREATE INDEX IF NOT EXISTS file_hash_index
ON file_storage(peer_id, file_hash);

CREATE INDEX IF NOT EXISTS file_name_index
ON file_storage(peer_id, file_path, file_type);