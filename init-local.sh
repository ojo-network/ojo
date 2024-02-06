#!/bin/sh

rm -r ~/.ojo

OJO_RPC=26657
OJO_P2P=26656
CELESTIA_RPC=36657

# Start the Celestia Devnet
CONTAINER_NAME="celestia_devnet"
running=$(docker ps --filter name=$CONTAINER_NAME --format "{{.Names}}" | grep $CONTAINER_NAME)
if [ -z "$running" ]; then
    echo "Starting celestia devnet"
else
    echo "Celestia devnet is already running. Restarting..."
    docker rm -f $CONTAINER_NAME
fi

docker run -t -d --name $CONTAINER_NAME \
    -p 36650:26650 -p 36657:26657 -p 36658:26658 -p 36659:26659 -p 9091:9090 \
    ghcr.io/rollkit/local-celestia-devnet:v0.12.5
sleep 5

# set variables for the chain
VALIDATOR_NAME=validator1
CHAIN_ID=ojo-dev-01
KEY_NAME=ojo-key
KEY_2_NAME=ojo-key-2
CHAINFLAG="--chain-id ${CHAIN_ID}"
TOKEN_AMOUNT="10000000000000000000000000uojo"
STAKING_AMOUNT="1000000000uojo"

# create a random Namespace ID for your rollup to post blocks to
NAMESPACE=$(openssl rand -hex 8)

# query the DA Layer start height, in this case we are querying
# our local devnet at port 26657, the RPC. The RPC endpoint is
# to allow users to interact with Celestia's nodes by querying
# the node's state and broadcasting transactions on the Celestia
# network. The default port is 26657.
DA_BLOCK_HEIGHT=$(curl http://0.0.0.0:$CELESTIA_RPC/block | jq -r '.result.block.header.height')

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
echo -e "\n\n\n\n\n Your NAMESPACE is $NAMESPACE \n\n Your DA_BLOCK_HEIGHT is $DA_BLOCK_HEIGHT \n\n\n\n\n"

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

# set the staking amounts in the genesis transaction
ojod gentx $KEY_NAME $STAKING_AMOUNT --chain-id $CHAIN_ID --keyring-backend test

# collect genesis transactions
ojod collect-gentxs

# copy centralized sequencer address into genesis.json
# Note: validator and sequencer are used interchangeably here
ADDRESS=$(jq -r '.address' ~/.ojo/config/priv_validator_key.json)
PUB_KEY=$(jq -r '.pub_key' ~/.ojo/config/priv_validator_key.json)
jq --argjson pubKey "$PUB_KEY" '. + {"validators": [{"address": "'$ADDRESS'", "pub_key": $pubKey, "power": "1000", "name": "Rollkit Sequencer"}]}' ~/.ojo/config/genesis.json > temp.json && mv temp.json ~/.ojo/config/genesis.json

AUTH_TOKEN=$(celestia light auth write)

# create a restart-local.sh file to restart the chain later
[ -f restart-local.sh ] && rm restart-local.sh
echo "DA_BLOCK_HEIGHT=$DA_BLOCK_HEIGHT" >> restart-local.sh
echo "NAMESPACE=$NAMESPACE" >> restart-local.sh
echo "AUTH_TOKEN=$AUTH_TOKEN" >> restart-local.sh

echo "ojod start --rollkit.aggregator true --rollkit.da_layer celestia --rollkit.da_config='{\"base_url\":\"http://localhost:$CELESTIA_RPC\",\"timeout\":60000000000,\"fee\":600000,\"gas_limit\":6000000,\"auth_token\":\"'\$AUTH_TOKEN'\"}' --rollkit.namespace_id \$NAMESPACE --rollkit.da_start_height \$DA_BLOCK_HEIGHT --rpc.laddr tcp://127.0.0.1:$OJO_RPC --minimum-gas-prices="0.025uojo" --p2p.laddr \"0.0.0.0:$OJO_P2P\"" >> restart-local.sh

# start the chain
ojod start --rollkit.aggregator true --rollkit.da_layer celestia --rollkit.da_config='{"base_url":"http://localhost:'$CELESTIA_RPC'","timeout":60000000000,"fee":600000,"gas_limit":6000000,"auth_token":"'$AUTH_TOKEN'"}' --rollkit.namespace_id $NAMESPACE --rollkit.da_start_height $DA_BLOCK_HEIGHT --rpc.laddr tcp://127.0.0.1:$OJO_RPC --p2p.laddr "0.0.0.0:$OJO_P2P" --minimum-gas-prices="0.025uojo"

# uncomment the next command if you are using lazy aggregation
# ojod start --rollkit.aggregator true --rollkit.da_layer celestia --rollkit.da_config='{"base_url":"http://localhost:26658","timeout":60000000000,"fee":600000,"gas_limit":6000000,"auth_token":"'$AUTH_TOKEN'"}' --rollkit.namespace_id $NAMESPACE --rollkit.da_start_height $DA_BLOCK_HEIGHT --rollkit.lazy_aggregator
