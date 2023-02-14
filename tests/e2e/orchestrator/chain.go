package orchestrator

type Chain struct {
	chainId  string
	mnemonic string
}

func NewChain(chainId string) *Chain {
	mnemonic, _ := createMnemonic()
	return &Chain{
		chainId:  chainId,
		mnemonic: mnemonic,
	}
}
