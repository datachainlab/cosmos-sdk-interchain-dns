package client

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/core/04-channel/types"

	"github.com/datachainlab/cosmos-sdk-interchain-dns/x/ibc-dns/client/keeper"
	"github.com/datachainlab/cosmos-sdk-interchain-dns/x/ibc-dns/client/types"
	commontypes "github.com/datachainlab/cosmos-sdk-interchain-dns/x/ibc-dns/common/types"
	servertypes "github.com/datachainlab/cosmos-sdk-interchain-dns/x/ibc-dns/server/types"
)

// NewHandler returns a handler
func NewHandler(keeper types.MsgServer) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgRegisterDomain:
			res, err := keeper.RegisterDomain(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgDomainAssociationCreate:
			res, err := keeper.DomainAssociationCreate(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		default:
			return nil, commontypes.ErrUnknownRequest
		}
	}
}

// NewPacketReceiver returns a new PacketReceiver
func NewPacketReceiver(keeper keeper.Keeper) commontypes.PacketReceiver {
	return func(ctx sdk.Context, packet channeltypes.Packet) (*sdk.Result, []byte, error) {
		data, err := commontypes.DeserializeJSONPacketData(servertypes.PacketCdc(), packet.GetData())
		if err != nil {
			return nil, nil, err
		}
		switch data := data.(type) {
		case *servertypes.DomainAssociationResultPacketData:
			return handleDomainAssociationResultPacketData(ctx, keeper, packet, data)
		default:
			return nil, nil, commontypes.ErrUnknownRequest
		}
	}
}

func handleDomainAssociationResultPacketData(
	ctx sdk.Context,
	keeper keeper.Keeper,
	packet channeltypes.Packet,
	data *servertypes.DomainAssociationResultPacketData,
) (*sdk.Result, []byte, error) {
	switch data.Status {
	case servertypes.STATUS_OK:
		err := keeper.ReceiveDomainAssociationResultPacketData(
			ctx,
			packet,
			data,
		)
		if err != nil {
			return nil, nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "failed to handle a packet 'DomainAssociationResultPacketData: %v'", err)
		}
	case servertypes.STATUS_FAILED:
		// TODO cleanup
	default:
		return nil, nil, fmt.Errorf("unknown status '%v'", data.Status)
	}

	ack := servertypes.NewDomainAssociationResultPacketAcknowledgement()
	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, ack.GetBytes(), nil
}

// NewPacketAcknowledgementReceiver returns a new PacketAcknowledgementReceiver
func NewPacketAcknowledgementReceiver(keeper keeper.Keeper) commontypes.PacketAcknowledgementReceiver {
	return func(ctx sdk.Context, packet channeltypes.Packet, ack []byte) (*sdk.Result, error) {
		ackData, err := commontypes.DeserializeJSONPacketAck(servertypes.PacketCdc(), packet.GetData())
		if err != nil {
			return nil, err
		}
		switch ackData := ackData.(type) {
		case *servertypes.RegisterDomainPacketAcknowledgement:
			return handleRegisterDomainPacketAcknowledgement(ctx, keeper, ackData, packet)
		case *servertypes.DomainAssociationCreatePacketAcknowledgement:
			return handleDomainAssociationCreatePacketAcknowledgement(ctx, keeper, ackData)
		default:
			return nil, commontypes.ErrUnknownRequest
		}
	}
}

func handleRegisterDomainPacketAcknowledgement(ctx sdk.Context, k keeper.Keeper, ack *servertypes.RegisterDomainPacketAcknowledgement, packet channeltypes.Packet) (*sdk.Result, error) {
	if err := k.ReceiveRegisterDomainPacketAcknowledgement(ctx, ack.Status, ack.DomainName, packet); err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "failed to handle a packet 'RegisterDomainPacketAcknowledgement: %v'", err)
	}
	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleDomainAssociationCreatePacketAcknowledgement(ctx sdk.Context, k keeper.Keeper, ack *servertypes.DomainAssociationCreatePacketAcknowledgement) (*sdk.Result, error) {
	if err := k.ReceiveDomainAssociationCreatePacketAcknowledgement(ctx, ack.Status); err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "failed to handle a packet 'DomainAssociationCreatePacketAcknowledgement: %v'", err)
	}
	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}
