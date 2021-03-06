package server

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/core/04-channel/types"

	commontypes "github.com/datachainlab/cosmos-sdk-interchain-dns/x/ibc-dns/common/types"
	"github.com/datachainlab/cosmos-sdk-interchain-dns/x/ibc-dns/server/keeper"
	"github.com/datachainlab/cosmos-sdk-interchain-dns/x/ibc-dns/server/types"
)

// NewPacketReceiver returns a receiver to handle received packets
func NewPacketReceiver(keeper keeper.Keeper) commontypes.PacketReceiver {
	return func(ctx sdk.Context, packet channeltypes.Packet) (*sdk.Result, []byte, error) {
		data, err := commontypes.DeserializeJSONPacketData(types.PacketCdc(), packet.GetData())
		if err != nil {
			return nil, nil, err
		}
		switch data := data.(type) {
		case *types.RegisterDomainPacketData:
			return handlePacketRegisterChannelDomain(ctx, keeper, packet, data)
		case *types.DomainAssociationCreatePacketData:
			return handleDomainAssociationCreatePacketData(ctx, keeper, packet, data)
		default:
			return nil, nil, commontypes.ErrUnknownRequest
		}
	}
}

func handlePacketRegisterChannelDomain(ctx sdk.Context, keeper keeper.Keeper, packet channeltypes.Packet, data *types.RegisterDomainPacketData) (*sdk.Result, []byte, error) {
	var status uint32
	if err := keeper.ReceivePacketRegisterDomain(ctx, packet, data); err != nil {
		ctx.Logger().Info("failed to handle a packet 'PacketRegisterChannelDomain'", "err", err)
		status = types.STATUS_FAILED
	} else {
		status = types.STATUS_OK
	}
	ack := types.NewRegisterDomainPacketAcknowledgement(status, data.DomainName)
	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, ack.GetBytes(), nil
}

func handleDomainAssociationCreatePacketData(ctx sdk.Context, keeper keeper.Keeper, packet channeltypes.Packet, data *types.DomainAssociationCreatePacketData) (*sdk.Result, []byte, error) {
	ack, completed := keeper.ReceiveDomainAssociationCreatePacketData(ctx, packet, data)
	if completed {
		err := keeper.SendDomainAssociationResultPacketData(ctx, ack.Status, data.SrcClient, data.DstClient)
		if err != nil {
			return nil, nil, err
		}
	}

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, ack.GetBytes(), nil
}

// NewPacketAcknowledgementReceiver returns a receiver to handle received acks
func NewPacketAcknowledgementReceiver(keeper keeper.Keeper) commontypes.PacketAcknowledgementReceiver {
	return func(ctx sdk.Context, packet channeltypes.Packet, ack []byte) (*sdk.Result, error) {
		ackData, err := commontypes.DeserializeJSONPacketAck(types.PacketCdc(), packet.GetData())
		if err != nil {
			return nil, err
		}
		switch ackData := ackData.(type) {
		case *types.DomainAssociationResultPacketAcknowledgement:
			return handleDomainAssociationResultPacketAcknowledgement(ctx, keeper, ackData)
		default:
			return nil, commontypes.ErrUnknownRequest
		}
	}
}

func handleDomainAssociationResultPacketAcknowledgement(ctx sdk.Context, k keeper.Keeper, ackData *types.DomainAssociationResultPacketAcknowledgement) (*sdk.Result, error) {
	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}
