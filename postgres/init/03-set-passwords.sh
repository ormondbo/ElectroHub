#!/usr/bin/env bash
set -Eeuo pipefail

: "${SHOP_APP_PASSWORD:?SHOP_APP_PASSWORD is required}"
: "${SHOP_READ_1_PASSWORD:?SHOP_READ_1_PASSWORD is required}"
: "${SHOP_READ_2_PASSWORD:?SHOP_READ_2_PASSWORD is required}"

psql \
    --username "$POSTGRES_USER" \
    --dbname "$POSTGRES_DB" \
    --set ON_ERROR_STOP=1 \
    --set app_password="$SHOP_APP_PASSWORD" \
    --set read_1_password="$SHOP_READ_1_PASSWORD" \
    --set read_2_password="$SHOP_READ_2_PASSWORD" <<'SQL'
ALTER ROLE shop_app PASSWORD :'app_password';
ALTER ROLE shop_read_1 PASSWORD :'read_1_password';
ALTER ROLE shop_read_2 PASSWORD :'read_2_password';
SQL
