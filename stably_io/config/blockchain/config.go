package blockchainconfig

import (
	coreumconfig "coreumservice/go/stably_io/config/blockchain/coreum"
	"coreumservice/go/stably_io/utils"
)

type Blockchain struct {
	// Specific chain configs
	Coreum *coreumconfig.Coreum
}

func GetConfig(stage utils.Stage) *Blockchain {
	return &Blockchain{
		Coreum: coreumconfig.GetConfig(stage),
	}
}

func GetConfigDefault() *Blockchain {
	return GetConfig(utils.GetStage())
}
