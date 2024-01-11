#!/bin/sh

BLOCK_TIME="${BLOCK_TIME:-5}"

if [ -z "$GENESIS_PATH" ]; then
  echo "GENESIS_PATH not set. Running anvil with only --block-time option."
  anvil --block-time "$BLOCK_TIME" --host 0.0.0.0
else
  anvil --block-time "$BLOCK_TIME" --init "$GENESIS_PATH" --host 0.0.0.0
fi

set -e
if [ "${1#-}" != "${1}" ] || [ -z "$(command -v "${1}")" ]; then
  set -- node "$@"
fi
exec "$@"
