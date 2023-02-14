#!/bin/bash

set -ex

ojod init val01 --chain-id=$OJO_CHAIN_ID
echo $MNEMONIC | ojod keys add val01 --recover --keyring-backend=test
ojod add-genesis-account $(ojod keys show val01 -a --keyring-backend=test) 1000000000000stake,500000000000uojo
ojod gentx val01 500000000000stake --chain-id=test-osmosis --keyring-backend=test
ojod collect-gentxs

ojod start
