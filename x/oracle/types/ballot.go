package types

import (
	"bytes"
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ojo-network/ojo/util/decmath"
)

// VoteForTally is a convenience wrapper to reduce redundant lookup cost.
type VoteForTally struct {
	Denom        string
	ExchangeRate sdk.Dec
	Voter        sdk.ValAddress
	Power        int64
}

// NewVoteForTally returns a new VoteForTally instance.
func NewVoteForTally(rate sdk.Dec, denom string, voter sdk.ValAddress, power int64) VoteForTally {
	return VoteForTally{
		ExchangeRate: rate,
		Denom:        denom,
		Voter:        voter,
		Power:        power,
	}
}

// ExchangeRateBallot is a convenience wrapper around a ExchangeRateVote slice.
type ExchangeRateBallot []VoteForTally

// ToMap return organized exchange rate map by validator.
func (pb ExchangeRateBallot) ToMap() map[string]sdk.Dec {
	exchangeRateMap := make(map[string]sdk.Dec)
	for _, vote := range pb {
		if vote.ExchangeRate.IsPositive() {
			exchangeRateMap[vote.Voter.String()] = vote.ExchangeRate
		}
	}

	return exchangeRateMap
}

// Power returns the total amount of voting power in the ballot.
func (pb ExchangeRateBallot) Power() int64 {
	var totalPower int64
	for _, vote := range pb {
		totalPower += vote.Power
	}

	return totalPower
}

// ExchangeRates returns the exchange rates in the ballot as a list.
func (pb ExchangeRateBallot) ExchangeRates() []sdk.Dec {
	exchangeRates := []sdk.Dec{}
	for _, vote := range pb {
		if vote.ExchangeRate.BigInt().BitLen() <= sdk.MaxBitLen {
			exchangeRates = append(exchangeRates, vote.ExchangeRate)
		}
	}

	return exchangeRates
}

// Median returns the median of the ExchangeRateVote.
func (pb ExchangeRateBallot) Median() (sdk.Dec, error) {
	rates := pb.ExchangeRates()
	if len(rates) == 0 {
		return sdk.ZeroDec(), nil
	}

	return decmath.Median(rates)
}

// StandardDeviation returns the standard deviation around the median of the
// ExchangeRateVote.
func (pb ExchangeRateBallot) StandardDeviation(median sdk.Dec) (sdk.Dec, error) {
	rates := pb.ExchangeRates()
	if len(rates) == 0 {
		return sdk.ZeroDec(), nil
	}

	return decmath.MedianDeviation(median, rates)
}

// Len implements sort.Interface
func (pb ExchangeRateBallot) Len() int {
	return len(pb)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (pb ExchangeRateBallot) Less(i, j int) bool {
	if pb[i].ExchangeRate.LT(pb[j].ExchangeRate) {
		return true
	}
	if pb[i].ExchangeRate.Equal(pb[j].ExchangeRate) {
		return bytes.Compare(pb[i].Voter, pb[j].Voter) < 0
	}
	return false
}

// Swap implements sort.Interface.
func (pb ExchangeRateBallot) Swap(i, j int) {
	pb[i], pb[j] = pb[j], pb[i]
}

// BallotDenom is a convenience wrapper for setting rates deterministically.
type BallotDenom struct {
	Ballot ExchangeRateBallot
	Denom  string
}

// BallotMapToSlice returns an array of sorted exchange rate ballots.
func BallotMapToSlice(votes map[string]ExchangeRateBallot) []BallotDenom {
	b := make([]BallotDenom, len(votes))
	i := 0
	for denom, ballot := range votes {
		b[i] = BallotDenom{
			Denom:  denom,
			Ballot: ballot,
		}
		i++
	}
	sort.Slice(b, func(i, j int) bool {
		return b[i].Denom < b[j].Denom
	})
	return b
}

// Claim is an interface that directs its rewards to an attached bank account.
type Claim struct {
	Power             int64
	MandatoryWinCount int64
	Recipient         sdk.ValAddress
}

// NewClaim generates a Claim instance.
func NewClaim(power, mandatoryWinCount int64, recipient sdk.ValAddress) Claim {
	return Claim{
		Power:             power,
		MandatoryWinCount: mandatoryWinCount,
		Recipient:         recipient,
	}
}

// ClaimMapToSlices returns an array of sorted exchange rate ballots and a second
// array with validators not in the includeMap filtered out.
func ClaimMapToSlices(claims map[string]Claim, includeMap map[string]bool) ([]Claim, []Claim) {
	c := make([]Claim, len(claims))
	r := make([]Claim, len(includeMap))
	i := 0
	j := 0
	for _, claim := range claims {
		if _, ok := includeMap[claim.Recipient.String()]; ok {
			r[j] = Claim{
				Power:             claim.Power,
				MandatoryWinCount: claim.MandatoryWinCount,
				Recipient:         claim.Recipient,
			}
			j++
		}
		c[i] = Claim{
			Power:             claim.Power,
			MandatoryWinCount: claim.MandatoryWinCount,
			Recipient:         claim.Recipient,
		}
		i++
	}
	sort.Slice(c, func(i, j int) bool {
		return c[i].Recipient.String() < c[j].Recipient.String()
	})
	sort.Slice(r, func(i, j int) bool {
		return r[i].Recipient.String() < r[j].Recipient.String()
	})
	return c, r
}
