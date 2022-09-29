package genesistool

import (
	"encoding/hex"
	"github.com/tchain/go-tchain-sdk/utils"
	"strings"
)

// RegEncode regulators encode
func RegEncode(regulators []string) string {
	addrs := make([]utils.Address, 0, len(regulators))
	for _, val := range regulators {
		addrs = append(addrs, utils.StringToAddress(val))
	}

	regs := make([]string, utils.AddressLength*len(addrs))
	for idx, addr := range addrs {
		regs[idx] = hex.EncodeToString(addr.Bytes())
	}
	return "0x" + strings.Join(regs, "")
}
