#!/usr/bin/env bash

if [[ $# -le 0 ]]; then
    echo Require more than one argument to protoc.
    exit 1
fi

SCRIPTPATH="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
ROOTDIR="$(dirname "$SCRIPTPATH")"

vendor=false
if [[ $ROOTDIR = *"/vendor/istio.io/istio" ]]; then
  vendor=true
fi

# Ensure expected GOPATH setup
if [ "$vendor" = false ] && [ "$ROOTDIR" != "${GOPATH-$HOME/go}/src/istio.io/istio" ]; then
  die "Istio not found in GOPATH/src/istio.io/"
fi

gen_img=gcr.io/istio-testing/protoc:2018-06-12

docker run  -i --volume /var/run/docker.sock:/var/run/docker.sock \
  --rm --entrypoint /usr/bin/protoc -v "$ROOTDIR:$ROOTDIR" -w "$(pwd)" $gen_img "$@"
