package configutils

import (
	"coreumservice/go/stably_io/utils"

	"github.com/pkg/errors"
)

// Use *ValueFunc instead of value to only fetch the value of
// the target stage.
//
// Each stage's value function may require some external execution
// thus it should only be invoked if the corresponding stage is selected.
func SwitchOnStage[T any](stage utils.Stage,
	prodValueFunc func() T,
	betaValueFunc func() T,
	testValueFunc func() T,
	localValueFunc func() T,
) T {
	switch stage {
	case utils.Prod:
		return prodValueFunc()
	case utils.Beta:
		return betaValueFunc()
	case utils.Test:
		return testValueFunc()
	case utils.Local:
		return localValueFunc()
	default:
		panic(errors.Errorf("Not supported stage: %v", stage))
	}
}
