#!/bin/bash

set -ex

ojod init $NODE_MONIKER --chain-id=$OJO_CHAIN_ID

echo $VAL1_MNEMONIC | ojod keys add $VAL1_MONIKER --recover --keyring-backend=test
ojod add-genesis-account $(ojod keys show $VAL1_MONIKER -a --keyring-backend=test) 600000000000uojo
ojod gentx $VAL1_MONIKER 500000000000uojo --chain-id=$OJO_CHAIN_ID --keyring-backend=test

echo $VAL2_MNEMONIC | ojod keys add $VAL2_MONIKER --recover --keyring-backend=test
ojod add-genesis-account $(ojod keys show $VAL2_MONIKER -a --keyring-backend=test) 600000000000uojo
ojod gentx $VAL2_MONIKER 500000000000uojo --chain-id=$OJO_CHAIN_ID --keyring-backend=test

ojod collect-gentxs

ojod start
