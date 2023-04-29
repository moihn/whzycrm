#!/bin/bash

docker run --rm -it -u $(id -u):$(id -g) --tmpfs $HOME -e HOME=$HOME -w /app --volume $PWD:/app mcr.microsoft.com/dotnet/sdk:6.0 dotnet publish -a arm64 -c Release

