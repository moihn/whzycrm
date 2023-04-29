set -e

export GOBIN="$PWD/cmd/whzy_rest_api"
(cd cmd/whzy_rest_api && go install)
ARCH=$(arch)
docker build -t whzycrm . --build-arg BS_USER=$(id -un) --build-arg BS_UID=$(id -u) --build-arg BS_GROUP=$(id -gn) --build-arg BS_GID=$(id -g) -f Dockerfile."$ARCH"
