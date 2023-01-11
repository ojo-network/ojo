# Ojo Devnet

```text
Chain ID: ojo-devnet
```

[Ojo Devnet Explorer](https://eyes.ojo.network.umee.cc/ojo-devnet)

Each official release, rc, or alpha release will trigger a complete network destruction and rebuild via github actions. The chain id and endpoints remain stable.

## Automation Notes

You can automatically fund your developer wallet by adding it to [this file](./pulumi/testnet/Pulumi.devnet.yaml) (and submitting a PR)

Example YAML Wallet Snippet:

```yaml
config:
  testnet:config:
    genesisAccounts:
    - address: ojo1gqjmdr64quxvjm5pt3ycgzyv4fc2zavc4z4zj7
      funding: [ 1000000000000uojo ]
```

## Load Balancer Endpoints

* [https://rpc.ojo-devnet.network.umee.cc](https://rpc.ojo-devnet.network.umee.cc)
* [https://api.ojo-devnet.network.umee.cc](https://api.ojo-devnet.network.umee.cc)

## Node Direct Endpoints

* [https://rpc.devnet-n0.ojo-devnet.network.umee.cc](https://rpc.devnet-n0.ojo-devnet.network.umee.cc)
* [https://rpc.devnet-n1.ojo-devnet.network.umee.cc](https://rpc.devnet-n1.ojo-devnet.network.umee.cc)
* [https://rpc.devnet-n2.ojo-devnet.network.umee.cc](https://rpc.devnet-n2.ojo-devnet.network.umee.cc)
* [https://api.devnet-n0.ojo-devnet.network.umee.cc](https://api.devnet-n0.ojo-devnet.network.umee.cc)
* [https://api.devnet-n1.ojo-devnet.network.umee.cc](https://api.devnet-n1.ojo-devnet.network.umee.cc)
* [https://api.devnet-n2.ojo-devnet.network.umee.cc](https://api.devnet-n2.ojo-devnet.network.umee.cc)
