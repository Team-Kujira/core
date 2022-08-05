<!--
order: 2
-->

# State

## ExchangeRatePrevote

`ExchangeRatePrevote` containing validator voter's prevote for a given denom for the current `VotePeriod`.

- ExchangeRatePrevote: `0x01<denom_Bytes><valAddress_Bytes> -> amino(ExchangeRatePrevote)`

```go
type ValAddress []byte
type VoteHash []byte

type ExchangeRatePrevote struct {
	Hash        VoteHash       // Vote hex hash to protect centralize data source problem
	Denom       string         // Ticker name of native denom
	Voter       sdk.ValAddress // Voter val address
	SubmitBlock int64
}
```

## ExchangeRateVote

`ExchangeRateVote` containing validator voter's vote for a given denom for the current `VotePeriod`.

- ExchangeRateVote: `0x02<denom_Bytes><valAddress_Bytes> -> amino(ExchangeRateVote)`

```go
type ExchangeRateVote struct {
	ExchangeRate sdk.Dec        // ExchangeRate of denom in USD
	Denom        string         // Ticker name of native denom
	Voter        sdk.ValAddress // voter val address of validator
}
```

## ExchangeRate

An `sdk.Dec` that stores the current exchange rate for a given denom.

You can get the active list of denoms (denominations with votes past `VoteThreshold`) with `k.GetActiveDenoms()`.

- ExchangeRate: `0x03<denom_Bytes> -> amino(sdk.Dec)`

## FeederDelegation

An `sdk.AccAddress` (`kujira-` account) address of `operator`'s delegated price feeder.

- FeederDelegation: `0x04<valAddress_Bytes> -> amino(sdk.AccAddress)`

## MissCounter

An `int64` representing the number of `VotePeriods` that validator `operator` missed during the current `SlashWindow`.

- MissCounter: `0x05<valAddress_Bytes> -> amino(int64)`

## AggregateExchangeRatePrevote

`AggregateExchangeRatePrevote` containing validator voter's aggregated prevote for all denoms for the current `VotePeriod`.

- AggregateExchangeRatePrevote: `0x06<valAddress_Bytes> -> amino(AggregateExchangeRatePrevote)`

```go
// AggregateVoteHash is hash value to hide vote exchange rates
// which is formatted as hex string in SHA256("{salt}:{exchange rate}{denom},...,{exchange rate}{denom}:{voter}")
type AggregateVoteHash []byte

type AggregateExchangeRatePrevote struct {
	Hash        AggregateVoteHash // Vote hex hash to protect centralize data source problem
	Voter       sdk.ValAddress    // Voter val address
	SubmitBlock int64
}
```

## AggregateExchangeRateVote

`AggregateExchangeRateVote` containing validator voter's aggregate vote for all denoms for the current `VotePeriod`.

- AggregateExchangeRateVote: `0x07<valAddress_Bytes> -> amino(AggregateExchangeRateVote)`

```go
type ExchangeRateTuple struct {
	Denom        string  `json:"denom"`
	ExchangeRate sdk.Dec `json:"exchange_rate"`
}

type ExchangeRateTuples []ExchangeRateTuple

type AggregateExchangeRateVote struct {
	ExchangeRateTuples ExchangeRateTuples // ExchangeRates of Luna in target fiat currencies
	Voter              sdk.ValAddress     // voter val address of validator
}
```
