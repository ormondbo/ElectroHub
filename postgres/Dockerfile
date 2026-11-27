FROM postgres:18-alpine

COPY postgres/init/*.sql /docker-entrypoint-initdb.d/
COPY postgres/init/*.sh /docker-entrypoint-initdb.d/

HEALTHCHECK --interval=5s --timeout=3s --retries=10 \
    CMD pg_isready -U "$POSTGRES_USER" -d "$POSTGRES_DB" || exit 1
