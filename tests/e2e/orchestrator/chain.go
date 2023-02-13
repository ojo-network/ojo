package orchestrator

type Chain struct {
	chainId        string
	val01_mnemonic string
}

func NewChain(chainId string) *Chain {
	mnemonic, _ := createMnemonic()
	return &Chain{
		chainId:        chainId,
		val01_mnemonic: mnemonic,
	}
}
