package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (p ParamUpdatePlan) String() string {
	due := p.DueAt()
	return fmt.Sprintf(`Oracle Param Update Plan
  Title: %s
  %s
  Description: %s.`, p.Title, due, p.Description)
}

// ValidateBasic does basic validation of a ParamUpdatePlan
func (p ParamUpdatePlan) ValidateBasic() error {
	if len(p.Title) == 0 {
		return ErrInvalidRequest.Wrap("name cannot be empty")
	}
	if p.Height <= 0 {
		return ErrInvalidRequest.Wrap("height must be greater than 0")
	}

	return nil
}

// ShouldExecute returns true if the Plan is ready to execute given the current context
func (p ParamUpdatePlan) ShouldExecute(ctx sdk.Context) bool {
	return p.Height == ctx.BlockHeight()
}

// DueAt is a string representation of when this plan is due to be executed
func (p ParamUpdatePlan) DueAt() string {
	return fmt.Sprintf("height: %d", p.Height)
}
