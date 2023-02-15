#!/bin/bash

set -ex

sed -i "s/address = \"ojo1zypqa76je7pxsdwkfah6mu9a583sju6xzthge3\"/address = \"$ACCOUNT_ADDRESS\"/g" price-feeder.toml
sed -i "s/chain_id = \"ojo-testnet\"/chain_id = \"$CHAIN_ID\"/g" price-feeder.toml
sed -i "s/validator = \"ojovaloper1zypqa76je7pxsdwkfah6mu9a583sju6x6tnq6w\"/validator = \"$ACCOUNT_VALIDATOR\"/g" price-feeder.toml

sed -i "s/grpc_endpoint = \"localhost:9090\"/grpc_endpoint = \"$GRPC_ENDPOINT\"/g" price-feeder.toml
sed -i "s/tmrpc_endpoint = \"localhost:26657\"/tmrpc_endpoint = \"$TMRPC_ENDPOINT\"/g" price-feeder.toml

price-feeder price-feeder.toml --log-level debug
