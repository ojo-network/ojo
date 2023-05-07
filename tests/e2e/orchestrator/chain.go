package orchestrator

import (
	"fmt"
	"os"

	"github.com/cosmos/cosmos-sdk/codec"
	tmrand "github.com/tendermint/tendermint/libs/rand"

	ojoapp "github.com/ojo-network/ojo/app"
	"github.com/ojo-network/ojo/app/params"
	"github.com/ojo-network/ojo/client/tx"
)

const (
	keyringPassphrase = "testpassphrase"
	keyringAppName    = "testnet"
)

var (
	encodingConfig params.EncodingConfig
	cdc            codec.Codec
)

func init() {
	encodingConfig = ojoapp.MakeEncodingConfig()
	cdc = encodingConfig.Codec
}

type chain struct {
	dataDir    string
	id         string
	validators []*validator
	accounts   []*tx.OjoAccount
}

func newChain() (*chain, error) {
	tmpDir, err := os.MkdirTemp("", "ojo-e2e-testnet-")
	if err != nil {
		return nil, err
	}

	return &chain{
		id:      "chain-" + tmrand.NewRand().Str(6),
		dataDir: tmpDir,
	}, nil
}

func (c *chain) configDir() string {
	return fmt.Sprintf("%s/%s", c.dataDir, c.id)
}

func (c *chain) createAndInitValidators(count int) error {
	for i := 0; i < count; i++ {
		node := c.createValidator(i)

		// generate genesis files
		if err := node.init(); err != nil {
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
