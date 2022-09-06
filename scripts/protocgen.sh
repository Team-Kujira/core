#!/usr/bin/env bash

set -eo pipefail

# get protoc executions
go get github.com/regen-network/cosmos-proto/protoc-gen-gocosmos 2>/dev/null

# get cosmos sdk from github
go get github.com/cosmos/cosmos-sdk@v0.45.4 2>/dev/null

echo "Generating gogo proto code"
cd proto

buf mod update

proto_dirs=$(find ./kujira -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do
  for file in $(find "${dir}" -maxdepth 1 -name '*.proto'); do
    if grep go_package $file &>/dev/null; then
      buf generate --template buf.gen.gogo.yaml $file
    fi
  done
done

cd ..

# move proto files to the right places
#
# Note: Proto files are suffixed with the current binary version.
cp -r kujira/* ./
rm -rf kujira

go mod tidy -compat=1.18