package orchestrator

import (
	"encoding/json"
	"fmt"
	"os"

	cmtrand "github.com/cometbft/cometbft/libs/rand"

	"github.com/cosmos/cosmos-sdk/codec"
	//"github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/ojo-network/ojo/client/tx"
	appparams "github.com/ojo-network/ojo/app/params"
)

const (
	keyringPassphrase = "testpassphrase"
	keyringAppName    = "testnet"
)

var encodingConfig appparams.EncodingConfig

type chain struct {
	cdc        codec.Codec
	dataDir    string
	id         string
	validators []*validator
	accounts   []*tx.OjoAccount
}

func newChain(cdc codec.Codec) (*chain, error) {
	tmpDir, err := os.MkdirTemp("", "ojo-e2e-testnet-")
	if err != nil {
		return nil, err
	}

	return &chain{
		id:      "chain-" + cmtrand.NewRand().Str(6),
		dataDir: tmpDir,
		cdc:     cdc,
	}, nil
}

func (c *chain) configDir() string {
	return fmt.Sprintf("%s/%s", c.dataDir, c.id)
}

func (c *chain) createAndInitValidators(count int, gen map[string]json.RawMessage) error {
	for i := 0; i < count; i++ {
		node := c.createValidator(i)

		// generate genesis files
		if err := node.init(gen); err != nil {
			return err
		}

		c.validators = append(c.validators, node)

		// create keys
		if err := node.createKey("val"); err != nil {
			return err
		}
		if err := node.createNodeKey(); err != nil {
			return err
		}
		if err := node.createConsensusKey(); err != nil {
			return err
		}
	}

	return nil
}

func (c *chain) createAccounts(numAccounts int) error {
	for i := 0; i < numAccounts; i++ {
		newAccount, err := tx.NewOjoAccount(fmt.Sprintf("account-%d", i))
		if err != nil {
			return err
		}
		c.accounts = append(c.accounts, newAccount)
	}
	return nil
}

func (c *chain) createValidator(index int) *validator {
	return &validator{
		chain:   c,
		index:   index,
		moniker: "ojo",
	}
}
