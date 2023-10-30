package utils_test

import (
	"testing"

	"coreumservice/go/stably_io/utils"

	"github.com/stretchr/testify/require"
)

func Test_GetEnvVar(t *testing.T) {
	t.Run("Non-empty case", func(it *testing.T) {
		// Pre-set the env var first
		it.Setenv("TESTING_ENV_VAR_543543", "testing_value")

		// Make sure the value is returned
		value := utils.GetEnvVar("TESTING_ENV_VAR_543543", true)
		require.Equal(it, "testing_value", value)
		value = utils.GetEnvVar("TESTING_ENV_VAR_543543", false)
		require.Equal(it, "testing_value", value)
	})

	t.Run("Empty case with panic", func(it *testing.T) {
		// Make sure it will panic
		defer func() {
			r := recover()
			require.NotNil(it, r)
		}()
		_ = utils.GetEnvVar("NOT_EXISTING_ENV_VAR_543543", true)
	})

	t.Run("Empty case without panic", func(it *testing.T) {
		value := utils.GetEnvVar("NOT_EXISTING_ENV_VAR_543543", false)
		require.Empty(it, value)
	})
}
