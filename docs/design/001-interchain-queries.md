# Design Doc 001: Interchain Queries

## Status

Draft

## Abstract

In order for the Ojo network to effectively provide services to other chains, we need a method of transport at Layer 1. We must develop three mechanisms: a *Host* module on Ojo; a *Client* module for consumer chains to use; and a *Relayer* to move messages in between the first two. The consumer chain will request data, which the relayer will then pick up, query Ojo for, and then respond back to the original chain.

*Host* Module - requires relayers to "prove" that they have relayed information accurately in order to gain rewards. Clients will also be able to pay for services rendered by relayers.

*Client* Module - accepts validator responses to requests made by the *Client*; requires a *Proof* for information to be accepted. Also is responsible for ensuring the client has paid.

*Relayer* - allows relayers to move information between the *Host* and *Client*, and is responsible for proving that data was accurate to be paid. In this implementation, validators also act as relayers.

## Context

There are two main existing solutions :

### [Quicksilver's Interchain Queries](https://github.com/ingenuity-build/quicksilver/tree/main/x/interchainquery/keeper)

Quicksilver requires accurate about other chains in order to utilize interchain accounts for their remote liquid staking solution.

This is a KV store solution which requires the *Client* and *Relayer* infrastructure, but does not require a *Host* module. Quicksilver also requires validators to prove that data is accurate before accepting it on-chain.

### [Strangelove's IBC Module](https://github.com/strangelove-ventures/ibc-go/tree/feature/icq_implementation/modules/apps/icq)

The ICQ IBC module fully leverages IBC standards, and requires the *Client*, *Relayer*, and *Host* pieces of infrastructure. It also requires a proof to the *Client* that the information from the *Host* is accurate.

There are a few issues with bringing both of these solutions in-house:

1. Ease of API - The quicksilver & strangelove implementations are for agnostic data, and in order to consume relayed packets, the *Client* chain would have to parse this. This should be a part of the *Client* module.
2. Economic Factor - Neither implementation allows for relayers to be paid for their services; in fact, in both implementation if the gas fee of a chain is above 0, it *costs* coins to run a relayer between both chains.

## Specification

### Host Module

#### Encryption

In order for on-chain data to not be intercepted and relayed, any packets which are queryable must be stored in an encrypted fashion using the public keys of all the active-set validators. This data will be decrypted off-chain by the `Relayer` component using the private key of that validator, then relayed.

#### Msgs

- `SubmitDataRequest(dataType, chainID)` - Registers a requested dataset to be relayed to `chainID`.

- `VerifyRelay(requestID, proof)` - Allows validators to claim rewards for relaying information.

#### Queries

- `QueryRequestedData(requestID)` - Returns an encrypted form of the data requested by `SubmitDataRequest`.

- `QueryDataRequests()` - Returns a set of active `requestIDs`.

### Client Module

#### Msgs

- `FulfillDataRequest(decryptedData, proof)` - Fulfills the data request submitted on the `Host` module. Requires `Proof` that the information from the `Host` actually belongs to the host using an [ABCI RequestQuery.](https://github.com/strangelove-ventures/ibc-go/tree/feature/icq_implementation/modules/apps/icq#abci-query).

### Relayer

#### API

-  `DecryptResponsePacket(requestID, encryptedData)` - Allows only active validators to decrypt a given packet of data with a secret key. This is so that non-validators are unable to relay.

#### Responsibility

The relayer will loop through:

1. Querying for outstanding Data Requests on the *Host* Chain.
2. Querying for the data requested.
3. Decrypting the data.
4. Submitting a `FulfillDataRequest` on the *Client* chain.
5. Submitting a `VerifyRelay` on the *Host* chain.

### Proposed API

#### Client Module

The client module should have an easy-to-use API specific to the type of data being relayed. At this level, data should not be generic. An example of the Keeper APIs we could implement here are:

- `GetPrice(denom) sdk.Dec` - returns relayed price of the asset `denom`.
- `GetReservesProof(denom) sdk.Dec` - returns an `sdk.Dec` in the range of `[0, 1]`, determining how valid a set of reserves are.

#### Host Module

The host module should have APIs that other modules can use to store datasets. Other modules should be concerned with aggregating data such as `pricing`, `proof of reserves`, and then submitting that information to be stored by `x/icq`.

These keepers should provide an argnostic API, which would look like:

`SetDatapoint(key, value) error`

### Outcomes

> What systems will be affected?

This will affect validator requirements, on-chain and off-chain.

> Are there any logging, monitoring or observability needs?

We will want to strictly consider logging for the `Relayer` tool.

> Are there any security considerations?

We need to ensure bad actors will not:

1. Relay information without being a validator on `Ojo`.
2. Relay **bad** information from `Ojo` to the `Client` chain.
3. Relay information between these modules without submitting a `SubmitDataRequest` msg.
4. Abuse `SubmitDataRequest` to send malicious information to the receiving chain.
5. Request information that is unavailable.

> Will these changes require a breaking (major) release?

Yes

> Does this change require coordination with other teams?

Mainly the validator set.

## Alternative Approaches

As mentioned in [Context](##Context), there are existing ICQ methods that do not allow for:

1. An easy-to-use API.
2. A reward system for relayers.

Because of these points, we need this altered version.

## Consequences

### Backwards Compatibility

N/A. This is a new set of modules.

### Positive

* Economic factor
* Allows for secure information
* Ease of implementation on the *Client* side
* Not dependent on IBC versions

### Negative

* Development time & support
* ABCI blocktime impact

### Neutral

* Does not fully leverage IBC standards

## Further Discussions

* Could relayers not have to be validators, using some sort of ante system?


## References


- [Strangelove's ICQ](https://github.com/strangelove-ventures/ibc-go/tree/feature/icq_implementation/modules/apps/icq)
- [Quicksilver's ICQ](https://github.com/ingenuity-build/quicksilver)
