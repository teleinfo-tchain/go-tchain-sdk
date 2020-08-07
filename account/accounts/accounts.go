package accounts

import (
	"github.com/bif/bif-sdk-go/utils"
)

type Account struct {
	Address utils.Address `json:"address"` // Ethereum account address derived from the key
	URL     URL           `json:"url"`     // Optional resource locator within a backend
}
