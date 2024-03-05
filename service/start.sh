#!/bin/sh

echo "run db migration"

set -e

/app/migrate -path /app/migrations -database "$POSTGRESQL_URL" -verbose up

echo "start the app"
exec "$@"