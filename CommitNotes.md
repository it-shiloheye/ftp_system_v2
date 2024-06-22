# Commit Notes

### 22nd June 2024, 23:59 GMT +3
```sh
1. Shifting from server-client to peer-db-peer 
2. Using React-SWC-TS as browser frontend
3. Finalising connection to db
# Please enter the commit message for your changes. Lines starting
# with '#' will be ignored, and an empty message aborts the commit.
#
# On branch main
# Your branch is up to date with 'origin/main'.
#
# Changes to be committed:
#	modified:   .gitignore
#	deleted:    _lib/.gitignore
#	deleted:    _lib/db_access/db_main.go
#	deleted:    _lib/db_access/sql/init.sql
#	deleted:    _lib/exports.go
#	new file:   client/.eslintrc.cjs
#	new file:   client/.gitignore
#	new file:   client/README.md
#	new file:   client/index.html
#	deleted:    client/main.go
#	new file:   client/public/vite.svg
#	new file:   client/src/App.css
#	new file:   client/src/App.tsx
#	new file:   client/src/assets/react.svg
#	new file:   client/src/index.css
#	new file:   client/src/main.tsx
#	new file:   client/src/vite-env.d.ts
#	new file:   client/vite.config.ts
#	modified:   dev_kill_script.ps1
#	modified:   dev_script.ps1
#	modified:   go.mod
#	modified:   go.sum
#	new file:   init.example.sql
#	renamed:    _lib/.air.toml -> lib/.air.toml
#	new file:   lib/.gitignore
#	renamed:    _lib/base/atomic_json.go -> lib/base/atomic_json.go
#	renamed:    _lib/base/base.go -> lib/base/base.go
#	renamed:    _lib/base/init.go -> lib/base/init.go
#	new file:   lib/base/ip_handling.go
#	renamed:    _lib/base/mutexed_map.go -> lib/base/mutexed_map.go
#	renamed:    _lib/base/mutexed_queue.go -> lib/base/mutexed_queue.go
#	renamed:    _lib/context/context.go -> lib/context/context.go
#	new file:   lib/db_access/db_main.go
#	new file:   lib/db_access/generated/client_queries.sql.go
#	new file:   lib/db_access/generated/db.go
#	new file:   lib/db_access/generated/models.go
#	new file:   lib/db_access/sql/queries/client_queries.sql
#	new file:   lib/db_access/sql/schema/client_id.sql
#	renamed:    _lib/file_handler/v2/bytes_store.go -> lib/file_handler/v2/bytes_store.go
#	renamed:    _lib/file_handler/v2/file_basic.go -> lib/file_handler/v2/file_basic.go
#	renamed:    _lib/file_handler/v2/file_hash.go -> lib/file_handler/v2/file_hash.go
#	renamed:    _lib/file_handler/v2/lock_file.go -> lib/file_handler/v2/lock_file.go
#	renamed:    _lib/logging/fake_logger.go -> lib/logging/fake_logger.go
#	renamed:    _lib/logging/log_item/error_type.go -> lib/logging/log_item/error_type.go
#	renamed:    _lib/logging/logging_struct.go -> lib/logging/logging_struct.go
#	renamed:    client/network_client/network_client.go -> lib/network_client/network_client.go
#	renamed:    client/network_client/network_engine.go -> lib/network_client/network_engine.go
#	new file:   lib/sqlc.yaml
#	renamed:    _lib/tls_handler/v2/cert_data.go -> lib/tls_handler/v2/cert_data.go
#	renamed:    _lib/tls_handler/v2/cert_handler_2.go -> lib/tls_handler/v2/cert_handler_2.go
#	modified:   package.json
#	renamed:    server/.air.toml -> server-peer/.air.toml
#	renamed:    server/config/config.go -> server-peer/config/config.go
#	new file:   server-peer/main.go
#	renamed:    client/remove-item.ps1 -> server-peer/remove-item.ps1
#	renamed:    server/server/init_server.go -> server-peer/server/init_server.go
#	renamed:    server/server/register_routes.go -> server-peer/server/register_routes.go
#	renamed:    server/server/server_loop.go -> server-peer/server/server_loop.go
#	deleted:    server/main.go
#	deleted:    server/remove-item.ps1
#
```