#!/bin/bash
set -e
mkdir -p artifacts/icon
mkdir -p artifacts/archway
mkdir -p artifacts/evm
LOCAL_X_CALL_REPO=".xcall-multi"

clone_xCall_multi() {
  echo "Cloning xcall-multi repo..."
  X_CALL_BRANCH="${1:-main}"
  rm -rf "$LOCAL_X_CALL_REPO"
  git clone -b "$X_CALL_BRANCH" --single-branch "https://github.com/icon-project/xcall-multi.git" "$LOCAL_X_CALL_REPO"
  sed -i 's/docker-compose/docker compose/g' "${LOCAL_X_CALL_REPO}"/Makefile
}

build_xCall_contracts() {
  echo "Generating optimized contracts of xcall-multi contracts..."
  clone_xCall_multi "${1:-main}"
  cd "$LOCAL_X_CALL_REPO"
  #  ./scripts/optimize-cosmwasm.sh //not required right now
  make build-java-docker
  make build-solidity-docker
  #  cp artifacts/archway/*.wasm ../artifacts/archway/
  cp artifacts/icon/*.jar ../artifacts/icon/
  cp -R artifacts/evm/ ../artifacts/evm/
  cd -
}

if [ "$1" = "build" ]; then
  shift
  build_xCall_contracts "$@"
fi
