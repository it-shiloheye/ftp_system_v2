CREATE DATABASE ftp_system_db;
CREATE USER ftp_system_server WITH ENCRYPTED PASSWORD 'yourpass';
GRANT ALL PRIVILEGES ON DATABASE yourdbname TO youruser;