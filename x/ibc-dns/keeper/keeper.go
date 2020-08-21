package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/capability"
	porttypes "github.com/cosmos/cosmos-sdk/x/ibc/05-port/types"
	"github.com/datachainlab/cosmos-sdk-interchain-dns/x/ibc-dns/common/types"
)

type Keeper struct {
	portKeeper   types.PortKeeper
	scopedKeeper capability.ScopedKeeper
}

func NewKeeper(portKeeper types.PortKeeper, scopedKeeper capability.ScopedKeeper) Keeper {
	return Keeper{
		portKeeper:   portKeeper,
		scopedKeeper: scopedKeeper,
	}
}

// BindPort defines a wrapper function for the ort TPCKeeper's function in
// order to expose it to module's InitGenesis function
func (k Keeper) BindPort(ctx sdk.Context, portID string) (*capability.Capability, error) {
	cap := k.portKeeper.BindPort(ctx, portID)
	if err := k.ClaimCapability(ctx, cap, porttypes.PortPath(portID)); err != nil {
		return nil, err
	}
	return cap, nil
}

// ClaimCapability allows the transfer module that can claim a capability that IBC module
// passes to it
func (k Keeper) ClaimCapability(ctx sdk.Context, cap *capability.Capability, name string) error {
	return k.scopedKeeper.ClaimCapability(ctx, cap, name)
}