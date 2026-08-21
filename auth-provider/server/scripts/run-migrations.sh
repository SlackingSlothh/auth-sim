#!/bin/sh
set -e
# Retry running migrate until it succeeds. Compose already waits for the
# DB container healthcheck, so this loop only handles transient connection
# timing issues.
# echo "Environment vars:"
# env | grep -E 'DB_|POSTGRES_' || true

MAX_ATTEMPTS=20
DELAY=2

echo "Starting migrations (will retry up to $MAX_ATTEMPTS times)"
for attempt in $(seq 1 $MAX_ATTEMPTS); do
  echo "Attempt $attempt/$MAX_ATTEMPTS: running migrate"
  if migrate -path /migrations -database "postgres://$DB_USER:$DB_PASSWORD@db_primary:5432/$DB_NAME?sslmode=disable" up; then
    echo "Migrate succeeded"
    exit 0
  fi
  echo "Migrate failed on attempt $attempt — retrying in $DELAY seconds"
  sleep $DELAY
done


echo "Migrate failed after $MAX_ATTEMPTS attempts"
exit 1
