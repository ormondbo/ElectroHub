#!/usr/bin/env bash
set -Eeuo pipefail

: "${POSTGRES_USER:?POSTGRES_USER is required}"
: "${AUTH_POSTGRES_DB:?AUTH_POSTGRES_DB is required}"
: "${AUTH_POSTGRES_USER:?AUTH_POSTGRES_USER is required}"
: "${AUTH_POSTGRES_PASSWORD:?AUTH_POSTGRES_PASSWORD is required}"

psql \
    --username "$POSTGRES_USER" \
    --dbname postgres \
    --set ON_ERROR_STOP=1 \
    --set auth_user="$AUTH_POSTGRES_USER" \
    --set auth_password="$AUTH_POSTGRES_PASSWORD" <<'SQL'
SELECT format(
    'CREATE ROLE %I LOGIN NOSUPERUSER NOCREATEDB NOCREATEROLE NOREPLICATION PASSWORD %L',
    :'auth_user',
    :'auth_password'
)
WHERE NOT EXISTS (
    SELECT 1
    FROM pg_roles
    WHERE rolname = :'auth_user'
)
\gexec

ALTER ROLE :"auth_user" PASSWORD :'auth_password';
SQL

psql \
    --username "$POSTGRES_USER" \
    --dbname postgres \
    --set ON_ERROR_STOP=1 \
    --set auth_db="$AUTH_POSTGRES_DB" \
    --set auth_user="$AUTH_POSTGRES_USER" <<'SQL'
SELECT format(
    'CREATE DATABASE %I OWNER %I',
    :'auth_db',
    :'auth_user'
)
WHERE NOT EXISTS (
    SELECT 1
    FROM pg_database
    WHERE datname = :'auth_db'
)
\gexec
SQL

psql \
    --username "$POSTGRES_USER" \
    --dbname "$AUTH_POSTGRES_DB" \
    --set ON_ERROR_STOP=1 \
    --set auth_user="$AUTH_POSTGRES_USER" <<'SQL'
REVOKE CREATE ON SCHEMA public FROM PUBLIC;

GRANT USAGE, CREATE ON SCHEMA public TO :"auth_user";

CREATE TABLE IF NOT EXISTS public.users (
    id UUID PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    phone VARCHAR(20) NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE public.users OWNER TO :"auth_user";
SQL