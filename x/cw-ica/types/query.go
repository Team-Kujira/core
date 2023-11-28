package types

// NewQueryInterchainAccountRequest creates and returns a new QueryInterchainAccountRequest
func NewQueryInterchainAccountRequest(connectionID, accountID, owner string) *QueryInterchainAccountRequest {
	return &QueryInterchainAccountRequest{
		ConnectionId: connectionID,
		AccountId:    accountID,
		Owner:        owner,
	}
}

// NewQueryInterchainAccountResponse creates and returns a new QueryInterchainAccountResponse
func NewQueryInterchainAccountResponse(interchainAccAddr string) *QueryInterchainAccountResponse {
	return &QueryInterchainAccountResponse{
		InterchainAccountAddress: interchainAccAddr,
	}
}
