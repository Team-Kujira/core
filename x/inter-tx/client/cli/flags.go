package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	// The connection end identifier on the controller chain
	FlagConnectionID = "connection-id"
	// The controller chain channel version
	FlagVersion = "version"
	// The account id for this interchain account
	FlagAccountID = "account-id"
	// The memo of a submitted ICA transaction
	FlagICAMemo = "ica-memo"
	// The timeout of a submitted ICA transaction, nanoseconds
	FlagICATimeout = "ica-timeout"
)

// common flagsets to add to various functions
var (
	fsConnectionID = flag.NewFlagSet("", flag.ContinueOnError)
	fsVersion      = flag.NewFlagSet("", flag.ContinueOnError)
	fsAccountID    = flag.NewFlagSet("", flag.ContinueOnError)
	fsICAMemo      = flag.NewFlagSet("", flag.ContinueOnError)
	fsICATimeout   = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	fsConnectionID.String(FlagConnectionID, "", "Connection ID")
	fsVersion.String(FlagVersion, "", "Version")
	fsAccountID.String(FlagAccountID, "", "Account ID")
	fsICAMemo.String(FlagICAMemo, "", "ICA Memo")
	fsICATimeout.Uint64(FlagICATimeout, 0, "ICA Timeout (ns)")
}
