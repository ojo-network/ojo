#!/bin/bash -eu

# USAGE:
# ./single-gen.sh <option of full path to ojod>

# Starts an ojo chain with only a single node. Best used with an ojod bin
# sitting in the same folder as the script, rather than using the one installed.
# Useful for upgrade testing, where two ojod versions can be placed in the
# folder to test.

# Without submitting any governance proposals, it seems ojod 1 and 2 releases
# can just start and continue off the same state back and forth without failing.
# e.g. run this with ojod1, stop ojod1, then run it with ojod2 to continue.

CWD="$( cd -- "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"
PRICE_FEEDER_CONFIG_PATH="${CWD}/../pricefeeder/price-feeder.example.toml"
export PRICE_FEEDER_CONFIG=$(realpath "$PRICE_FEEDER_CONFIG_PATH")
export PRICE_FEEDER_CHAIN_CONFIG="TRUE"
export PRICE_FEEDER_LOG_LEVEL="DEBUG"

NODE_BIN="${1:-$CWD/../build/ojod}"

# These options can be overridden by env
CHAIN_ID="${CHAIN_ID:-ojotest-1}"
CHAIN_DIR="${CHAIN_DIR:-$CWD/node-data}"
DENOM="${DENOM:-uojo}"
STAKE_DENOM="${STAKE_DENOM:-$DENOM}"
CLEANUP="${CLEANUP:-1}"
LOG_LEVEL="${LOG_LEVEL:-info}"
SCALE_FACTOR="${SCALE_FACTOR:-000000}"
VOTING_PERIOD="${VOTING_PERIOD:-20s}"

# Default 1 account keys + 1 user key with no special grants
VAL0_KEY="val"
VAL0_MNEMONIC="copper push brief egg scan entry inform record adjust fossil boss egg comic alien upon aspect dry avoid interest fury window hint race symptom"
VAL0_ADDR="ojo1y6xz2ggfc0pcsmyjlekh0j9pxh6hk87ymc9due"

USER_KEY="user"
USER_MNEMONIC="pony glide frown crisp unfold lawn cup loan trial govern usual matrix theory wash fresh address pioneer between meadow visa buffalo keep gallery swear"
USER_ADDR="ojo1usr9g5a4s2qrwl63sdjtrs2qd4a7huh6cuuhrc"

NEWLINE=$'\n'

hdir="$CHAIN_DIR/$CHAIN_ID"

if ! command -v jq &> /dev/null
then
  echo "⚠️ jq command could not be found!"
  echo "Install it by checking https://stedolan.github.io/jq/download/"
  exit 1
fi

echo "--- Chain ID = $CHAIN_ID"
echo "--- Chain Dir = $CHAIN_DIR"
echo "--- Coin Denom = $DENOM"
VERSION=$($NODE_BIN version)
echo "--- Binary Version = $VERSION"

killall "$NODE_BIN" &>/dev/null || true

if [[ "$CLEANUP" == 1 || "$CLEANUP" == "1" ]]; then
  rm -rf "$CHAIN_DIR"
  echo "Removed $CHAIN_DIR"
fi

# Folder for node
n0dir="$hdir/n0"

# Home flag for folder
home0="--home $n0dir"

# Config directories for node
n0cfgDir="$n0dir/config"

# Config files for nodes
n0cfg="$n0cfgDir/config.toml"

# App config file for node
n0app="$n0cfgDir/app.toml"

# Common flags
kbt="--keyring-backend test"
cid="--chain-id $CHAIN_ID"

# Check if the node-data dir has been initialized already
if [[ ! -d "$hdir" ]]; then
  echo "====================================="
  echo "STARTING NEW CHAIN WITH GENESIS STATE"
  echo "====================================="

  echo "--- Creating $NODE_BIN validator with chain-id=$CHAIN_ID"

  # Build genesis file and create accounts
  if [[ "$STAKE_DENOM" != "$DENOM" ]]; then
    coins="1000000$SCALE_FACTOR$STAKE_DENOM,1000000$SCALE_FACTOR$DENOM"
  else
    coins="1000000$SCALE_FACTOR$DENOM"
  fi
  coins_user="1000000$SCALE_FACTOR$DENOM"

  echo "--- Initializing home..."

  # Initialize the home directory of node
  $NODE_BIN $home0 $cid init n0

  echo "--- Enabling node API"
  sed -i -s '108s/enable = false/enable = true/' $n0app

  # Generate new random key
  # $NODE_BIN $home0 keys add val $kbt &>/dev/null

  echo "--- Importing keys..."
  echo "$VAL0_MNEMONIC$NEWLINE"
  yes "$VAL0_MNEMONIC$NEWLINE" | $NODE_BIN $home0 keys add $VAL0_KEY $kbt --recover
  yes "$USER_MNEMONIC$NEWLINE" | $NODE_BIN $home0 keys add $USER_KEY $kbt --recover

  echo "--- Adding addresses..."
  $NODE_BIN $home0 keys show $VAL0_KEY -a $kbt
  $NODE_BIN $home0 keys show $VAL0_KEY -a --bech val $kbt
  $NODE_BIN $home0 keys show $USER_KEY -a $kbt
  $NODE_BIN $home0 add-genesis-account $($NODE_BIN $home0 keys show $VAL0_KEY -a $kbt) $coins &>/dev/null
  $NODE_BIN $home0 add-genesis-account $($NODE_BIN $home0 keys show $USER_KEY -a $kbt) $coins_user &>/dev/null


  echo "--- Patching genesis..."
  if [[ "$STAKE_DENOM" == "$DENOM" ]]; then
    jq '.consensus_params["block"]["time_iota_ms"]="5000"
      | .app_state["crisis"]["constant_fee"]["denom"]="'$DENOM'"
      | .app_state["gov"]["deposit_params"]["min_deposit"][0]["denom"]="'$DENOM'"
      | .app_state["mint"]["params"]["mint_denom"]="'$DENOM'"
      | .app_state["staking"]["params"]["bond_denom"]="'$DENOM'"
      | .app_state["gov"]["voting_params"]["voting_period"]="10s"' \
        $n0cfgDir/genesis.json > $n0cfgDir/tmp_genesis.json && mv $n0cfgDir/tmp_genesis.json $n0cfgDir/genesis.json

  fi

  jq '.app_state["gov"]["voting_params"]["voting_period"]="'$VOTING_PERIOD'"' $n0cfgDir/genesis.json > $n0cfgDir/tmp_genesis.json && mv $n0cfgDir/tmp_genesis.json $n0cfgDir/genesis.json

  echo "--- Creating gentx..."
  $NODE_BIN $home0 gentx $VAL0_KEY 1000$SCALE_FACTOR$STAKE_DENOM $kbt $cid

  $NODE_BIN $home0 collect-gentxs > /dev/null

  echo "--- Validating genesis..."
  $NODE_BIN $home0 validate-genesis

  echo "--- Final Genesis ---"
  cat $n0cfgDir/genesis.json

  # Use perl for cross-platform compatibility
  # Example usage: perl -i -pe 's/^param = ".*?"/param = "100"/' config.toml

  echo "--- Modifying config..."
  perl -i -pe 's|addr_book_strict = true|addr_book_strict = false|g' $n0cfg
  perl -i -pe 's|external_address = ""|external_address = "tcp://127.0.0.1:26657"|g' $n0cfg
  perl -i -pe 's|"tcp://127.0.0.1:26657"|"tcp://0.0.0.0:26657"|g' $n0cfg
  perl -i -pe 's|allow_duplicate_ip = false|allow_duplicate_ip = true|g' $n0cfg
  perl -i -pe 's|log_level = "info"|log_level = "'$LOG_LEVEL'"|g' $n0cfg
  perl -i -pe 's|timeout_commit = ".*?"|timeout_commit = "5s"|g' $n0cfg

  echo "--- Modifying app..."
  perl -i -pe 's|minimum-gas-prices = ""|minimum-gas-prices = "0.05uojo"|g' $n0app

else
  echo "===================================="
  echo "CONTINUING CHAIN FROM PREVIOUS STATE"
  echo "===================================="
fi # data dir check

log_path=$hdir.n0.log

# Start the instance
echo "--- Starting node..."
echo
echo "Logs:"
echo "  * tail -f $log_path"
echo
echo "Env for easy access:"
echo "export H1='--home $hdir'"
echo
echo "Command Line Access:"
echo "  * $NODE_BIN --home $hdir status"

$NODE_BIN $home0 start --api.enable true --grpc.address="0.0.0.0:9090" --grpc-web.enable=false --log_level trace > $log_path 2>&1 &

# Adds 1 sec to create the log and makes it easier to debug it on CI
sleep 1

cat $log_path
