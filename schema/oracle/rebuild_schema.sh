#!/bin/bash

source .env

PDB_ADMIN=system
PDB_ADMIN_CONN_STRING="${PDB_ADMIN}/${PDB_ADMIN_PWD}@//localhost:1521/XEPDB1"
PDB_APP_CONN_STRING="${APP_DB_USER}/${APP_DB_PWD}@//localhost:1521/XEPDB1"
SQLPLUS_EXEC="docker exec -i ${DB_DOCKER_CONTAINER} sqlplus -S -L"

# create schema user in database
${SQLPLUS_EXEC} "${PDB_ADMIN_CONN_STRING}" <<EOF
-- create schema user
DEFINE NEW_USER = '${APP_DB_USER^^}'
DEFINE NEW_USER_PW = '${APP_DB_PWD}'

drop user &NEW_USER cascade
/

create user &NEW_USER IDENTIFIED by &NEW_USER_PW
/

GRANT CREATE TABLE TO &NEW_USER
/

GRANT CREATE SEQUENCE TO &NEW_USER
/

GRANT CREATE SESSION TO &NEW_USER
/

GRANT CREATE ANY PROCEDURE TO &NEW_USER
/

-- ALTER USER &NEW_USER quota unlimited on data
-- /
EOF

# use application db user to create objects in schema
while IFS= read -r line || [ -n "$line" ]
do
    echo "Loading ${line}" 
    cat "${line}" | ${SQLPLUS_EXEC} "${PDB_APP_CONN_STRING}"
done <objects.lst
