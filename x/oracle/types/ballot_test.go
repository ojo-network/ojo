package types

import (
	"fmt"
	"math"
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
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

func TestPBMedian(t *testing.T) {
	tests := []struct {
		inputs      []int64
		powers      []int64
		isValidator []bool
		median      sdk.Dec
	}{
		{
			// Supermajority one number
			[]int64{1, 2, 10, 100000},
			[]int64{1, 1, 100, 1},
			[]bool{true, true, true, true},
			sdk.NewDec(6),
		},
		{
			// Adding fake validator doesn't change outcome
			[]int64{1, 2, 10, 100000, 10000000000},
			[]int64{1, 1, 100, 1, 10000},
			[]bool{true, true, true, true, false},
			sdk.NewDec(10),
		},
		{
			[]int64{1, 2, 3, 4},
			[]int64{1, 100, 100, 1},
			[]bool{true, true, true, true},
			sdk.MustNewDecFromStr("2.5"),
		},
		{
			// No votes
			[]int64{},
			[]int64{},
			[]bool{true, true, true, true},
			sdk.NewDec(0),
		},
	}

	for _, tc := range tests {
		pb := ExchangeRateBallot{}
		for i, input := range tc.inputs {
			valAddr := sdk.ValAddress(secp256k1.GenPrivKey().PubKey().Address())

			power := tc.powers[i]
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

		median, err := pb.Median()
		require.NoError(t, err)
		require.Equal(t, tc.median, median)
	}
}

func TestPBStandardDeviation(t *testing.T) {
	tests := []struct {
		inputs            []sdk.Dec
		powers            []int64
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
			sdk.MustNewDecFromStr("49997.000142508550309932"),
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
			[]sdk.Dec{
				sdk.MustNewDecFromStr("1.0"),
				sdk.MustNewDecFromStr("2.0"),
				sdk.MustNewDecFromStr("3.0"),
				sdk.MustNewDecFromStr("4.00"),
			},
			[]int64{1, 100, 100, 1},
			[]bool{true, true, true, true},
			sdk.MustNewDecFromStr("1.118033988749894848"),
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

			power := tc.powers[i]
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
		median, err := pb.Median()
		require.NoError(t, err)
		stdDev, _ := pb.StandardDeviation(median)
		require.NoError(t, err)
		require.Equal(t, tc.standardDeviation, stdDev)
	}
}

func TestPBStandardDeviation_Overflow(t *testing.T) {
	valAddr := sdk.ValAddress(secp256k1.GenPrivKey().PubKey().Address())
	overflowRate, err := sdk.NewDecFromStr("1000000000000000000000000000000000000000000000000000000000000.0")
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
	median, err := pb.Median()
	require.NoError(t, err)
	deviation, err := pb.StandardDeviation(median)
	require.NoError(t, err)
	expectedDevation := sdk.MustNewDecFromStr("616.5")
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

func TestClaimMapToSlice(t *testing.T) {
	valAddress := GenerateRandomValAddr(1)
	claim := NewClaim(10, 4, valAddress[0])
	claimSlice := ClaimMapToSlice(map[string]Claim{
		"testClaim":    claim,
		"anotherClaim": claim,
	})
	require.Equal(t, []Claim{claim, claim}, claimSlice)
}
