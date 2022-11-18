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

In order to solve #1, the *Client* module will contain a simple API

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

### User requirements

> Specify expected user behavior

### Proposed API

> Describe new API or changed API

### Outcomes

> This section does not need to be filled in at the beginning, but must
> be completed prior to the merging of the implementation.
> Here are some common questions that get answered as part of the detailed design:
>
> - What systems will be affected?
> - Are there any logging, monitoring or observability needs?
> - Are there any security considerations?
> - Will these changes require a breaking (major) release?
> - Does this change require coordination with other teams?

## Alternative Approaches

> This section contains information around alternative options that are considered
> before making a decision. It should contain a explanation on why the alternative
> approach(es) were not chosen.

## Test Cases [optional]

> How will the changes be tested?
> Test cases in the form of example scenarios for an implementation are mandatory for designs that are affecting important parts of the system functionality. Design docs can include links to test cases if applicable.

## Consequences

> This section describes the consequences, after applying the decision. All
> consequences should be summarized here, not just the "positive" ones.

### Backwards Compatibility

> All design docs that introduce backwards incompatibilities must include a section describing these incompatibilities and their severity. The doc must explain how the author proposes to deal with these incompatibilities. Submissions without a sufficient backwards compatibility treatise may be rejected outright.

### Positive

### Negative

### Neutral

## Further Discussions

> This section should contain potential followups or issues to be solved in future iterations (usually referencing comments from a pull-request discussion).

## Comments

> Optional. Provide additional important comments.
> If the proposal is rejected, document why it was rejected.

## References

> Are there any relevant PR comments, issues that led up to this, or articles
> referenced for why we made the given design choice? If so link them here!

- {reference link}
