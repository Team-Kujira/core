#!/bin/bash
# microtick and bitcanna contributed significantly here.
# Pebbledb state sync script.
# invoke like: bash scripts/ss.bash



## USAGE RUNDOWN
# Not for use on live nodes
# For use when testing.
# Assumes that ~/.evmosd doesn't exist
# can be modified to suit your purposes if ~/.evmosd does already exist


set -uxe

# Set Golang environment variables.
export GOPATH=~/go
export PATH=$PATH:~/go/bin

# Install with pebbledb 
go mod edit -replace github.com/tendermint/tm-db=github.com/baabeetaa/tm-db@pebble
go mod tidy
go install -ldflags '-w -s -X github.com/cosmos/cosmos-sdk/types.DBBackend=pebbledb -X github.com/tendermint/tm-db.ForceSync=1' -tags pebbledb ./...

# go install ./...

# NOTE: ABOVE YOU CAN USE ALTERNATIVE DATABASES, HERE ARE THE EXACT COMMANDS
# go install -ldflags '-w -s -X github.com/cosmos/cosmos-sdk/types.DBBackend=rocksdb' -tags rocksdb ./...
# go install -ldflags '-w -s -X github.com/cosmos/cosmos-sdk/types.DBBackend=badgerdb' -tags badgerdb ./...
# go install -ldflags '-w -s -X github.com/cosmos/cosmos-sdk/types.DBBackend=boltdb' -tags boltdb ./...

# Initialize chain.
kujirad init test

# Get Genesis
wget -O ~/.kujira/config/genesis.json https://github.com/Team-Kujira/networks/raw/master/mainnet/kaiyo-1.json


# Get "trust_hash" and "trust_height".
INTERVAL=1000
LATEST_HEIGHT=$(curl -s https://kujira-rpc.polkachu.com/block | jq -r .result.block.header.height)
BLOCK_HEIGHT=$(($LATEST_HEIGHT-$INTERVAL)) 
TRUST_HASH=$(curl -s "https://kujira-rpc.polkachu.com/block?height=$BLOCK_HEIGHT" | jq -r .result.block_id.hash)

# Print out block and transaction hash from which to sync state.
echo "trust_height: $BLOCK_HEIGHT"
echo "trust_hash: $TRUST_HASH"

# Export state sync variables.
export KUJIRAD_STATESYNC_ENABLE=true
export KUJIRAD_P2P_MAX_NUM_OUTBOUND_PEERS=200
export KUJIRAD_STATESYNC_RPC_SERVERS="https://kujira-rpc.polkachu.com:443,https://kujira-rpc.polkachu.com:443"
export KUJIRAD_STATESYNC_TRUST_HEIGHT=$BLOCK_HEIGHT
export KUJIRAD_STATESYNC_TRUST_HASH=$TRUST_HASH

# Fetch and set list of seeds from chain registry.
export KUJIRAD_P2P_SEEDS=$(curl -s https://raw.githubusercontent.com/cosmos/chain-registry/master/kujira/chain.json | jq -r '[foreach .peers.seeds[] as $item (""; "\($item.id)@\($item.address)")] | join(",")')

# Start chain.
kujirad start --x-crisis-skip-assert-invariants --db_backend pebbledb
