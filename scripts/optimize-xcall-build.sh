#!/bin/bash
set -e
mkdir -p artifacts/icon
mkdir -p artifacts/archway
mkdir -p artifacts/evm
mkdir -p artifacts/sui
LOCAL_X_CALL_REPO=".xcall-multi"
LOCAL_ARTIFACT_XCALL_SUI="xcall"

clone_xCall_multi() {
  echo "Cloning xcall-multi repo..."
  X_CALL_BRANCH="${1:-main}"
  rm -rf "$LOCAL_X_CALL_REPO"
  git clone -b "$X_CALL_BRANCH" --single-branch "https://github.com/icon-project/xcall-multi.git" "$LOCAL_X_CALL_REPO"
  sed -i 's/docker-compose/docker compose/g' "${LOCAL_X_CALL_REPO}"/Makefile
  cd artifacts/sui
  git clone --bare -b "$X_CALL_BRANCH" --single-branch "https://github.com/icon-project/xcall-multi.git" "$LOCAL_ARTIFACT_XCALL_SUI"
  cd -
}

build_xCall_contracts() {
  echo "Generating optimized contracts of xcall-multi contracts..."
  clone_xCall_multi "${1:-main}"
  cd "$LOCAL_X_CALL_REPO"
  make build-wasm-docker
  make build-java-docker
  make build-solidity-docker
  cp artifacts/archway/*.wasm ../artifacts/archway/
  cp artifacts/icon/*.jar ../artifacts/icon/
  cp -R artifacts/evm/ ../artifacts/evm/
  cd -
}

build_sui_docker(){
  cd test/chains/sui/data
  docker build -t mysten/sui-tools-w-git .
  cd -
}

if [ "$1" = "build" ]; then
  shift
  build_xCall_contracts "$@"
  build_sui_docker
fi
