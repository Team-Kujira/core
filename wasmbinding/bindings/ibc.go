package bindings

import (
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	"github.com/cosmos/ibc-go/v6/modules/core/03-connection/types"
	commitmenttypes "github.com/cosmos/ibc-go/v6/modules/core/23-commitment/types"
	host "github.com/cosmos/ibc-go/v6/modules/core/24-host"
	"github.com/cosmos/ibc-go/v6/modules/core/exported"
	ibckeeper "github.com/cosmos/ibc-go/v6/modules/core/keeper"
	ibctmtypes "github.com/cosmos/ibc-go/v6/modules/light-clients/07-tendermint/types"
)

type VerifyMemberShipQuery struct {
	Connection     string `json:"connection"`
	RevisionNumber uint64 `json:"revision_number"`
	RevisionHeight uint64 `json:"revision_height"`
	Proof          []byte `josn:"proof"`
	Value          []byte `json:"value"`
}

type VerifyNonMemberShipQuery struct {
	Connection     string `json:"connection"`
	RevisionNumber uint64 `json:"revision_number"`
	RevisionHeight uint64 `json:"revision_height"`
	Proof          []byte `josn:"proof"`
}

// ----- moved from ibc-go/modules/core/03-connection -----
// getClientStateAndVerificationStore returns the client state and associated KVStore for the provided client identifier.
// If the client type is localhost then the core IBC KVStore is returned, otherwise the client prefixed store is returned.
func getClientStateAndVerificationStore(ctx sdk.Context, keeper ibckeeper.Keeper, clientID string, ibcStoreKey *storetypes.KVStoreKey) (exported.ClientState, sdk.KVStore, error) {
	clientState, found := keeper.ClientKeeper.GetClientState(ctx, clientID)
	if !found {
		return nil, nil, sdkerrors.Wrap(clienttypes.ErrClientNotFound, clientID)
	}

	store := keeper.ClientKeeper.ClientStore(ctx, clientID)
	if clientID == exported.Localhost {
		store = ctx.KVStore(ibcStoreKey)
	}

	return clientState, store, nil
}

// getClientStateAndStore checks the client and connection status and returns the client state and associated KVStore
func getClientStateAndStore(ctx sdk.Context, keeper ibckeeper.Keeper, ibcStoreKey *storetypes.KVStoreKey, connection types.ConnectionEnd) (exported.ClientState, storetypes.KVStore, error) {
	// verify the client is ACTIVE
	clientID := connection.GetClientID()
	clientState, clientStore, err := getClientStateAndVerificationStore(ctx, keeper, clientID, ibcStoreKey)
	if err != nil {
		return clientState, clientStore, err
	}

	_, ok := clientState.(*ibctmtypes.ClientState)
	if !ok {
		return clientState, clientStore, sdkerrors.Wrapf(clienttypes.ErrInvalidClient, "invalid client type %T, expected %T", clientState, &ibctmtypes.ClientState{})
	}

	if status := clientState.Status(ctx, clientStore, keeper.Codec()); status != exported.Active {
		return clientState, clientStore, sdkerrors.Wrapf(clienttypes.ErrClientNotActive, "client (%s) status is %s", clientID, status)
	}

	// verify the connection state is OPEN
	if connection.State != types.OPEN {
		return clientState, clientStore, sdkerrors.Wrapf(
			types.ErrInvalidConnectionState,
			"connection state is not OPEN (got %s)", connection.State.String(),
		)
	}

	return clientState, clientStore, nil
}

// getConsStateAndMerklePath generates merkle path using connection and get consensus state using height from the client store
func getConsStateAndMerklePath(keeper ibckeeper.Keeper, clientState exported.ClientState, clientStore storetypes.KVStore, connection types.ConnectionEnd, height clienttypes.Height) (*ibctmtypes.ConsensusState, *commitmenttypes.MerklePath, error) {
	if clientState.GetLatestHeight().LT(height) {
		return nil, nil, sdkerrors.Wrapf(
			sdkerrors.ErrInvalidHeight,
			"client state height < proof height (%d < %d), please ensure the client has been updated", clientState.GetLatestHeight(), height,
		)
	}

	merklePathWithoutPrefix := commitmenttypes.NewMerklePath(host.FullClientStatePath(connection.GetCounterparty().GetClientID()))
	merklePath, err := commitmenttypes.ApplyPrefix(connection.GetCounterparty().GetPrefix(), merklePathWithoutPrefix)
	if err != nil {
		return nil, nil, err
	}

	consensusState, err := ibctmtypes.GetConsensusState(clientStore, keeper.Codec(), height)
	if err != nil {
		return nil, nil, sdkerrors.Wrap(clienttypes.ErrConsensusStateNotFound, "please ensure the proof was constructed against a height that exists on the client")
	}

	return consensusState, &merklePath, nil
}

func HandleIBCQuery(ctx sdk.Context, keeper ibckeeper.Keeper, ibcStoreKey *storetypes.KVStoreKey, q *IBCQuery) error {
	switch {
	case q.VerifyMemberShip != nil:
		connectionID := q.VerifyMemberShip.Connection
		connection, found := keeper.ConnectionKeeper.GetConnection(ctx, connectionID)
		if !found {
			return sdkerrors.Wrap(types.ErrConnectionNotFound, connectionID)
		}

		clientState, clientStore, err := getClientStateAndStore(ctx, keeper, ibcStoreKey, connection)
		if err != nil {
			return err
		}

		var merkleProof commitmenttypes.MerkleProof
		if err := keeper.Codec().Unmarshal(q.VerifyMemberShip.Proof, &merkleProof); err != nil {
			return sdkerrors.Wrap(commitmenttypes.ErrInvalidProof, "failed to unmarshal proof into ICS 23 commitment merkle proof")
		}

		height := clienttypes.NewHeight(q.VerifyMemberShip.RevisionNumber, q.VerifyMemberShip.RevisionHeight)
		consState, merklePath, err := getConsStateAndMerklePath(keeper, clientState, clientStore, connection, height)
		if err != nil {
			return err
		}

		if err := merkleProof.VerifyMembership(clientState.(*ibctmtypes.ClientState).ProofSpecs, consState.GetRoot(), merklePath, q.VerifyMemberShip.Value); err != nil { //nolint
			return err
		}

		return nil

	case q.VerifyNonMemberShip != nil:
		connectionID := q.VerifyNonMemberShip.Connection
		connection, found := keeper.ConnectionKeeper.GetConnection(ctx, connectionID)
		if !found {
			return sdkerrors.Wrap(types.ErrConnectionNotFound, connectionID)
		}

		clientState, clientStore, err := getClientStateAndStore(ctx, keeper, ibcStoreKey, connection)
		if err != nil {
			return err
		}

		var merkleProof commitmenttypes.MerkleProof
		if err := keeper.Codec().Unmarshal(q.VerifyNonMemberShip.Proof, &merkleProof); err != nil {
			return sdkerrors.Wrap(commitmenttypes.ErrInvalidProof, "failed to unmarshal proof into ICS 23 commitment merkle proof")
		}

		height := clienttypes.NewHeight(q.VerifyNonMemberShip.RevisionNumber, q.VerifyNonMemberShip.RevisionHeight)
		consState, merklePath, err := getConsStateAndMerklePath(keeper, clientState, clientStore, connection, height)
		if err != nil {
			return err
		}

		if err := merkleProof.VerifyNonMembership(clientState.(*ibctmtypes.ClientState).ProofSpecs, consState.GetRoot(), merklePath); err != nil { //nolint
			return err
		}

		return nil
	default:
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized query request")
	}
}
