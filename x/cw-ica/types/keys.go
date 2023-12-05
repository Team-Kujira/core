package types

import (
	"fmt"
)

const (
	ModuleName = "cwica"

	StoreKey = ModuleName

	RouterKey = ModuleName

	QuerierRoute = ModuleName
)

const (
	// CallbackDataKeyPrefix is the prefix of CallbackData
	CallbackDataKeyPrefix = "CallbackData/"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

func PacketID(portID string, channelID string, sequence uint64) string {
	return fmt.Sprintf("%s.%s.%d", portID, channelID, sequence)
}

// CallbackDataKey returns the store key to retrieve a CallbackData from the index fields
func CallbackDataKey(callbackKey string) []byte {
	var key []byte

	callbackKeyBytes := []byte(callbackKey)
	key = append(key, callbackKeyBytes...)
	key = append(key, []byte("/")...)

	return key
}
