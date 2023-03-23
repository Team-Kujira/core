package cli

import (
	"time"

	"github.com/spf13/cobra"
	// "github.com/cosmos/cosmos-sdk/client/flags"
)

var DefaultRelativePacketTimeoutTimestamp = uint64((time.Duration(10) * time.Minute).Nanoseconds())

//nolint:unused
const (
	flagPacketTimeoutTimestamp = "packet-timeout-timestamp"
	listSeparator              = ","
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{}

	return cmd
}
