#! /bin/bash
# Copyright 2023-present Kuei-chun Chen. All rights reserved.
# build.sh

die() { echo "$*" 1>&2 ; exit 1; }
VERSION="v$(cat version)-$(git log -1 --date=format:"%Y%m%d" --format="%ad")"
REPO=$(basename "$(dirname "$(pwd)")")/$(basename "$(pwd)")
LDFLAGS="-X main.version=$VERSION -X main.repo=$REPO"
TAG="simagix/bond"
[[ "$(which go)" = "" ]] && die "go command not found"

gover=$(go version | cut -d' ' -f3)
if [ "$gover" \< "go1.18" ]; then
    [[ "$GOPATH" = "" ]] && die "GOPATH not set"
    [[ "${GOPATH}/src/github.com/$REPO" != "$(pwd)" ]] && die "building bond should be under ${GOPATH}/src/github.com/$REPO"
fi

if [ ! -f go.sum ]; then
    go mod tidy
fi

mkdir -p dist
if [ "$1" == "docker" ]; then
  BR=$(git branch --show-current)
  if [[ "${BR}" == "main" ]]; then
    BR="latest"
  fi
  docker build --no-cache -f Dockerfile -t ${TAG}:${BR} .
  docker run ${TAG}:${BR} /bond -version
  # docker rmi -f $(docker images -f "dangling=true" -q) > /dev/null 2>&1
elif [ "$1" == "dist" ]; then
  [[ "$(which uname)" = "" ]] && die "uname command not found"
  ofile="./dist/bond-$(uname|tr '[:upper:]' '[:lower:]')-$(uname -m)"
  go build -ldflags "$LDFLAGS" -o ${ofile} main/bond.go
else
  rm -f ./dist/bond
  go build -ldflags "$LDFLAGS" -o ./dist/bond main/bond.go
  if [[ -f ./dist/bond ]]; then
    ./dist/bond -version
  fi
fi
