docker build -t whzycrm . --build-arg BS_USER=$(id -un) --build-arg BS_UID=$(id -u) --build-arg BS_GROUP=$(id -gn) --build-arg BS_GID=$(id -g)
