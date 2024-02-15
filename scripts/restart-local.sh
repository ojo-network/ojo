DA_BLOCK_HEIGHT=2403
ojod start --rollkit.aggregator --rollkit.da_address=:26650 --rollkit.da_start_height $DA_BLOCK_HEIGHT --rpc.laddr tcp://127.0.0.1:36657 --grpc.address 127.0.0.1:9290 --p2p.laddr "0.0.0.0:36656" --minimum-gas-prices=0.025uojo
