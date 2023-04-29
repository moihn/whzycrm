#!/bin/sh

set -x

export PROJ_NAME=${PROJ_NAME:-whzycrmqa}
export DB_NAME=${DB_NAME:-whzydb}
export DB_DBA_PW=${DB_DBA_PW:-dbasecret}
export DB_APPUSER_PW=${DB_APPUSER_PW:-appsecret}
export TMP_SCHEMA=${TMP_SCHEMA:-tmp}
export PSQL_CMD=${PSQL_CMD:-docker compose -p ${PROJ_NAME} exec -T db psql}

# sort *.sh script files in the folder and run them in alphabetic order
SORTED_SCRIPTS=$(ls)
for script in "${SORTED_SCRIPTS}"
do
	sh -x ${script}
done
