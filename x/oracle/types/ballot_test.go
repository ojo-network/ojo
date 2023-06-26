package types

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"testing"

	"github.com/cometbft/cometbft/crypto/secp256k1"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestToMap(t *testing.T) {
	tests := struct {
		votes   []VoteForTally
		isValid []bool
	}{
		[]VoteForTally{
			{
				Voter:        sdk.ValAddress(secp256k1.GenPrivKey().PubKey().Address()),
				Denom:        OjoDenom,
				ExchangeRate: sdk.NewDec(1600),
				Power:        100,
			},
			{
				Voter:        sdk.ValAddress(secp256k1.GenPrivKey().PubKey().Address()),
				Denom:        OjoDenom,
				ExchangeRate: sdk.ZeroDec(),
				Power:        100,
			},
			{
				Voter:        sdk.ValAddress(secp256k1.GenPrivKey().PubKey().Address()),
				Denom:        OjoDenom,
				ExchangeRate: sdk.NewDec(1500),
				Power:        100,
			},
		},
		[]bool{true, false, true},
	}

	pb := ExchangeRateBallot(tests.votes)
	mapData := pb.ToMap()

	for i, vote := range tests.votes {
		exchangeRate, ok := mapData[vote.Voter.String()]
		if tests.isValid[i] {
			require.True(t, ok)
			require.Equal(t, exchangeRate, vote.ExchangeRate)
		} else {
			require.False(t, ok)
		}
	}
}

func TestSqrt(t *testing.T) {
	num := sdk.NewDecWithPrec(144, 4)
	floatNum, err := strconv.ParseFloat(num.String(), 64)
	require.NoError(t, err)

	floatNum = math.Sqrt(floatNum)
	num, err = sdk.NewDecFromStr(fmt.Sprintf("%f", floatNum))
	require.NoError(t, err)

	require.Equal(t, sdk.NewDecWithPrec(12, 2), num)
}

func TestPBPower(t *testing.T) {
	ctx := sdk.NewContext(nil, tmproto.Header{}, false, nil)
	valAccAddrs, sk := GenerateRandomTestCase()
	pb := ExchangeRateBallot{}
	ballotPower := int64(0)

	for i := 0; i < len(sk.Validators()); i++ {
		power := sk.Validator(ctx, valAccAddrs[i]).GetConsensusPower(sdk.DefaultPowerReduction)
		vote := NewVoteForTally(
			sdk.ZeroDec(),
			OjoDenom,
			valAccAddrs[i],
			power,
		)

		pb = append(pb, vote)
		require.NotEqual(t, int64(0), vote.Power)

		ballotPower += vote.Power
	}

	require.Equal(t, ballotPower, pb.Power())

	// Mix in a fake validator, the total power should not have changed.
	pubKey := secp256k1.GenPrivKey().PubKey()
	faceValAddr := sdk.ValAddress(pubKey.Address())
	fakeVote := NewVoteForTally(
		sdk.OneDec(),
		OjoDenom,
		faceValAddr,
		0,
	)

	pb = append(pb, fakeVote)
	require.Equal(t, ballotPower, pb.Power())
}

func TestPBWeightedMedian(t *testing.T) {
	tests := []struct {
		inputs      []int64
		weights     []int64
		isValidator []bool
		median      sdk.Dec
		errMsg      string
	}{
		{
			// Supermajority one number
			[]int64{1, 2, 10, 100000},
			[]int64{1, 1, 100, 1},
			[]bool{true, true, true, true},
			sdk.NewDec(10),
			"",
		},
		{
			// Adding fake validator doesn't change outcome
			[]int64{1, 2, 10, 100000, 10000000000},
			[]int64{1, 1, 100, 1, 10000},
			[]bool{true, true, true, true, false},
			sdk.NewDec(10),
			"",
		},
		{
			// Tie votes
			[]int64{1, 2, 3, 4},
			[]int64{1, 100, 100, 1},
			[]bool{true, true, true, true},
			sdk.NewDec(2),
			"",
		},
		{
			// No votes
			[]int64{},
			[]int64{},
			[]bool{true, true, true, true},
			sdk.NewDec(0),
			"",
		},
		{
			// Out of order
			[]int64{1, 2, 10, 3},
			[]int64{1, 1, 100, 1},
			[]bool{true, true, true, true},
			sdk.NewDec(10),
			"ballot must be sorted before this operation",
		},
	}

	for _, tc := range tests {
		pb := ExchangeRateBallot{}
		for i, input := range tc.inputs {
			valAddr := sdk.ValAddress(secp256k1.GenPrivKey().PubKey().Address())

			power := tc.weights[i]
			if !tc.isValidator[i] {
				power = 0
			}

			vote := NewVoteForTally(
				sdk.NewDec(int64(input)),
				OjoDenom,
				valAddr,
				power,
			)

			pb = append(pb, vote)
		}

		median, err := pb.WeightedMedian()
		if tc.errMsg == "" {
			require.NoError(t, err)
			require.Equal(t, tc.median, median)
		} else {
			require.ErrorContains(t, err, tc.errMsg)
		}
	}
}

func TestPBStandardDeviation(t *testing.T) {
	tests := []struct {
		inputs            []sdk.Dec
		weights           []int64
		isValidator       []bool
		standardDeviation sdk.Dec
	}{
		{
			// Supermajority one number
			[]sdk.Dec{
				sdk.MustNewDecFromStr("1.0"),
				sdk.MustNewDecFromStr("2.0"),
				sdk.MustNewDecFromStr("10.0"),
				sdk.MustNewDecFromStr("100000.00"),
			},
			[]int64{1, 1, 100, 1},
			[]bool{true, true, true, true},
			sdk.MustNewDecFromStr("49995.000362536252310906"),
		},
		{
			// Adding fake validator doesn't change outcome
			[]sdk.Dec{
				sdk.MustNewDecFromStr("1.0"),
				sdk.MustNewDecFromStr("2.0"),
				sdk.MustNewDecFromStr("10.0"),
				sdk.MustNewDecFromStr("100000.00"),
				sdk.MustNewDecFromStr("10000000000"),
			},
			[]int64{1, 1, 100, 1, 10000},
			[]bool{true, true, true, true, false},
			sdk.MustNewDecFromStr("4472135950.751005519905537611"),
		},
		{
			// Tie votes
			[]sdk.Dec{
				sdk.MustNewDecFromStr("1.0"),
				sdk.MustNewDecFromStr("2.0"),
				sdk.MustNewDecFromStr("3.0"),
				sdk.MustNewDecFromStr("4.00"),
			},
			[]int64{1, 100, 100, 1},
			[]bool{true, true, true, true},
			sdk.MustNewDecFromStr("1.224744871391589049"),
		},
		{
			// No votes
			[]sdk.Dec{},
			[]int64{},
			[]bool{true, true, true, true},
			sdk.NewDecWithPrec(0, 0),
		},
	}

	for _, tc := range tests {
		pb := ExchangeRateBallot{}
		for i, input := range tc.inputs {
			valAddr := sdk.ValAddress(secp256k1.GenPrivKey().PubKey().Address())

			power := tc.weights[i]
			if !tc.isValidator[i] {
				power = 0
			}

			vote := NewVoteForTally(
				input,
				OjoDenom,
				valAddr,
				power,
			)

			pb = append(pb, vote)
		}
		stdDev, _ := pb.StandardDeviation()

		require.Equal(t, tc.standardDeviation, stdDev)
	}
}

func TestPBStandardDeviation_Overflow(t *testing.T) {
	valAddr := sdk.ValAddress(secp256k1.GenPrivKey().PubKey().Address())
	overflowRate, err := sdk.NewDecFromStr("100000000000000000000000000000000000000000000000000000000.0")
	require.NoError(t, err)
	pb := ExchangeRateBallot{
		NewVoteForTally(
			sdk.OneDec(),
			OjoSymbol,
			valAddr,
			2,
		),
		NewVoteForTally(
			sdk.NewDec(1234),
			OjoSymbol,
			valAddr,
			2,
		),
		NewVoteForTally(
			overflowRate,
			OjoSymbol,
			valAddr,
			1,
		),
	}

	deviation, err := pb.StandardDeviation()
	require.NoError(t, err)
	expectedDevation := sdk.MustNewDecFromStr("871.862661203013097586")
	require.Equal(t, expectedDevation, deviation)
}

func TestBallotMapToSlice(t *testing.T) {
	valAddress := GenerateRandomValAddr(1)

	pb := ExchangeRateBallot{
		NewVoteForTally(
			sdk.NewDec(1234),
			OjoSymbol,
			valAddress[0],
			2,
		),
		NewVoteForTally(
			sdk.NewDec(12345),
			OjoSymbol,
			valAddress[0],
			1,
		),
	}

	ballotSlice := BallotMapToSlice(map[string]ExchangeRateBallot{
		OjoDenom:     pb,
		IbcDenomAtom: pb,
	})
	require.Equal(t, []BallotDenom{{Ballot: pb, Denom: IbcDenomAtom}, {Ballot: pb, Denom: OjoDenom}}, ballotSlice)
}

func TestExchangeRateBallotSwap(t *testing.T) {
	valAddress := GenerateRandomValAddr(2)

	voteTallies := []VoteForTally{
		NewVoteForTally(
			sdk.NewDec(1234),
			OjoSymbol,
			valAddress[0],
			2,
		),
		NewVoteForTally(
			sdk.NewDec(12345),
			OjoSymbol,
			valAddress[1],
			1,
		),
	}

	pb := ExchangeRateBallot{voteTallies[0], voteTallies[1]}

	require.Equal(t, pb[0], voteTallies[0])
	require.Equal(t, pb[1], voteTallies[1])
	pb.Swap(1, 0)
	require.Equal(t, pb[1], voteTallies[0])
	require.Equal(t, pb[0], voteTallies[1])
}

func TestStandardDeviationUnsorted(t *testing.T) {
	valAddress := GenerateRandomValAddr(1)
	pb := ExchangeRateBallot{
		NewVoteForTally(
			sdk.NewDec(1234),
			OjoSymbol,
			valAddress[0],
			2,
		),
		NewVoteForTally(
			sdk.NewDec(12),
			OjoSymbol,
			valAddress[0],
			1,
		),
	}

	deviation, err := pb.StandardDeviation()
	require.ErrorIs(t, err, ErrBallotNotSorted)
	require.Equal(t, "0.000000000000000000", deviation.String())
}

func TestClaimMapToSlices(t *testing.T) {
	valAddresses := GenerateRandomValAddr(2)
	claim1 := NewClaim(10, 1, 4, valAddresses[0])
	claim2 := NewClaim(10, 1, 4, valAddresses[1])
	claimSlice, rewardSlice := ClaimMapToSlices(
		map[string]Claim{
			"testClaim":    claim1,
			"anotherClaim": claim2,
		},
		[]string{
			claim1.Recipient.String(),
		},
	)
	require.Contains(t, claimSlice, claim1, claim2)
	require.Equal(t, []Claim{claim1}, rewardSlice)
}

func TestExchangeRateBallotSort(t *testing.T) {
	v1 := VoteForTally{ExchangeRate: sdk.MustNewDecFromStr("0.2"), Voter: sdk.ValAddress{0, 1}}
	v1Cpy := VoteForTally{ExchangeRate: sdk.MustNewDecFromStr("0.2"), Voter: sdk.ValAddress{0, 1}}
	v2 := VoteForTally{ExchangeRate: sdk.MustNewDecFromStr("0.1"), Voter: sdk.ValAddress{0, 1, 1}}
	v3 := VoteForTally{ExchangeRate: sdk.MustNewDecFromStr("0.1"), Voter: sdk.ValAddress{0, 1}}
	v4 := VoteForTally{ExchangeRate: sdk.MustNewDecFromStr("0.5"), Voter: sdk.ValAddress{1}}

	tcs := []struct {
		got      ExchangeRateBallot
		expected ExchangeRateBallot
	}{
		{
			got:      ExchangeRateBallot{v1, v2, v3, v4},
			expected: ExchangeRateBallot{v3, v2, v1, v4},
		},
		{
			got:      ExchangeRateBallot{v1},
			expected: ExchangeRateBallot{v1},
		},
		{
			got:      ExchangeRateBallot{v1, v1Cpy},
			expected: ExchangeRateBallot{v1, v1Cpy},
		},
	}
	for i, tc := range tcs {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			sort.Sort(tc.got)
			require.Exactly(t, tc.expected, tc.got)
		})
	}
}
