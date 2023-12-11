package bindings

import (
	"cosmossdk.io/errors"
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types" //nolint
	"github.com/cosmos/ibc-go/v8/modules/core/03-connection/types"
	commitmenttypes "github.com/cosmos/ibc-go/v8/modules/core/23-commitment/types"
	"github.com/cosmos/ibc-go/v8/modules/core/exported"
	ibckeeper "github.com/cosmos/ibc-go/v8/modules/core/keeper"
	ibctmtypes "github.com/cosmos/ibc-go/v8/modules/light-clients/07-tendermint"
)

type VerifyMembershipQuery struct {
	Connection     string `json:"connection"`
	RevisionNumber uint64 `json:"revision_number"`
	RevisionHeight uint64 `json:"revision_height"`
	Proof          []byte `json:"proof"`
	Value          []byte `json:"value"`
	PathPrefix     string `json:"path_prefix"`
	PathKey        string `json:"path_key"`
}

type VerifyNonMembershipQuery struct {
	Connection     string `json:"connection"`
	RevisionNumber uint64 `json:"revision_number"`
	RevisionHeight uint64 `json:"revision_height"`
	Proof          []byte `json:"proof"`
	PathPrefix     string `json:"path_prefix"`
	PathKey        string `json:"path_key"`
}

type IbcVerifyResponse struct{}

// ----- moved from ibc-go/modules/core/03-connection -----
// getClientStateAndVerificationStore returns the client state and associated KVStore for the provided client identifier.
// If the client type is localhost then the core IBC KVStore is returned, otherwise the client prefixed store is returned.
func getClientStateAndVerificationStore(ctx sdk.Context, keeper ibckeeper.Keeper, clientID string, ibcStoreKey *storetypes.KVStoreKey) (exported.ClientState, storetypes.KVStore, error) {
	clientState, found := keeper.ClientKeeper.GetClientState(ctx, clientID)
	if !found {
		return nil, nil, errors.Wrap(clienttypes.ErrClientNotFound, clientID)
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
		return clientState, clientStore, errors.Wrapf(clienttypes.ErrInvalidClient, "invalid client type %T, expected %T", clientState, &ibctmtypes.ClientState{})
	}

	if status := clientState.Status(ctx, clientStore, keeper.Codec()); status != exported.Active {
		return clientState, clientStore, errors.Wrapf(clienttypes.ErrClientNotActive, "client (%s) status is %s", clientID, status)
	}

	// verify the connection state is OPEN
	if connection.State != types.OPEN {
		return clientState, clientStore, errors.Wrapf(
			types.ErrInvalidConnectionState,
			"connection state is not OPEN (got %s)", connection.State.String(),
		)
	}

	return clientState, clientStore, nil
}

// getConsStateAndMerklePath generates merkle path using path and get consensus state using height from the client store
func getConsStateAndMerklePath(keeper ibckeeper.Keeper, clientState exported.ClientState, clientStore storetypes.KVStore, height clienttypes.Height, pathPrefix string, pathKey string) (*ibctmtypes.ConsensusState, *commitmenttypes.MerklePath, error) {
	if clientState.GetLatestHeight().LT(height) {
		return nil, nil, errors.Wrapf(
			sdkerrors.ErrInvalidHeight,
			"client state height < proof height (%d < %d), please ensure the client has been updated", clientState.GetLatestHeight(), height,
		)
	}

	merklePath := commitmenttypes.NewMerklePath(pathPrefix, pathKey)

	consensusState, found := ibctmtypes.GetConsensusState(clientStore, keeper.Codec(), height)
	if !found {
		return nil, nil, errors.Wrap(clienttypes.ErrConsensusStateNotFound, "please ensure the proof was constructed against a height that exists on the client")
	}

	return consensusState, &merklePath, nil
}

func HandleIBCQuery(ctx sdk.Context, keeper ibckeeper.Keeper, ibcStoreKey *storetypes.KVStoreKey, q *IbcQuery) error {
	switch {
	case q.VerifyMembership != nil:
		connectionID := q.VerifyMembership.Connection
		connection, found := keeper.ConnectionKeeper.GetConnection(ctx, connectionID)
		if !found {
			return errors.Wrap(types.ErrConnectionNotFound, connectionID)
		}

		clientState, clientStore, err := getClientStateAndStore(ctx, keeper, ibcStoreKey, connection)
		if err != nil {
			return err
		}

		var merkleProof commitmenttypes.MerkleProof
		if err := keeper.Codec().Unmarshal(q.VerifyMembership.Proof, &merkleProof); err != nil {
			return errors.Wrap(commitmenttypes.ErrInvalidProof, "failed to unmarshal proof into ICS 23 commitment merkle proof")
		}

		height := clienttypes.NewHeight(q.VerifyMembership.RevisionNumber, q.VerifyMembership.RevisionHeight)
		consState, merklePath, err := getConsStateAndMerklePath(keeper, clientState, clientStore, height, q.VerifyMembership.PathPrefix, q.VerifyMembership.PathKey)
		if err != nil {
			return err
		}

		if err := merkleProof.VerifyMembership(clientState.(*ibctmtypes.ClientState).ProofSpecs, consState.GetRoot(), *merklePath, q.VerifyMembership.Value); err != nil { //nolint
			return err
		}

		return nil

	case q.VerifyNonMembership != nil:
		connectionID := q.VerifyNonMembership.Connection
		connection, found := keeper.ConnectionKeeper.GetConnection(ctx, connectionID)
		if !found {
			return errors.Wrap(types.ErrConnectionNotFound, connectionID)
		}

		clientState, clientStore, err := getClientStateAndStore(ctx, keeper, ibcStoreKey, connection)
		if err != nil {
			return err
		}

		var merkleProof commitmenttypes.MerkleProof
		if err := keeper.Codec().Unmarshal(q.VerifyNonMembership.Proof, &merkleProof); err != nil {
			return errors.Wrap(commitmenttypes.ErrInvalidProof, "failed to unmarshal proof into ICS 23 commitment merkle proof")
		}

		height := clienttypes.NewHeight(q.VerifyNonMembership.RevisionNumber, q.VerifyNonMembership.RevisionHeight)
		consState, merklePath, err := getConsStateAndMerklePath(keeper, clientState, clientStore, height, q.VerifyMembership.PathPrefix, q.VerifyMembership.PathKey)
		if err != nil {
			return err
		}

		if err := merkleProof.VerifyNonMembership(clientState.(*ibctmtypes.ClientState).ProofSpecs, consState.GetRoot(), merklePath); err != nil { //nolint
			return err
		}

		return nil
	default:
		return errors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized query request")
	}
}
