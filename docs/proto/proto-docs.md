<!-- This file is auto-generated. Please do not modify it yourself. -->
 # Protobuf Documentation
 <a name="top"></a>

 ## Table of Contents
 
 - [kujira/oracle/oracle.proto](#kujira/oracle/oracle.proto)
     - [Denom](#kujira.oracle.Denom)
     - [ExchangeRateTuple](#kujira.oracle.ExchangeRateTuple)
     - [Params](#kujira.oracle.Params)
   
 - [kujira/oracle/genesis.proto](#kujira/oracle/genesis.proto)
     - [FeederDelegation](#kujira.oracle.FeederDelegation)
     - [GenesisState](#kujira.oracle.GenesisState)
     - [MissCounter](#kujira.oracle.MissCounter)
   
 - [kujira/oracle/query.proto](#kujira/oracle/query.proto)
     - [QueryActivesRequest](#kujira.oracle.QueryActivesRequest)
     - [QueryActivesResponse](#kujira.oracle.QueryActivesResponse)
     - [QueryExchangeRateRequest](#kujira.oracle.QueryExchangeRateRequest)
     - [QueryExchangeRateResponse](#kujira.oracle.QueryExchangeRateResponse)
     - [QueryExchangeRatesRequest](#kujira.oracle.QueryExchangeRatesRequest)
     - [QueryExchangeRatesResponse](#kujira.oracle.QueryExchangeRatesResponse)
     - [QueryMissCounterRequest](#kujira.oracle.QueryMissCounterRequest)
     - [QueryMissCounterResponse](#kujira.oracle.QueryMissCounterResponse)
     - [QueryParamsRequest](#kujira.oracle.QueryParamsRequest)
     - [QueryParamsResponse](#kujira.oracle.QueryParamsResponse)
     - [QueryVoteTargetsRequest](#kujira.oracle.QueryVoteTargetsRequest)
     - [QueryVoteTargetsResponse](#kujira.oracle.QueryVoteTargetsResponse)
   
     - [Query](#kujira.oracle.Query)
   
 - [kujira/oracle/tx.proto](#kujira/oracle/tx.proto)
     - [MsgAddRequiredDenom](#kujira.oracle.MsgAddRequiredDenom)
     - [MsgAddRequiredDenomResponse](#kujira.oracle.MsgAddRequiredDenomResponse)
     - [MsgRemoveRequiredDenom](#kujira.oracle.MsgRemoveRequiredDenom)
     - [MsgRemoveRequiredDenomResponse](#kujira.oracle.MsgRemoveRequiredDenomResponse)
     - [MsgUpdateParams](#kujira.oracle.MsgUpdateParams)
     - [MsgUpdateParamsResponse](#kujira.oracle.MsgUpdateParamsResponse)
   
     - [Msg](#kujira.oracle.Msg)
   
 - [Scalar Value Types](#scalar-value-types)

 
 
 <a name="kujira/oracle/oracle.proto"></a>
 <p align="right"><a href="#top">Top</a></p>

 ## kujira/oracle/oracle.proto
 

 
 <a name="kujira.oracle.Denom"></a>

 ### Denom
 Denom - the object to hold configurations of each denom

 
 | Field | Type | Label | Description |
 | ----- | ---- | ----- | ----------- |
 | `name` | [string](#string) |  |  |
 
 

 

 
 <a name="kujira.oracle.ExchangeRateTuple"></a>

 ### ExchangeRateTuple
 ExchangeRateTuple - struct to store interpreted exchange rates data to store

 
 | Field | Type | Label | Description |
 | ----- | ---- | ----- | ----------- |
 | `denom` | [string](#string) |  |  |
 | `exchange_rate` | [string](#string) |  |  |
 
 

 

 
 <a name="kujira.oracle.Params"></a>

 ### Params
 Params defines the parameters for the oracle module.

 
 | Field | Type | Label | Description |
 | ----- | ---- | ----- | ----------- |
 | `vote_period` | [uint64](#uint64) |  |  |
 | `vote_threshold` | [string](#string) |  |  |
 | `max_deviation` | [string](#string) |  |  |
 | `required_denoms` | [string](#string) | repeated |  |
 | `slash_fraction` | [string](#string) |  |  |
 | `slash_window` | [uint64](#uint64) |  |  |
 | `min_valid_per_window` | [string](#string) |  |  |
 | `reward_band` | [string](#string) |  | Deprecated |
 | `whitelist` | [Denom](#kujira.oracle.Denom) | repeated |  |
 
 

 

  <!-- end messages -->

  <!-- end enums -->

  <!-- end HasExtensions -->

  <!-- end services -->

 
 
 <a name="kujira/oracle/genesis.proto"></a>
 <p align="right"><a href="#top">Top</a></p>

 ## kujira/oracle/genesis.proto
 

 
 <a name="kujira.oracle.FeederDelegation"></a>

 ### FeederDelegation
 FeederDelegation is the address for where oracle feeder authority are
delegated to. By default this struct is only used at genesis to feed in
default feeder addresses.

 
 | Field | Type | Label | Description |
 | ----- | ---- | ----- | ----------- |
 | `feeder_address` | [string](#string) |  |  |
 | `validator_address` | [string](#string) |  |  |
 
 

 

 
 <a name="kujira.oracle.GenesisState"></a>

 ### GenesisState
 GenesisState defines the oracle module's genesis state.

 
 | Field | Type | Label | Description |
 | ----- | ---- | ----- | ----------- |
 | `params` | [Params](#kujira.oracle.Params) |  |  |
 | `exchange_rates` | [ExchangeRateTuple](#kujira.oracle.ExchangeRateTuple) | repeated |  |
 | `miss_counters` | [MissCounter](#kujira.oracle.MissCounter) | repeated |  |
 
 

 

 
 <a name="kujira.oracle.MissCounter"></a>

 ### MissCounter
 MissCounter defines an miss counter and validator address pair used in
oracle module's genesis state

 
 | Field | Type | Label | Description |
 | ----- | ---- | ----- | ----------- |
 | `validator_address` | [string](#string) |  |  |
 | `miss_counter` | [uint64](#uint64) |  |  |
 
 

 

  <!-- end messages -->

  <!-- end enums -->

  <!-- end HasExtensions -->

  <!-- end services -->

 
 
 <a name="kujira/oracle/query.proto"></a>
 <p align="right"><a href="#top">Top</a></p>

 ## kujira/oracle/query.proto
 

 
 <a name="kujira.oracle.QueryActivesRequest"></a>

 ### QueryActivesRequest
 QueryActivesRequest is the request type for the Query/Actives RPC method.

 

 

 
 <a name="kujira.oracle.QueryActivesResponse"></a>

 ### QueryActivesResponse
 QueryActivesResponse is response type for the
Query/Actives RPC method.

 
 | Field | Type | Label | Description |
 | ----- | ---- | ----- | ----------- |
 | `actives` | [string](#string) | repeated | actives defines a list of the denomination which oracle prices aggreed upon. |
 
 

 

 
 <a name="kujira.oracle.QueryExchangeRateRequest"></a>

 ### QueryExchangeRateRequest
 QueryExchangeRateRequest is the request type for the Query/ExchangeRate RPC method.

 
 | Field | Type | Label | Description |
 | ----- | ---- | ----- | ----------- |
 | `denom` | [string](#string) |  | denom defines the denomination to query for. |
 
 

 

 
 <a name="kujira.oracle.QueryExchangeRateResponse"></a>

 ### QueryExchangeRateResponse
 QueryExchangeRateResponse is response type for the
Query/ExchangeRate RPC method.

 
 | Field | Type | Label | Description |
 | ----- | ---- | ----- | ----------- |
 | `exchange_rate` | [string](#string) |  | exchange_rate defines the exchange rate of whitelisted assets |
 
 

 

 
 <a name="kujira.oracle.QueryExchangeRatesRequest"></a>

 ### QueryExchangeRatesRequest
 QueryExchangeRatesRequest is the request type for the Query/ExchangeRates RPC method.

 

 

 
 <a name="kujira.oracle.QueryExchangeRatesResponse"></a>

 ### QueryExchangeRatesResponse
 QueryExchangeRatesResponse is response type for the
Query/ExchangeRates RPC method.

 
 | Field | Type | Label | Description |
 | ----- | ---- | ----- | ----------- |
 | `exchange_rates` | [cosmos.base.v1beta1.DecCoin](#cosmos.base.v1beta1.DecCoin) | repeated | exchange_rates defines a list of the exchange rate for all whitelisted denoms. |
 
 

 

 
 <a name="kujira.oracle.QueryMissCounterRequest"></a>

 ### QueryMissCounterRequest
 QueryMissCounterRequest is the request type for the Query/MissCounter RPC method.

 
 | Field | Type | Label | Description |
 | ----- | ---- | ----- | ----------- |
 | `validator_addr` | [string](#string) |  | validator defines the validator address to query for. |
 
 

 

 
 <a name="kujira.oracle.QueryMissCounterResponse"></a>

 ### QueryMissCounterResponse
 QueryMissCounterResponse is response type for the
Query/MissCounter RPC method.

 
 | Field | Type | Label | Description |
 | ----- | ---- | ----- | ----------- |
 | `miss_counter` | [uint64](#uint64) |  | miss_counter defines the oracle miss counter of a validator |
 
 

 

 
 <a name="kujira.oracle.QueryParamsRequest"></a>

 ### QueryParamsRequest
 QueryParamsRequest is the request type for the Query/Params RPC method.

 

 

 
 <a name="kujira.oracle.QueryParamsResponse"></a>

 ### QueryParamsResponse
 QueryParamsResponse is the response type for the Query/Params RPC method.

 
 | Field | Type | Label | Description |
 | ----- | ---- | ----- | ----------- |
 | `params` | [Params](#kujira.oracle.Params) |  | params defines the parameters of the module. |
 
 

 

 
 <a name="kujira.oracle.QueryVoteTargetsRequest"></a>

 ### QueryVoteTargetsRequest
 QueryVoteTargetsRequest is the request type for the Query/VoteTargets RPC method.

 

 

 
 <a name="kujira.oracle.QueryVoteTargetsResponse"></a>

 ### QueryVoteTargetsResponse
 QueryVoteTargetsResponse is response type for the
Query/VoteTargets RPC method.

 
 | Field | Type | Label | Description |
 | ----- | ---- | ----- | ----------- |
 | `vote_targets` | [string](#string) | repeated | vote_targets defines a list of the denomination in which everyone should vote in the current vote period. |
 
 

 

  <!-- end messages -->

  <!-- end enums -->

  <!-- end HasExtensions -->

 
 <a name="kujira.oracle.Query"></a>

 ### Query
 Query defines the gRPC querier service.

 | Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
 | ----------- | ------------ | ------------- | ------------| ------- | -------- |
 | `ExchangeRate` | [QueryExchangeRateRequest](#kujira.oracle.QueryExchangeRateRequest) | [QueryExchangeRateResponse](#kujira.oracle.QueryExchangeRateResponse) | ExchangeRate returns exchange rate of a denom | GET|/oracle/denoms/{denom}/exchange_rate|
 | `ExchangeRates` | [QueryExchangeRatesRequest](#kujira.oracle.QueryExchangeRatesRequest) | [QueryExchangeRatesResponse](#kujira.oracle.QueryExchangeRatesResponse) | ExchangeRates returns exchange rates of all denoms | GET|/oracle/denoms/exchange_rates|
 | `Actives` | [QueryActivesRequest](#kujira.oracle.QueryActivesRequest) | [QueryActivesResponse](#kujira.oracle.QueryActivesResponse) | Actives returns all active denoms | GET|/oracle/denoms/actives|
 | `MissCounter` | [QueryMissCounterRequest](#kujira.oracle.QueryMissCounterRequest) | [QueryMissCounterResponse](#kujira.oracle.QueryMissCounterResponse) | MissCounter returns oracle miss counter of a validator | GET|/oracle/validators/{validator_addr}/miss|
 | `Params` | [QueryParamsRequest](#kujira.oracle.QueryParamsRequest) | [QueryParamsResponse](#kujira.oracle.QueryParamsResponse) | Params queries all parameters. | GET|/oracle/params|
 
  <!-- end services -->

 
 
 <a name="kujira/oracle/tx.proto"></a>
 <p align="right"><a href="#top">Top</a></p>

 ## kujira/oracle/tx.proto
 

 
 <a name="kujira.oracle.MsgAddRequiredDenom"></a>

 ### MsgAddRequiredDenom
 MsgAddRequiredDenom represents a message to add a denom to the whitelist

 
 | Field | Type | Label | Description |
 | ----- | ---- | ----- | ----------- |
 | `authority` | [string](#string) |  |  |
 | `symbol` | [string](#string) |  |  |
 
 

 

 
 <a name="kujira.oracle.MsgAddRequiredDenomResponse"></a>

 ### MsgAddRequiredDenomResponse
 MsgAddRequiredDenomResponse defines the Msg/AddRequiredDenom response type.

 

 

 
 <a name="kujira.oracle.MsgRemoveRequiredDenom"></a>

 ### MsgRemoveRequiredDenom
 MsgRemoveRequiredDenom represents a message to remove a denom from the whitelist

 
 | Field | Type | Label | Description |
 | ----- | ---- | ----- | ----------- |
 | `authority` | [string](#string) |  |  |
 | `symbol` | [string](#string) |  |  |
 
 

 

 
 <a name="kujira.oracle.MsgRemoveRequiredDenomResponse"></a>

 ### MsgRemoveRequiredDenomResponse
 MsgRemoveRequiredDenomResponse defines the Msg/RemoveRequiredDenom response type.

 

 

 
 <a name="kujira.oracle.MsgUpdateParams"></a>

 ### MsgUpdateParams
 

 
 | Field | Type | Label | Description |
 | ----- | ---- | ----- | ----------- |
 | `authority` | [string](#string) |  |  |
 | `params` | [Params](#kujira.oracle.Params) |  |  |
 
 

 

 
 <a name="kujira.oracle.MsgUpdateParamsResponse"></a>

 ### MsgUpdateParamsResponse
 MsgUpdateParamsResponse defines the Msg/UpdateParams response type.

 

 

  <!-- end messages -->

  <!-- end enums -->

  <!-- end HasExtensions -->

 
 <a name="kujira.oracle.Msg"></a>

 ### Msg
 Msg defines the oracle Msg service.

 | Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
 | ----------- | ------------ | ------------- | ------------| ------- | -------- |
 | `AddRequiredDenom` | [MsgAddRequiredDenom](#kujira.oracle.MsgAddRequiredDenom) | [MsgAddRequiredDenomResponse](#kujira.oracle.MsgAddRequiredDenomResponse) | AddRequiredDenom adds a new price to the required list of prices | |
 | `RemoveRequiredDenom` | [MsgRemoveRequiredDenom](#kujira.oracle.MsgRemoveRequiredDenom) | [MsgRemoveRequiredDenomResponse](#kujira.oracle.MsgRemoveRequiredDenomResponse) | RemoveRequiredDenom removes a price from the required list of prices | |
 | `UpdateParams` | [MsgUpdateParams](#kujira.oracle.MsgUpdateParams) | [MsgUpdateParamsResponse](#kujira.oracle.MsgUpdateParamsResponse) | UpdateParams sets new module params | |
 
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
 