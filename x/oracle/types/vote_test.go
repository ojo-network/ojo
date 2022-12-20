package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	appparams "github.com/ojo-network/ojo/app/params"
	"github.com/ojo-network/ojo/x/oracle/types"
)

func init() {
	appparams.SetAddressPrefixes()
}

func TestAggregateExchangeRatePrevoteString(t *testing.T) {
	addr := sdk.ValAddress(sdk.AccAddress([]byte("addr1_______________")))
	aggregateVoteHash := types.GetAggregateVoteHash("salt", "OJO:100,ATOM:100", addr)
	aggregateExchangeRatePreVote := types.NewAggregateExchangeRatePrevote(
		aggregateVoteHash,
		addr,
		100,
	)

	require.Equal(t, "hash: ccd44c2be8cec771f4bc8a0b33895bd44e3459b9\nvoter: ojovaloper1v9jxgu33ta047h6lta047h6lta047h6ludnc0y\nsubmit_block: 100\n", aggregateExchangeRatePreVote.String())
}

func TestAggregateExchangeRateVoteString(t *testing.T) {
	aggregateExchangeRatePreVote := types.NewAggregateExchangeRateVote(
		types.ExchangeRateTuples{
			types.NewExchangeRateTuple(types.OjoDenom, sdk.OneDec()),
		},
		sdk.ValAddress(sdk.AccAddress([]byte("addr1_______________"))),
	)

	require.Equal(t, "exchange_rate_tuples:\n    - denom: uojo\n      exchange_rate: \"1.000000000000000000\"\nvoter: ojovaloper1v9jxgu33ta047h6lta047h6lta047h6ludnc0y\n", aggregateExchangeRatePreVote.String())
}

func TestExchangeRateTuplesString(t *testing.T) {
	exchangeRateTuple := types.NewExchangeRateTuple(types.OjoDenom, sdk.OneDec())
	require.Equal(t, exchangeRateTuple.String(), "denom: uojo\nexchange_rate: \"1.000000000000000000\"\n")

	exchangeRateTuples := types.ExchangeRateTuples{
		exchangeRateTuple,
		types.NewExchangeRateTuple(types.IbcDenomAtom, sdk.SmallestDec()),
	}
	require.Equal(t, "- denom: uojo\n  exchange_rate: \"1.000000000000000000\"\n- denom: ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2\n  exchange_rate: \"0.000000000000000001\"\n", exchangeRateTuples.String())
}

func TestParseExchangeRateTuples(t *testing.T) {
	valid := "uojo:123.0,uatom:123.123"
	_, err := types.ParseExchangeRateTuples(valid)
	require.NoError(t, err)

	duplicatedDenom := "uojo:100.0,uatom:123.123,uatom:121233.123"
	_, err = types.ParseExchangeRateTuples(duplicatedDenom)
	require.Error(t, err)

	invalidCoins := "123.123"
	_, err = types.ParseExchangeRateTuples(invalidCoins)
	require.Error(t, err)

	invalidCoinsWithValid := "uojo:123.0,123.1"
	_, err = types.ParseExchangeRateTuples(invalidCoinsWithValid)
	require.Error(t, err)

	zeroCoinsWithValid := "uojo:0.0,uatom:123.1"
	_, err = types.ParseExchangeRateTuples(zeroCoinsWithValid)
	require.Error(t, err)

	negativeCoinsWithValid := "uojo:-1234.5,uatom:123.1"
	_, err = types.ParseExchangeRateTuples(negativeCoinsWithValid)
	require.Error(t, err)

	multiplePricesPerRate := "uojo:123: uojo:456,uusdc:789"
	_, err = types.ParseExchangeRateTuples(multiplePricesPerRate)
	require.Error(t, err)

	res, err := types.ParseExchangeRateTuples("")
	require.Nil(t, err)
	require.Nil(t, res)
}
