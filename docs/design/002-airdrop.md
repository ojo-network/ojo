# Design Doc 002: Airdrop Module

## Status

Draft

## Abstract

For the Ojo Network's Airdrop, we intend to require that users are only able to claim the other half of their tokens by staking their entire first half with a validator.

## Context

We want to incentivize airdrop recipients to stake their tokens. An effective way to do this is by rewarding users with the other half of their airdrop allocation if they stake to a validator.

We also need to ensure that we're able to create vesting accounts for the initial airdrop amount, and controlling how long afterwards it will take to unlock. An issue here is that vesting accounts in cosmos can only be created for accounts which do not yet exist.

## Specification

### Msgs

- `CreateAirdropAccount(address, tokensToReceive, amountExpectedToDelegate, vestingLength)` - Create an airdrop account which will unlock into a vesting account, if all the tokens in amountExpectedToDelegate are delegated. This transaction can only occur at genesis.
- `ClaimAirdrop(fromAddress, toAddress)` - Allows an airdrop recipient to claim the 2nd portion of the airdrop specified in the `CreateAirdropAccount` message.
  - This transaction will create a new Delayed Vesting Account at `toAddress` with the amount of tokens in `tokensToReceive`. This account will vest as long as `vestingLength` above. This transaction fails if `amountExpectedToDelegate` is not met from the `fromAddress` account. Emits an event once the airdrop has been claimed.

### Proposed API

- `QueryAirdropAccount` - Returns an existing airdrop account, along with whether or not it has already been claimed.

### Outcomes

> What systems will be affected?

This is mostly agnostic, since we've decided to use existing vesting account solutions.

> Are there any logging, monitoring or observability needs?

Event emissions & regular error logging on the keeper is necessary.

> Are there any security considerations?

We want to make sure users are unable to:
* Claim one airdrop multiple times.
* Claim one airdrop to multiple accounts at once.
* Keep other users from blocking their ability to claim their airdrop.
* Authorize an airdrop claim for an account other than their own.
* Create airdrop accounts after genesis.

> Will these changes require a breaking (major) release?

Yes

> Does this change require coordination with other teams?

No

## Alternative Approaches

### Forking the vesting account code to add a conditional vesting account option

This would likely result in cosmos sdk issues that we would not forsee, and would take much more effort. It would also likely impact block scanners and confuse users.

### Implement our own form of vesting accounts

Load of effort here would be a lot more than necessary.

## Consequences

- Users need to have multiple accounts

### Backwards Compatibility

Yes

### Positive

- Less development time, simple implementation.
- No frontend changes necessary from block scanners.
- Conditional vesting accounts without the need to send liquid tokens.

### Negative

- Multiple vesting accounts needed from the user.
- Additional KV store space.

### Neutral

N/A

## Further Discussions

- Is there a way to implement this without two user accounts?

## References

- [Genesis Vesting Account](https://docs.cosmos.network/v0.45/modules/auth/05_vesting.html#genesis-initialization)
