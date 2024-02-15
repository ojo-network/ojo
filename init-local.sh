#!/bin/sh

rm -r ~/.ojo

DENOM="${DENOM:-uojo}"
STAKE_DENOM="${STAKE_DENOM:-$DENOM}"

# set variables for the chain
VALIDATOR_NAME=validator1
CHAIN_ID=ojo-dev-01
KEY_NAME=ojo-key
KEY_2_NAME=ojo-key-2
CHAINFLAG="--chain-id ${CHAIN_ID}"
TOKEN_AMOUNT="10000000000000000000000000uojo"
STAKING_AMOUNT="1000000000uojo"

# query the DA Layer start height, in this case we are querying
# our local devnet at port 26657, the RPC. The RPC endpoint is
# to allow users to interact with Celestia's nodes by querying
# the node's state and broadcasting transactions on the Celestia
# network. The default port is 26657.
DA_BLOCK_HEIGHT=$(curl http://0.0.0.0:26657/block | jq -r '.result.block.header.height')

# rollkit logo
cat <<'EOF'

                 :=+++=.
              -++-    .-++:
          .=+=.           :++-.
       -++-                  .=+=: .
   .=+=:                        -%@@@*
  +%-                       .=#@@@@@@*
    -++-                 -*%@@@@@@%+:
       .=*=.         .=#@@@@@@@%=.
      -++-.-++:    =*#@@@@@%+:.-++-=-
  .=+=.       :=+=.-: @@#=.   .-*@@@@%
  =*=:           .-==+-    :+#@@@@@@%-
     :++-               -*@@@@@@@#=:
        =%+=.       .=#@@@@@@@#%:
     -++:   -++-   *+=@@@@%+:   =#*##-
  =*=.         :=+=---@*=.   .=*@@@@@%
  .-+=:            :-:    :+%@@@@@@%+.
      :=+-             -*@@@@@@@#=.
         .=+=:     .=#@@@@@@%*-
             -++-  *=.@@@#+:
                .====+*-.

   ______         _  _  _     _  _
   | ___ \       | || || |   (_)| |
   | |_/ /  ___  | || || | __ _ | |_
   |    /  / _ \ | || || |/ /| || __|
   | |\ \ | (_) || || ||   < | || |_
   \_| \_| \___/ |_||_||_|\_\|_| \__|
EOF

# echo variables for the chain
echo -e "\n Your DA_BLOCK_HEIGHT is $DA_BLOCK_HEIGHT \n"

# build the ojo chain with Rollkit
make install

# reset any existing genesis/chain data
ojod tendermint unsafe-reset-all

# initialize the validator with the chain ID you set
ojod init $VALIDATOR_NAME --chain-id $CHAIN_ID

# add keys for key 1 and key 2 to keyring-backend test
ojod keys add $KEY_NAME --keyring-backend test
ojod keys add $KEY_2_NAME --keyring-backend test

# add these as genesis accounts
ojod add-genesis-account $KEY_NAME $TOKEN_AMOUNT --keyring-backend test
ojod add-genesis-account $KEY_2_NAME $TOKEN_AMOUNT --keyring-backend test

# patch genesis
echo "--- Patching genesis..."
jq '.consensus_params["block"]["time_iota_ms"]="5000"
    | .app_state["crisis"]["constant_fee"]["denom"]="'$DENOM'"
    | .app_state["gov"]["deposit_params"]["min_deposit"][0]["denom"]="'$DENOM'"
    | .app_state["mint"]["params"]["mint_denom"]="'$DENOM'"
    | .app_state["staking"]["params"]["bond_denom"]="'$DENOM'"
    | .app_state["gov"]["voting_params"]["voting_period"]="10s"' \
    ~/.ojo/config/genesis.json > temp.json && mv temp.json ~/.ojo/config/genesis.json

# set the staking amounts in the genesis transaction
ojod gentx $KEY_NAME $STAKING_AMOUNT --chain-id $CHAIN_ID --keyring-backend test

# collect genesis transactions
ojod collect-gentxs

# copy centralized sequencer address into genesis.json
# Note: validator and sequencer are used interchangeably here
ADDRESS=$(jq -r '.address' ~/.ojo/config/priv_validator_key.json)
PUB_KEY=$(jq -r '.pub_key' ~/.ojo/config/priv_validator_key.json)
jq --argjson pubKey "$PUB_KEY" '.consensus["validators"]=[{"address": "'$ADDRESS'", "pub_key": $pubKey, "power": "1000", "name": "Rollkit Sequencer"}]' ~/.ojo/config/genesis.json > temp.json && mv temp.json ~/.ojo/config/genesis.json

# create a restart-local.sh file to restart the chain later
[ -f restart-local.sh ] && rm restart-local.sh
echo "DA_BLOCK_HEIGHT=$DA_BLOCK_HEIGHT" >> restart-local.sh

echo "ojod start --rollkit.aggregator --rollkit.da_address=":26650" --rollkit.da_start_height \$DA_BLOCK_HEIGHT --rpc.laddr tcp://127.0.0.1:36657 --grpc.address 127.0.0.1:9290 --p2p.laddr \"0.0.0.0:36656\" --minimum-gas-prices="0.025uojo"" >> restart-local.sh

# start the chain
ojod start --rollkit.aggregator --rollkit.da_address=":26650" --rollkit.da_start_height $DA_BLOCK_HEIGHT --rpc.laddr tcp://127.0.0.1:36657 --grpc.address 127.0.0.1:9290 --p2p.laddr "0.0.0.0:36656" --minimum-gas-prices="0.025uojo"
