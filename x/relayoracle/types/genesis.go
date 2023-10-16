package types

import (
	"fmt"

	host "github.com/cosmos/ibc-go/v7/modules/core/24-host"
)

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		PortId: PortID,
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	if err := host.PortIdentifierValidator(gs.PortId); err != nil {
		return err
	}

	err := gs.Params.Validate()
	if err != nil {
		return err
	}

	requestIds := make(map[uint64]struct{})
	for _, request := range gs.Requests {
		if _, found := requestIds[request.RequestID]; found {
			return fmt.Errorf("duplicated request id: %d", request.RequestID)
		}

		requestIds[request.RequestID] = struct{}{}
	}

	// check if all pending requests and results have a valid request object
	for _, pending := range gs.GetPendingRequestIds() {
		if _, found := requestIds[pending]; !found {
			return fmt.Errorf("pending request id not found: %d", pending)
		}
	}

	for _, result := range gs.GetResults() {
		if _, found := requestIds[result.RequestID]; !found {
			return fmt.Errorf("result request id not found: %d", result.RequestID)
		}
	}

	return nil
}
