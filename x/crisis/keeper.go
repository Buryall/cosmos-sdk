package crisis

import (
	"fmt"
	"time"

	"github.com/tendermint/tendermint/libs/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/crisis/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// Keeper - crisis keeper
type Keeper struct {
	routes         []types.InvarRoute
	paramSpace     params.Subspace
	invCheckPeriod uint

	distrKeeper  types.DistributionKeeper
	supplyKeeper types.SupplyKeeper
}

// NewKeeper creates a new Keeper object
func NewKeeper(paramSpace params.Subspace, invCheckPeriod uint,
	distrKeeper types.DistributionKeeper, supplyKeeper types.SupplyKeeper) Keeper {

	return Keeper{
		routes:         []types.InvarRoute{},
		paramSpace:     paramSpace.WithKeyTable(types.ParamKeyTable()),
		invCheckPeriod: invCheckPeriod,
		distrKeeper:    distrKeeper,
		supplyKeeper:   supplyKeeper,
	}
}

// register routes for the
func (k *Keeper) RegisterRoute(moduleName, route string, invar sdk.Invariant) {
	invarRoute := types.NewInvarRoute(moduleName, route, invar)
	k.routes = append(k.routes, invarRoute)
}

// Routes - return the keeper's invariant routes
func (k Keeper) Routes() []types.InvarRoute {
	return k.routes
}

// Invariants returns all the registered Crisis keeper invariants.
func (k Keeper) Invariants() []sdk.Invariant {
	var invars []sdk.Invariant
	for _, route := range k.routes {
		invars = append(invars, route.Invar)
	}
	return invars
}

// assert all invariants
func (k Keeper) AssertInvariants(ctx sdk.Context, logger log.Logger) {

	start := time.Now()
	invarRoutes := k.Routes()
	for _, ir := range invarRoutes {
		if err := ir.Invar(ctx); err != nil {

			// TODO: Include app name as part of context to allow for this to be
			// variable.
			panic(fmt.Errorf("invariant broken: %s\n"+
				"\tCRITICAL please submit the following transaction:\n"+
				"\t\t tx crisis invariant-broken %v %v", err, ir.ModuleName, ir.Route))
		}
	}

	end := time.Now()
	diff := end.Sub(start)

	logger.With("module", "x/crisis").Info("asserted all invariants", "duration", diff, "height", ctx.BlockHeight())
}

// DONTCOVER
