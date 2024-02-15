DA_BLOCK_HEIGHT=4
NAMESPACE=620dd396f249a456
ojod start --rollkit.aggregator --rollkit.da_address=http://localhost:$CELESTIA_RPC --rollkit.da_start_height 4 --rpc.laddr tcp://127.0.0.1:26657 --grpc.address 127.0.0.1:9290 --p2p.laddr "0.0.0.0:26656" --minimum-gas-prices=0.025uojo
