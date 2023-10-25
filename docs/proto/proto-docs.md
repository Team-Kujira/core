<!-- This file is auto-generated. Please do not modify it yourself. -->
 # Protobuf Documentation
 <a name="top"></a>

 ## Table of Contents
 
 - [kujira/scheduler/params.proto](#kujira/scheduler/params.proto)
     - [Params](#kujira.scheduler.Params)
   
 - [kujira/scheduler/hook.proto](#kujira/scheduler/hook.proto)
     - [Hook](#kujira.scheduler.Hook)
   
 - [kujira/scheduler/genesis.proto](#kujira/scheduler/genesis.proto)
     - [GenesisState](#kujira.scheduler.GenesisState)
   
 - [kujira/scheduler/query.proto](#kujira/scheduler/query.proto)
     - [QueryAllHookRequest](#kujira.scheduler.QueryAllHookRequest)
     - [QueryAllHookResponse](#kujira.scheduler.QueryAllHookResponse)
     - [QueryGetHookRequest](#kujira.scheduler.QueryGetHookRequest)
     - [QueryGetHookResponse](#kujira.scheduler.QueryGetHookResponse)
     - [QueryParamsRequest](#kujira.scheduler.QueryParamsRequest)
     - [QueryParamsResponse](#kujira.scheduler.QueryParamsResponse)
   
     - [Query](#kujira.scheduler.Query)
   
 - [kujira/scheduler/tx.proto](#kujira/scheduler/tx.proto)
     - [MsgCreateHook](#kujira.scheduler.MsgCreateHook)
     - [MsgCreateHookResponse](#kujira.scheduler.MsgCreateHookResponse)
     - [MsgDeleteHook](#kujira.scheduler.MsgDeleteHook)
     - [MsgDeleteHookResponse](#kujira.scheduler.MsgDeleteHookResponse)
     - [MsgUpdateHook](#kujira.scheduler.MsgUpdateHook)
     - [MsgUpdateHookResponse](#kujira.scheduler.MsgUpdateHookResponse)
   
     - [Msg](#kujira.scheduler.Msg)
   
 - [Scalar Value Types](#scalar-value-types)

 
 
 <a name="kujira/scheduler/params.proto"></a>
 <p align="right"><a href="#top">Top</a></p>

 ## kujira/scheduler/params.proto
 

 
 <a name="kujira.scheduler.Params"></a>

 ### Params
 Params defines the parameters for the module.

 

 

  <!-- end messages -->

  <!-- end enums -->

  <!-- end HasExtensions -->

  <!-- end services -->

 
 
 <a name="kujira/scheduler/hook.proto"></a>
 <p align="right"><a href="#top">Top</a></p>

 ## kujira/scheduler/hook.proto
 

 
 <a name="kujira.scheduler.Hook"></a>

 ### Hook
 

 
 | Field | Type | Label | Description |
 | ----- | ---- | ----- | ----------- |
 | `id` | [uint64](#uint64) |  |  |
 | `executor` | [string](#string) |  |  |
 | `contract` | [string](#string) |  |  |
 | `msg` | [bytes](#bytes) |  |  |
 | `frequency` | [int64](#int64) |  |  |
 | `funds` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |
 
 

 

  <!-- end messages -->

  <!-- end enums -->

  <!-- end HasExtensions -->

  <!-- end services -->

 
 
 <a name="kujira/scheduler/genesis.proto"></a>
 <p align="right"><a href="#top">Top</a></p>

 ## kujira/scheduler/genesis.proto
 

 
 <a name="kujira.scheduler.GenesisState"></a>

 ### GenesisState
 GenesisState defines the scheduler module's genesis state.

 
 | Field | Type | Label | Description |
 | ----- | ---- | ----- | ----------- |
 | `params` | [Params](#kujira.scheduler.Params) |  |  |
 | `hookList` | [Hook](#kujira.scheduler.Hook) | repeated |  |
 | `hookCount` | [uint64](#uint64) |  |  |
 
 

 

  <!-- end messages -->

  <!-- end enums -->

  <!-- end HasExtensions -->

  <!-- end services -->

 
 
 <a name="kujira/scheduler/query.proto"></a>
 <p align="right"><a href="#top">Top</a></p>

 ## kujira/scheduler/query.proto
 

 
 <a name="kujira.scheduler.QueryAllHookRequest"></a>

 ### QueryAllHookRequest
 

 
 | Field | Type | Label | Description |
 | ----- | ---- | ----- | ----------- |
 | `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  |  |
 
 

 

 
 <a name="kujira.scheduler.QueryAllHookResponse"></a>

 ### QueryAllHookResponse
 

 
 | Field | Type | Label | Description |
 | ----- | ---- | ----- | ----------- |
 | `Hook` | [Hook](#kujira.scheduler.Hook) | repeated |  |
 | `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  |  |
 
 

 

 
 <a name="kujira.scheduler.QueryGetHookRequest"></a>

 ### QueryGetHookRequest
 

 
 | Field | Type | Label | Description |
 | ----- | ---- | ----- | ----------- |
 | `id` | [uint64](#uint64) |  |  |
 
 

 

 
 <a name="kujira.scheduler.QueryGetHookResponse"></a>

 ### QueryGetHookResponse
 

 
 | Field | Type | Label | Description |
 | ----- | ---- | ----- | ----------- |
 | `Hook` | [Hook](#kujira.scheduler.Hook) |  |  |
 
 

 

 
 <a name="kujira.scheduler.QueryParamsRequest"></a>

 ### QueryParamsRequest
 QueryParamsRequest is request type for the Query/Params RPC method.

 

 

 
 <a name="kujira.scheduler.QueryParamsResponse"></a>

 ### QueryParamsResponse
 QueryParamsResponse is response type for the Query/Params RPC method.

 
 | Field | Type | Label | Description |
 | ----- | ---- | ----- | ----------- |
 | `params` | [Params](#kujira.scheduler.Params) |  | params holds all the parameters of this module. |
 
 

 

  <!-- end messages -->

  <!-- end enums -->

  <!-- end HasExtensions -->

 
 <a name="kujira.scheduler.Query"></a>

 ### Query
 Query defines the gRPC querier service.

 | Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
 | ----------- | ------------ | ------------- | ------------| ------- | -------- |
 | `Params` | [QueryParamsRequest](#kujira.scheduler.QueryParamsRequest) | [QueryParamsResponse](#kujira.scheduler.QueryParamsResponse) | Parameters queries the parameters of the module. | GET|/kujira/scheduler/params|
 | `Hook` | [QueryGetHookRequest](#kujira.scheduler.QueryGetHookRequest) | [QueryGetHookResponse](#kujira.scheduler.QueryGetHookResponse) | Queries a Hook by id. | GET|/kujira/scheduler/hook/{id}|
 | `HookAll` | [QueryAllHookRequest](#kujira.scheduler.QueryAllHookRequest) | [QueryAllHookResponse](#kujira.scheduler.QueryAllHookResponse) | Queries a list of Hook items. | GET|/kujira/scheduler/hook|
 
  <!-- end services -->

 
 
 <a name="kujira/scheduler/tx.proto"></a>
 <p align="right"><a href="#top">Top</a></p>

 ## kujira/scheduler/tx.proto
 

 
 <a name="kujira.scheduler.MsgCreateHook"></a>

 ### MsgCreateHook
 

 
 | Field | Type | Label | Description |
 | ----- | ---- | ----- | ----------- |
 | `authority` | [string](#string) |  |  |
 | `executor` | [string](#string) |  | The account that will execute the msg on the schedule |
 | `contract` | [string](#string) |  | The contract that the msg is called on |
 | `msg` | [bytes](#bytes) |  |  |
 | `frequency` | [int64](#int64) |  |  |
 | `funds` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |
 
 

 

 
 <a name="kujira.scheduler.MsgCreateHookResponse"></a>

 ### MsgCreateHookResponse
 

 

 

 
 <a name="kujira.scheduler.MsgDeleteHook"></a>

 ### MsgDeleteHook
 

 
 | Field | Type | Label | Description |
 | ----- | ---- | ----- | ----------- |
 | `authority` | [string](#string) |  |  |
 | `id` | [uint64](#uint64) |  |  |
 
 

 

 
 <a name="kujira.scheduler.MsgDeleteHookResponse"></a>

 ### MsgDeleteHookResponse
 

 

 

 
 <a name="kujira.scheduler.MsgUpdateHook"></a>

 ### MsgUpdateHook
 

 
 | Field | Type | Label | Description |
 | ----- | ---- | ----- | ----------- |
 | `authority` | [string](#string) |  |  |
 | `id` | [uint64](#uint64) |  |  |
 | `executor` | [string](#string) |  |  |
 | `contract` | [string](#string) |  |  |
 | `msg` | [bytes](#bytes) |  |  |
 | `frequency` | [int64](#int64) |  |  |
 | `funds` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |
 
 

 

 
 <a name="kujira.scheduler.MsgUpdateHookResponse"></a>

 ### MsgUpdateHookResponse
 

 

 

  <!-- end messages -->

  <!-- end enums -->

  <!-- end HasExtensions -->

 
 <a name="kujira.scheduler.Msg"></a>

 ### Msg
 Msg defines the scheduler Msg service.

 | Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
 | ----------- | ------------ | ------------- | ------------| ------- | -------- |
 | `CreateHook` | [MsgCreateHook](#kujira.scheduler.MsgCreateHook) | [MsgCreateHookResponse](#kujira.scheduler.MsgCreateHookResponse) | CreateHook adds a new hook to the scheduler | |
 | `UpdateHook` | [MsgUpdateHook](#kujira.scheduler.MsgUpdateHook) | [MsgUpdateHookResponse](#kujira.scheduler.MsgUpdateHookResponse) | UpdateHook updates an existing hook | |
 | `DeleteHook` | [MsgDeleteHook](#kujira.scheduler.MsgDeleteHook) | [MsgDeleteHookResponse](#kujira.scheduler.MsgDeleteHookResponse) | DeleteHook removes a hook from the scheduler | |
 
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
 