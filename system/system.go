package system

import "github.com/bif/bif-sdk-go/providers"

type System struct {
	Dpos     *Dpos
	//Voter    *Dpos
	//Cer      *Dpos
}

func NewSystem(provider providers.ProviderInterface) *System {
	system := new(System)
	system.Dpos = NewDpos(provider)
	return system
}