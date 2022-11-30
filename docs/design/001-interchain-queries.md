# Design Doc 001: Interchain Queries

## Status

Draft

## Abstract

In order for the Ojo network to effectively provide services to other chains, we need a method of transport at Layer 1. We must develop two mechanisms:

*Client* Module - accepts validator responses to requests made by the *Client*; requires a *Proof* for information to be accepted.
*Relayer* - allows relayers to move information between Ojo and the *Client*, and is responsible for submitting proofs.

The Client chain will request data, which the relayer will then pick up, query Ojo for the requested data, and then respond back to the original chain.

## Context

There are two main existing solutions :

### [Quicksilver's Interchain Queries](https://github.com/ingenuity-build/quicksilver/tree/main/x/interchainquery/keeper)

Quicksilver requires accurate information about other chains in order to utilize interchain accounts for their remote liquid staking solution.

This is a KV store solution which requires the *Client* and *Relayer* infrastructure, but does not require a *Host* module. Quicksilver also requires validators to prove that data is accurate before accepting it on-chain.

### [Strangelove's IBC Module](https://github.com/strangelove-ventures/ibc-go/tree/feature/icq_implementation/modules/apps/icq)

The ICQ IBC module fully leverages IBC standards, and requires the *Client*, *Relayer*, and *Host* pieces of infrastructure. It also requires a proof to the *Client* that the information from the *Host* is accurate.

There are a few issues with bringing both of these solutions in-house:

1. Ease of API - The quicksilver & strangelove implementations are for agnostic data, and in order to consume relayed packets, the *Client* chain would have to parse this. This should be a part of the *Client* module.
2. Economic Factor - Neither implementation allows for relayers to be paid for their services; in fact, in both implementation if the gas fee of a chain is above 0, it *costs* coins to run a relayer between both chains. The economic model in this case will be dealt with in a future design doc.

## Specification
### Client

#### Msgs

- `FulfillDataRequest(data, proof)` - Fulfills the data request submitted on the `Host` module. Requires `Proof` that the information from the `Host` actually belongs to the host using an [ABCI RequestQuery.](https://github.com/strangelove-ventures/ibc-go/tree/feature/icq_implementation/modules/apps/icq#abci-query).

### Relayer

#### Responsibility

The relayer will loop through:

1. Querying for outstanding Data Requests on the *Host* Chain.
2. Querying for the data requested.
3. Submitting a `FulfillDataRequest` on the *Client* chain.

### Proposed API

#### Client Module

The client module should have an easy-to-use API specific to the type of data being relayed. At this level, data should not be generic. An example of the Keeper APIs we could implement here are:

- `GetPrice(denom) sdk.Dec` - returns relayed price of the asset `denom`.
- `GetReservesProof(denom) sdk.Dec` - returns an `sdk.Dec` in the range of `[0, 1]`, determining how valid a set of reserves are.

### Outcomes

> What systems will be affected?

This will affect validator requirements, on-chain and off-chain.

> Are there any logging, monitoring or observability needs?

We will want to strictly consider logging for the `Relayer` tool.

> Are there any security considerations?

Currently, we need to be concerned about bad actors not being able to:

1. Relay **bad** information from `Ojo` to the `Client` chain.
2. Request information that is unavailable.

In the future, we need to ensure that:

1. Relayers get paid for services rendered

> Will these changes require a breaking (major) release?

Yes

> Does this change require coordination with other teams?

Mainly the validator set.

## Alternative Approaches

There are existing ICQ methods that do not allow for:

1. An easy-to-use API.
2. A reward system for relayers.

Because of these points, we need this altered version.

## Consequences

### Backwards Compatibility

N/A. This is a new set of software components.

### Positive

* Future economic factor
* Allows for secure information
* Ease of implementation on the *Client* side

### Negative

* Development time & support
* ABCI blocktime impact

### Neutral

* Does not fully leverage IBC standards
* Requires use of ABCI Queries & IBC to get the host chain's header

## Further Discussions

* How do we make relayers profitable?

## References


- [Strangelove's ICQ](https://github.com/strangelove-ventures/ibc-go/tree/feature/icq_implementation/modules/apps/icq)
- [Quicksilver's ICQ](https://github.com/ingenuity-build/quicksilver)
