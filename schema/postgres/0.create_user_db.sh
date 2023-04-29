#!/bin/sh

# run this script as postgres UNIX user
DB_NAME=${DB_NAME:-appdb}
DB_DBA_PW=${DB_DBA_PW:-dbasecret}
DB_APPUSER_PW=${DB_APPUSER_PW:-appsecret}
PSQL_CMD=${PSQL_CMD:-docker compose exec -T db psql}

# 1. create a whzy_dba user that can login
${PSQL_CMD} -U postgres <<EOF
CREATE USER ${DB_NAME}_dba PASSWORD '${DB_DBA_PW}';
CREATE DATABASE ${DB_NAME} OWNER ${DB_NAME}_dba;
CREATE USER ${DB_NAME}_appuser PASSWORD '${DB_APPUSER_PW}';
GRANT pg_read_server_files TO ${DB_NAME}_appuser;
CREATE SCHEMA tmp AUTHORIZATION ${DB_NAME}_appuser;
EOF
