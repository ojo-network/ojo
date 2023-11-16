<!-- markdownlint-disable MD041 -->
<!-- markdownlint-disable MD013 -->

![Logo!](assets/ojo.png)

[![Project Status: WIP – Initial development is in progress, but there has not yet been a stable, usable release suitable for the public.](https://www.repostatus.org/badges/latest/active.svg)](https://www.repostatus.org/badges/latest/active.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/ojo-network/ojo?style=flat-square)](https://goreportcard.com/report/github.com/ojo-network/ojo)
[![Version](https://img.shields.io/github/v/tag/ojo-network/ojo.svg?style=flat-square)](https://github.com/ojo-network/ojo/releases/latest)
[![License: Apache-2.0](https://img.shields.io/github/license/ojo-network/ojo.svg?style=flat-square)](https://github.com/ojo-network/ojo/blob/main/LICENSE)
[![GitHub Super-Linter](https://img.shields.io/github/actions/workflow/status/ojo-network/ojo/lint.yml?branch=main)](https://github.com/marketplace/actions/super-linter)

> A Golang Implementation of the Ojo Network, a decentralized oracle
> with DeFi safety in mind.

Ojo is an oracle platform which other blockchains and smart contracts can use to receive
up-to-date and accurate data. This platform arose from our work at
[Umee](https://github.com/umee-network/umee), where we worked on developing our
own oracle based off of the [Terra Classic](https://github.com/terra-money/classic-core) design.

Ojo is able to provide pricing info via IBC, CosmWasm, and EVM. MoveVM support is coming shortly.

## Table of Contents

- [Table of Contents](#table-of-contents)
- [Releases](#releases)
- [Install](#install)
- [Networks](#networks)

## Releases

Our releases are tagged and binaries are produced [here](https://github.com/ojo-network/ojo/releases).

See [Release procedure](contributing.md#release-procedure) for more information about the release model.

## Install

To install the `ojod` binary:

```shell
$ make install
```

## Networks

Ojo currently has three active public networks:

| Network Name                                      | Type              | Docs                                               |
| :-----------------------------------------------: | :---------------: | :------------------------------------------------: |
| [Agamotto](https://agamotto.ojo.network/agamotto) | Mainnet           | [Docs](https://docs.ojo.network/networks/agamotto) |
| [Ditto](https://agamotto.ojo.network/ditto)       | Testnet           | N/A                                                |
| [Sauron](https://sauron.ojo.network/)             | Validator Testnet | [Docs](https://docs.ojo.network/networks/sauron)   |
