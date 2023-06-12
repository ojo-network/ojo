<!-- markdownlint-disable MD013 -->
<!-- markdownlint-disable MD024 -->

<!--
Changelog Guiding Principles:

Changelogs are for humans, not machines.
There should be an entry for every single version.
The same types of changes should be grouped.
Versions and sections should be linkable.
The latest version comes first.
The release date of each version is displayed.
Mention whether you follow Semantic Versioning.

Usage:

Change log entries are to be added to the Unreleased section under the
appropriate stanza (see below). Each entry should ideally include a tag and
the Github PR referenced in the following format:

* (<tag>) [#<PR-number>](https://github.com/ojo-network/ojo/pull/<PR-number>) <changelog entry>

Types of changes (Stanzas):

State Machine Breaking: for any changes that result in a divergent application state.
Features: for new features.
Improvements: for changes in existing functionality.
Deprecated: for soon-to-be removed features.
Bug Fixes: for any bug fixes.
Client Breaking: for breaking Protobuf, CLI, gRPC and REST routes used by clients.
API Breaking: for breaking exported Go APIs used by developers.

To release a new version, ensure an appropriate release branch exists. Add a
release version and date to the existing Unreleased section which takes the form
of:

## [<version>](https://github.com/ojo-network/ojo/releases/tag/<version>) - YYYY-MM-DD

Once the version is tagged and released, a PR should be made against the main
branch to incorporate the new changelog updates.

Ref: https://keepachangelog.com/en/1.0.0/
-->

# Changelog

## [Unreleased]

### Features

- [172](https://github.com/ojo-network/ojo/pull/172) Initial airdrop module ABCI & msg_server

## [v0.1.4]

### State Machine Breaking

- [203](https://github.com/ojo-network/ojo/pull/203) Update cosmos SDK to v0.46.13 - [barberry](https://forum.cosmos.network/t/cosmos-sdk-security-advisory-barberry/10825)

### Improvements

- [182](https://github.com/ojo-network/ojo/pull/182) Set standard for oracle parameter symbols

### Fixes

- [197](https://github.com/ojo-network/ojo/pull/197) Fix potential win count calculations.
- [202](https://github.com/ojo-network/ojo/pull/202) Migration to update ValidatorRewardSet from map to list.
- [207](https://github.com/ojo-network/ojo/pull/207) Store validator reward set in one key instead of by block

## [v0.1.3](https://github.com/ojo-network/ojo/releases/tag/v0.1.3) - 2023-05-09

### Features

- Initial release for Agamotto!
