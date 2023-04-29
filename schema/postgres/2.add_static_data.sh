#!/bin/sh
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

# run table creation as dba user
DB_NAME=${DB_NAME:-appdb}
PSQL_CMD=${PSQL_CMD:-docker compose exec -T db psql}

${PSQL_CMD} -U ${DB_NAME}_dba ${DB_NAME} <${SCRIPT_DIR}/static.sql

