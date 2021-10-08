package genesistool

import (
	"bytes"
	"github.com/bif/bif-sdk-go/utils"
	"github.com/bif/bif-sdk-go/utils/hexutil"
	"github.com/bif/bif-sdk-go/utils/rlp"
)

var (
	istanbulExtraVanity = 32 // Fixed number of extra-data bytes reserved for validator vanity
	istanbulExtraSeal   = 98 // 33 + 1 + 64

	hotStuffExtraVanity = 32        // Fixed number of extra-data bytes reserved for validator vanity
	hotStuffExtraSeal   = 98        // Fixed number of extra-data bytes reserved for validator seal
)

type istanbulExtra struct {
	Validators    []utils.Address
	Seal          []byte
	CommittedSeal [][]byte
}

func EncodeIstanbul(vanity string, validators []string) (string, error) {
	newVanity, err := hexutil.Decode(vanity)
	if err != nil {
		return "", err
	}

	if len(newVanity) < istanbulExtraVanity {
		newVanity = append(newVanity, bytes.Repeat([]byte{0x00}, istanbulExtraVanity-len(newVanity))...)
	}
	newVanity = newVanity[:istanbulExtraVanity]

	vals := make([]utils.Address, 0, len(validators))
	for _, val := range validators {
		vals = append(vals, utils.StringToAddress(val))
	}

	ist := &istanbulExtra{
		Validators:    vals,
		Seal:          make([]byte, istanbulExtraSeal),
		CommittedSeal: [][]byte{},
	}

	payload, err := rlp.EncodeToBytes(&ist)
	if err != nil {
		return "", err
	}

	return "0x" + utils.Bytes2Hex(append(newVanity, payload...)), nil
}

type hotStuffExtra struct {
	SpeakerAddr utils.Address

	Validators []utils.Address // This only exists in Genesis as it can be too long for every block

	Mask          []byte
	AggregatedKey []byte
	AggregatedSig []byte

	Seal []byte
}

func EncodeHotStuff(vanity string, validators []string) (string, error) {
	newVanity, err := hexutil.Decode(vanity)
	if err != nil {
		return "", err
	}

	if len(newVanity) < hotStuffExtraVanity {
		newVanity = append(newVanity, bytes.Repeat([]byte{0x00}, hotStuffExtraVanity-len(newVanity))...)
	}
	newVanity = newVanity[:hotStuffExtraVanity]

	vals := make([]utils.Address, 0, len(validators))
	for _, val := range validators {
		vals = append(vals, utils.StringToAddress(val))
	}

	ist := &hotStuffExtra{
		SpeakerAddr:   utils.Address{},
		Validators:    vals,
		Mask:          nil,
		AggregatedKey: nil,
		AggregatedSig: nil,
		Seal:          make([]byte, hotStuffExtraSeal),
	}

	payload, err := rlp.EncodeToBytes(&ist)
	if err != nil {
		return "", err
	}

	return "0x" + utils.Bytes2Hex(append(newVanity, payload...)), nil
}