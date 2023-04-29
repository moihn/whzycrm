Build Instruction
1. Create database container
* Get Oracle-XE Dockerfile from https://github.com/oracle/docker-images/tree/main/OracleDatabase/SingleInstance and follow the instruction to create Oracle 21c image.
* Create schema/.env file to add setup for below list of environment variables, using format key=value:
 - APP_DB_USER
 - APP_DB_PWD
 - DB_DOCKER_CONTAINER
 - PDB_ADMIN_PWD
* Create Oracle-XE 21c container using the docker/docker-compose.yml file.
* Run schema build using schema/rebuild_schema.sh