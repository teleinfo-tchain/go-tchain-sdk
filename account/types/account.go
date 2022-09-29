package types

import (
	"github.com/tchain/go-tchain-sdk/utils"
)

type Account struct {
	Address utils.Address `json:"address"` // Ethereum account address derived from the key
	URL     URL           `json:"url"`     // Optional resource locator within a backend
}
