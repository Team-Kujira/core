package abci_test

import (
	"encoding/json"
	"testing"

	"github.com/Team-Kujira/core/x/oracle/abci"
	"github.com/stretchr/testify/require"
)

func TestDecoding(t *testing.T) {
	resBody := []byte(`{"prices":{"BTC":"47375.706652541026694000","ETH":"2649.328939436595054949","USDT":"1.000661260343873178"}}`)
	prices := abci.PricesResponse{}
	err := json.Unmarshal(resBody, &prices)
	require.NoError(t, err)
}
