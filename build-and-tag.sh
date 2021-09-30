#!/usr/bin/env bash

set -e

function buildAndTag() {
  local version=$1
  local imagename="init-containers/vault-cert-generator"
  local fromline
  fromline=$(grep -e '^FROM ' Dockerfile | tail -n -1 | sed 's/^FROM[ \t]*//' | sed 's#.*/##' | sed 's/:/-/' | sed 's/#.*//' | sed -E 's/ +.*//')

  docker build \
    -t "${imagename}:${version}" \
    -t "${imagename}:${version}-${fromline}" \
    .
}

make clean
buildAndTag "0.1"
