#!/bin/sh
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

# run table creation as dba user
DB_NAME=${DB_NAME:-appdb}
TMP_SCHEMA=${TMP_SCHEMA:-appuser}
PSQL_CMD=${PSQL_CMD:-docker compose exec -T db psql}

${PSQL_CMD} -U ${DB_NAME}_appuser -d ${DB_NAME} <<EOF

	insert into vendor
	(name)
	values
	('鸿鲲饰品'),
	('天友玩具'),
	('展希');
	
EOF
