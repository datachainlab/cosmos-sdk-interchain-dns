package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	commontypes "github.com/datachainlab/cosmos-sdk-interchain-dns/x/ibc-dns/common/types"
	clienttypes "github.com/datachainlab/cosmos-sdk-interchain-dns/x/ibc-dns/client/types"
	servertypes "github.com/datachainlab/cosmos-sdk-interchain-dns/x/ibc-dns/server/types"
)

// ModuleCdc is the codec for the module
var ModuleCdc = codec.New()

func init() {
	RegisterCodec(ModuleCdc)
}

func RegisterCodec(cdc *codec.Codec) {
	commontypes.RegisterCodec(cdc)
	clienttypes.RegisterCodec(cdc)
	servertypes.RegisterCodec(cdc)
}