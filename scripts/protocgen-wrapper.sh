#!/bin/sh

if which docker > /dev/null; then
    docker run --volume "$(pwd)\:/workspace" --workdir /workspace ghcr.io/cosmos/proto-builder:0.12.1 sh ./scripts/protocgen.sh
elif which podman > /dev/null; then
    podman run --rm -v $(pwd):/workspace --workdir /workspace -u root ghcr.io/cosmos/proto-builder:0.12.1 sh ./scripts/protocgen.sh
fi