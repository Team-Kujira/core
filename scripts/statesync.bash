#!/bin/bash
# microtick and bitcanna contributed significantly here.
# invoke this script in your migaloo-chain folder like this:
# bash scripts/statesync.bash
set -uxe

# Set Golang environment variables.
export GOPATH=~/go
export PATH=$PATH:~/go/bin

# Install Migaloo.
go install ./...

# NOTE: ABOVE YOU CAN USE ALTERNATIVE DATABASES, HERE ARE THE EXACT COMMANDS
# go install -ldflags '-w -s -X github.com/cosmos/cosmos-sdk/types.DBBackend=rocksdb' -tags rocksdb ./...
# go install -ldflags '-w -s -X github.com/cosmos/cosmos-sdk/types.DBBackend=badgerdb' -tags badgerdb ./...
# go install -ldflags '-w -s -X github.com/cosmos/cosmos-sdk/types.DBBackend=boltdb' -tags boltdb ./...

# Check if genesis file exist, if not then download
GENESIS_FILE="$HOME/.migalood/config/genesis.json"
GENESIS_URL="https://github.com/White-Whale-Defi-Platform/migaloo-chain/raw/release/v2.0.x/networks/mainnet/genesis.json"

if [ ! -f "$GENESIS_FILE" ]; then
    # Initialize chain.
    migalood init test 
    mkdir -p "$(dirname "$GENESIS_FILE")"
    # Download genesis
    wget -O "$GENESIS_FILE" "$GENESIS_URL"
    
else
    echo "File already exists!!!"   
fi


# Set minimum gas price.
#  ADD SPACE IF RUN ON MAC
sed -i'' 's/minimum-gas-prices = ""/minimum-gas-prices = "0.0025uwhale"/' "$HOME/.migalood/config/app.toml"

# Get "trust_hash" and "trust_height".
INTERVAL=1000
LATEST_HEIGHT=$(curl -s https://whitewhale-rpc.lavenderfive.com/block | jq -r .result.block.header.height)
BLOCK_HEIGHT=$((LATEST_HEIGHT-INTERVAL)) 
TRUST_HASH=$(curl -s "https://whitewhale-rpc.lavenderfive.com/block?height=$BLOCK_HEIGHT" | jq -r .result.block_id.hash)

# Print out block and transaction hash from which to sync state.
echo "trust_height: $BLOCK_HEIGHT"
echo "trust_hash: $TRUST_HASH"

# Export state sync variables.
export MIGALOOD_STATESYNC_ENABLE=true
export MIGALOOD_P2P_MAX_NUM_OUTBOUND_PEERS=200
# replace the url below with a working one, get it from chain registry
export MIGALOOD_STATESYNC_RPC_SERVERS="https://whitewhale-mainnet-rpc.autostake.net:443,https://rpc-whitewhale.goldenratiostaking.net,https://whitewhale-rpc.lavenderfive.com,https://rpc-whitewhale.carbonzero.zone:443,https://rpc-whitewhale.whispernode.com:443,https://migaloo-rpc.kleomedes.network:443"
export MIGALOOD_STATESYNC_TRUST_HEIGHT=$BLOCK_HEIGHT
export MIGALOOD_STATESYNC_TRUST_HASH=$TRUST_HASH

# Fetch and set list of seeds from chain registry.
MIGALOOD_P2P_SEEDS=$(curl -s https://raw.githubusercontent.com/cosmos/chain-registry/master/migaloo/chain.json | jq -r '[foreach .peers.seeds[] as $item (""; "\($item.id)@\($item.address)")] | join(",")')
export MIGALOOD_P2P_SEEDS


# Start chain.
migalood start --x-crisis-skip-assert-invariants --minimum-gas-prices 0.00001uwhale