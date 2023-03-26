<!-- This file is auto-generated. Please do not modify it yourself. -->
 # Protobuf Documentation
 <a name="top"></a>

 ## Table of Contents
 
 - [denom/authorityMetadata.proto](#denom/authorityMetadata.proto)
     - [DenomAuthorityMetadata](#kujira.denom.DenomAuthorityMetadata)
   
 - [denom/params.proto](#denom/params.proto)
     - [Params](#kujira.denom.Params)
   
 - [denom/genesis.proto](#denom/genesis.proto)
     - [GenesisDenom](#kujira.denom.GenesisDenom)
     - [GenesisState](#kujira.denom.GenesisState)
   
 - [denom/query.proto](#denom/query.proto)
     - [QueryDenomAuthorityMetadataRequest](#kujira.denom.QueryDenomAuthorityMetadataRequest)
     - [QueryDenomAuthorityMetadataResponse](#kujira.denom.QueryDenomAuthorityMetadataResponse)
     - [QueryDenomsFromCreatorRequest](#kujira.denom.QueryDenomsFromCreatorRequest)
     - [QueryDenomsFromCreatorResponse](#kujira.denom.QueryDenomsFromCreatorResponse)
     - [QueryParamsRequest](#kujira.denom.QueryParamsRequest)
     - [QueryParamsResponse](#kujira.denom.QueryParamsResponse)
   
     - [Query](#kujira.denom.Query)
   
 - [denom/tx.proto](#denom/tx.proto)
     - [MsgBurn](#kujira.denom.MsgBurn)
     - [MsgBurnResponse](#kujira.denom.MsgBurnResponse)
     - [MsgChangeAdmin](#kujira.denom.MsgChangeAdmin)
     - [MsgChangeAdminResponse](#kujira.denom.MsgChangeAdminResponse)
     - [MsgCreateDenom](#kujira.denom.MsgCreateDenom)
     - [MsgCreateDenomResponse](#kujira.denom.MsgCreateDenomResponse)
     - [MsgMint](#kujira.denom.MsgMint)
     - [MsgMintResponse](#kujira.denom.MsgMintResponse)
   
     - [Msg](#kujira.denom.Msg)
   
 - [Scalar Value Types](#scalar-value-types)

 
 
 <a name="denom/authorityMetadata.proto"></a>
 <p align="right"><a href="#top">Top</a></p>

 ## denom/authorityMetadata.proto
 

 
 <a name="kujira.denom.DenomAuthorityMetadata"></a>

 ### DenomAuthorityMetadata
 DenomAuthorityMetadata specifies metadata for addresses that have specific
capabilities over a token factory denom. Right now there is only one Admin
permission, but is planned to be extended to the future.

 
 | Field | Type | Label | Description |
 | ----- | ---- | ----- | ----------- |
 | `Admin` | [string](#string) |  | Can be empty for no admin, or a valid kujira address |
 
 

 

  <!-- end messages -->

  <!-- end enums -->

  <!-- end HasExtensions -->

  <!-- end services -->

 
 
 <a name="denom/params.proto"></a>
 <p align="right"><a href="#top">Top</a></p>

 ## denom/params.proto
 

 
 <a name="kujira.denom.Params"></a>

 ### Params
 Params holds parameters for the denom module

 
 | Field | Type | Label | Description |
 | ----- | ---- | ----- | ----------- |
 | `creation_fee` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |
 
 

 

  <!-- end messages -->

  <!-- end enums -->

  <!-- end HasExtensions -->

  <!-- end services -->

 
 
 <a name="denom/genesis.proto"></a>
 <p align="right"><a href="#top">Top</a></p>

 ## denom/genesis.proto
 

 
 <a name="kujira.denom.GenesisDenom"></a>

 ### GenesisDenom
 

 
 | Field | Type | Label | Description |
 | ----- | ---- | ----- | ----------- |
 | `denom` | [string](#string) |  |  |
 | `authority_metadata` | [DenomAuthorityMetadata](#kujira.denom.DenomAuthorityMetadata) |  |  |
 
 

 

 
 <a name="kujira.denom.GenesisState"></a>

 ### GenesisState
 GenesisState defines the denom module's genesis state.

 
 | Field | Type | Label | Description |
 | ----- | ---- | ----- | ----------- |
 | `params` | [Params](#kujira.denom.Params) |  | params defines the paramaters of the module. |
 | `factory_denoms` | [GenesisDenom](#kujira.denom.GenesisDenom) | repeated |  |
 
 

 

  <!-- end messages -->

  <!-- end enums -->

  <!-- end HasExtensions -->

  <!-- end services -->

 
 
 <a name="denom/query.proto"></a>
 <p align="right"><a href="#top">Top</a></p>

 ## denom/query.proto
 

 
 <a name="kujira.denom.QueryDenomAuthorityMetadataRequest"></a>

 ### QueryDenomAuthorityMetadataRequest
 

 
 | Field | Type | Label | Description |
 | ----- | ---- | ----- | ----------- |
 | `denom` | [string](#string) |  |  |
 
 

 

 
 <a name="kujira.denom.QueryDenomAuthorityMetadataResponse"></a>

 ### QueryDenomAuthorityMetadataResponse
 

 
 | Field | Type | Label | Description |
 | ----- | ---- | ----- | ----------- |
 | `authority_metadata` | [DenomAuthorityMetadata](#kujira.denom.DenomAuthorityMetadata) |  |  |
 
 

 

 
 <a name="kujira.denom.QueryDenomsFromCreatorRequest"></a>

 ### QueryDenomsFromCreatorRequest
 

 
 | Field | Type | Label | Description |
 | ----- | ---- | ----- | ----------- |
 | `creator` | [string](#string) |  |  |
 
 

 

 
 <a name="kujira.denom.QueryDenomsFromCreatorResponse"></a>

 ### QueryDenomsFromCreatorResponse
 

 
 | Field | Type | Label | Description |
 | ----- | ---- | ----- | ----------- |
 | `denoms` | [string](#string) | repeated |  |
 
 

 

 
 <a name="kujira.denom.QueryParamsRequest"></a>

 ### QueryParamsRequest
 QueryParamsRequest is the request type for the Query/Params RPC method.

 

 

 
 <a name="kujira.denom.QueryParamsResponse"></a>

 ### QueryParamsResponse
 QueryParamsResponse is the response type for the Query/Params RPC method.

 
 | Field | Type | Label | Description |
 | ----- | ---- | ----- | ----------- |
 | `params` | [Params](#kujira.denom.Params) |  | params defines the parameters of the module. |
 
 

 

  <!-- end messages -->

  <!-- end enums -->

  <!-- end HasExtensions -->

 
 <a name="kujira.denom.Query"></a>

 ### Query
 Query defines the gRPC querier service.

 | Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
 | ----------- | ------------ | ------------- | ------------| ------- | -------- |
 | `Params` | [QueryParamsRequest](#kujira.denom.QueryParamsRequest) | [QueryParamsResponse](#kujira.denom.QueryParamsResponse) | Params returns the total set of minting parameters. | GET|/kujira/denoms/params|
 | `DenomAuthorityMetadata` | [QueryDenomAuthorityMetadataRequest](#kujira.denom.QueryDenomAuthorityMetadataRequest) | [QueryDenomAuthorityMetadataResponse](#kujira.denom.QueryDenomAuthorityMetadataResponse) |  | GET|/kujira/denoms/{denom}/authority_metadata|
 | `DenomsFromCreator` | [QueryDenomsFromCreatorRequest](#kujira.denom.QueryDenomsFromCreatorRequest) | [QueryDenomsFromCreatorResponse](#kujira.denom.QueryDenomsFromCreatorResponse) |  | GET|/kujira/denoms/by_creator/{creator}|
 
  <!-- end services -->

 
 
 <a name="denom/tx.proto"></a>
 <p align="right"><a href="#top">Top</a></p>

 ## denom/tx.proto
 

 
 <a name="kujira.denom.MsgBurn"></a>

 ### MsgBurn
 MsgBurn is the sdk.Msg type for allowing an admin account to burn
a token.  For now, we only support burning from the sender account.

 
 | Field | Type | Label | Description |
 | ----- | ---- | ----- | ----------- |
 | `sender` | [string](#string) |  |  |
 | `amount` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |
 
 

 

 
 <a name="kujira.denom.MsgBurnResponse"></a>

 ### MsgBurnResponse
 

 

 

 
 <a name="kujira.denom.MsgChangeAdmin"></a>

 ### MsgChangeAdmin
 MsgChangeAdmin is the sdk.Msg type for allowing an admin account to reassign
adminship of a denom to a new account

 
 | Field | Type | Label | Description |
 | ----- | ---- | ----- | ----------- |
 | `sender` | [string](#string) |  |  |
 | `denom` | [string](#string) |  |  |
 | `newAdmin` | [string](#string) |  |  |
 
 

 

 
 <a name="kujira.denom.MsgChangeAdminResponse"></a>

 ### MsgChangeAdminResponse
 

 

 

 
 <a name="kujira.denom.MsgCreateDenom"></a>

 ### MsgCreateDenom
 MsgCreateDenom is the sdk.Msg type for allowing an account to create
a new denom.  It requires a sender address and a unique nonce
(to allow accounts to create multiple denoms)

 
 | Field | Type | Label | Description |
 | ----- | ---- | ----- | ----------- |
 | `sender` | [string](#string) |  |  |
 | `nonce` | [string](#string) |  |  |
 
 

 

 
 <a name="kujira.denom.MsgCreateDenomResponse"></a>

 ### MsgCreateDenomResponse
 MsgCreateDenomResponse is the return value of MsgCreateDenom
It returns the full string of the newly created denom

 
 | Field | Type | Label | Description |
 | ----- | ---- | ----- | ----------- |
 | `new_token_denom` | [string](#string) |  |  |
 
 

 

 
 <a name="kujira.denom.MsgMint"></a>

 ### MsgMint
 MsgMint is the sdk.Msg type for allowing an admin account to mint
more of a token.

 
 | Field | Type | Label | Description |
 | ----- | ---- | ----- | ----------- |
 | `sender` | [string](#string) |  |  |
 | `amount` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |
 | `recipient` | [string](#string) |  |  |
 
 

 

 
 <a name="kujira.denom.MsgMintResponse"></a>

 ### MsgMintResponse
 

 

 

  <!-- end messages -->

  <!-- end enums -->

  <!-- end HasExtensions -->

 
 <a name="kujira.denom.Msg"></a>

 ### Msg
 Msg defines the Msg service.

 | Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
 | ----------- | ------------ | ------------- | ------------| ------- | -------- |
 | `CreateDenom` | [MsgCreateDenom](#kujira.denom.MsgCreateDenom) | [MsgCreateDenomResponse](#kujira.denom.MsgCreateDenomResponse) |  | |
 | `Mint` | [MsgMint](#kujira.denom.MsgMint) | [MsgMintResponse](#kujira.denom.MsgMintResponse) |  | |
 | `Burn` | [MsgBurn](#kujira.denom.MsgBurn) | [MsgBurnResponse](#kujira.denom.MsgBurnResponse) |  | |
 | `ChangeAdmin` | [MsgChangeAdmin](#kujira.denom.MsgChangeAdmin) | [MsgChangeAdminResponse](#kujira.denom.MsgChangeAdminResponse) | ForceTransfer is deactivated for now because we need to think through edge cases rpc ForceTransfer(MsgForceTransfer) returns (MsgForceTransferResponse); | |
 
  <!-- end services -->

 

 ## Scalar Value Types

 | .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
 | ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
 | <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
 | <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
 | <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
 | <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
 | <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
 | <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
 | <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
 | <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
 | <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
 | <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
 | <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
 | <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
 | <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
 | <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
 | <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |
 