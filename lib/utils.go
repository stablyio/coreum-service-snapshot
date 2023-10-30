package coreumservicelib

import (
	"coreumservice/go/stably_io/utils"
	"encoding/json"

	coreumconstant "github.com/CoreumFoundation/coreum/pkg/config/constant"
)

func ToJSONPretty(input interface{}) string {
	res, err := json.MarshalIndent(input, "", "    ")
	if err != nil {
		panic(err)
	}
	return string(res)
}

func GetChainIDByStage() string {
	chainID := ""
	switch utils.GetStage() {
	case utils.Prod:
		chainID = string(coreumconstant.ChainIDMain)
	case utils.Beta, utils.Local, utils.Test:
		chainID = string(coreumconstant.ChainIDTest)
	default:
		break
	}
	return chainID
}