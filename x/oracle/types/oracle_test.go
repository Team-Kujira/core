package types_test

import (
	"encoding/json"
	"testing"

	"github.com/Team-Kujira/core/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestVoteExtensionCompress(t *testing.T) {
	exchangeRatesStr := `[{"denom":"AKT","amount":"3.417096740560307550"},{"denom":"AMPKUJI","amount":"1.084075333349937652"},{"denom":"ARB","amount":"0.653568215050396212"},{"denom":"ATOM","amount":"5.790707106715531081"},{"denom":"AVAX","amount":"25.827959663657512637"},{"denom":"AXL","amount":"0.624201741733833997"},{"denom":"BAND","amount":"1.022584469058041796"},{"denom":"BNB","amount":"502.772950580502485398"},{"denom":"BTC","amount":"55783.302479024097639328"},{"denom":"CACAO","amount":"0.698803134648437975"},{"denom":"CRO","amount":"0.084405107433108864"},{"denom":"DOT","amount":"5.962045778513116897"},{"denom":"DYDX","amount":"1.269955883131320003"},{"denom":"DYM","amount":"1.335809901613850997"},{"denom":"ETH","amount":"2969.662966676288686015"},{"denom":"FET","amount":"1.139696855883139672"},{"denom":"FTM","amount":"0.428773762332588376"},{"denom":"FUZN","amount":"0.021324054049826744"},{"denom":"GLMR","amount":"0.192092486818396021"},{"denom":"HNT","amount":"3.087596429332989514"},{"denom":"INJ","amount":"19.452156583675507257"},{"denom":"JUNO","amount":"0.118688775145495572"},{"denom":"KAVA","amount":"0.374775225296475929"},{"denom":"KUJI","amount":"1.029789039653282812"},{"denom":"LINK","amount":"12.816128641011379658"},{"denom":"LUNA","amount":"0.364505753316675255"},{"denom":"LUNC","amount":"0.000068828162549657"},{"denom":"MATIC","amount":"0.491261448912043029"},{"denom":"MNTA","amount":"0.181188804107436880"},{"denom":"NSTK","amount":"0.038445180649993541"},{"denom":"NTRN","amount":"0.381527699701723989"},{"denom":"OSMO","amount":"0.471280360168268721"},{"denom":"PAXG","amount":"2335.525224879736093489"},{"denom":"QCKUJI","amount":"1.050603277885129415"},{"denom":"QCMNTA","amount":"0.194268552558550051"},{"denom":"RIO","amount":"1.059295830184573571"},{"denom":"RLB","amount":"0.078449612556799117"},{"denom":"SCRT","amount":"0.251650850228059606"},{"denom":"SHD","amount":"1.507075200077156108"},{"denom":"SOL","amount":"136.530800297842253396"},{"denom":"SOMM","amount":"0.031828093754901402"},{"denom":"STARS","amount":"0.010049914978090783"},{"denom":"STATOM","amount":"7.879848164467987846"},{"denom":"STKATOM","amount":"7.500369151595600714"},{"denom":"STOSMO","amount":"0.583630068136456088"},{"denom":"TIA","amount":"5.998201205777194726"},{"denom":"UNI","amount":"7.921914875907396220"},{"denom":"USDC","amount":"1.000020995501583086"},{"denom":"USDT","amount":"0.999735734500345704"},{"denom":"USK","amount":"0.988651090995423641"},{"denom":"WBTC","amount":"55811.491789331262939715"},{"denom":"WETH","amount":"2969.114861570696077782"},{"denom":"WHALE","amount":"0.005195497968885498"},{"denom":"WINK","amount":"0.029363420994439294"},{"denom":"WSTETH","amount":"3478.222037911750161926"}]`
	exchangeRates := sdk.DecCoins{}
	err := json.Unmarshal([]byte(exchangeRatesStr), &exchangeRates)
	require.NoError(t, err)

	prices := make(map[uint32][]byte)
	for index, rate := range exchangeRates {
		id := uint32(index)
		// id, err := types.SymbolToID(rate.Denom)
		// require.NoError(t, err)
		count := 0
		newAmount := rate.Amount
		for !newAmount.IsZero() {
			newAmount = newAmount.QuoInt64(10)
			count++
		}
		cuttingDecimals := count - 8
		for i := 0; i < cuttingDecimals; i++ {
			rate.Amount = rate.Amount.QuoInt64(10)
		}
		// for i := 0; i < cuttingDecimals; i++ {
		// 	rate.Amount = rate.Amount.MulInt64(10)
		// }
		prices[id] = append(rate.Amount.BigInt().Bytes(), byte(cuttingDecimals))
	}
	voteExt := types.VoteExtension{
		Height: 1,
		Prices: prices,
	}
	bz, err := voteExt.Marshal()
	require.NoError(t, err)
	_ = bz
	// require.Len(t, bz, 1241)
	// require.Len(t, bz, 776)

	compressed, err := voteExt.Compress()
	require.NoError(t, err)
	_ = compressed
	// require.Len(t, compressed, 1205)
	// require.Len(t, compressed, 733)
}
