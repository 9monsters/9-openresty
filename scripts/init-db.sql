-- Database setup for OpenResty Admin hub.
-- Run as postgres superuser:
--   psql -h 192.168.5.100 -U postgres -f scripts/init-db.sql
DO $$
BEGIN
  IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'openresty') THEN
    CREATE ROLE openresty LOGIN PASSWORD '1aac@Openresty';
  END IF;
END
$$;

SELECT 'CREATE DATABASE openresty OWNER openresty'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'openresty')\gexec

GRANT ALL PRIVILEGES ON DATABASE openresty TO openresty;

\c openresty
GRANT ALL ON SCHEMA public TO openresty;