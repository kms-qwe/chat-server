export MIGRATION_DSN="host=$DB_HOST port=$PG_INNER_PORT dbname=$POSTGRES_DB user=$POSTGRES_USER password=$POSTGRES_PASSWORD sslmode=disable"

sleep 3 && goose -dir "${MIGRATION_DIR}" postgres "${MIGRATION_DSN}" up -v && echo "SUCCESS"