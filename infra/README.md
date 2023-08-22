# Ojo Devnet

```text
Chain ID: ojo-internal-devnet
```

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
