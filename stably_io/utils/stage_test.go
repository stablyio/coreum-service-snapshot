package utils_test

import (
	"coreumservice/go/stably_io/utils"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_GetStage(t *testing.T) {
	t.Run("Valid cases", func(t *testing.T) {
		t.Setenv("STAGE", "prod")
		require.Equal(t, utils.Prod, utils.GetStage())

		t.Setenv("STAGE", "beta")
		require.Equal(t, utils.Beta, utils.GetStage())

		t.Setenv("STAGE", "local")
		require.Equal(t, utils.Local, utils.GetStage())

		t.Setenv("STAGE", "test")
		require.Equal(t, utils.Test, utils.GetStage())
	})

	t.Run("Invalid cases", func(t *testing.T) {
		t.Setenv("STAGE", "invalid")
		require.Panics(t, func() {
			utils.GetStage()
		})

		t.Setenv("STAGE", "ci")
		require.Panics(t, func() {
			utils.GetStage()
		})

		t.Setenv("STAGE", " ")
		require.Panics(t, func() {
			utils.GetStage()
		})

		t.Setenv("STAGE", "")
		require.Panics(t, func() {
			utils.GetStage()
		})
	})
}

func Test_Stage_String(t *testing.T) {
	require.Equal(t, "prod", utils.Prod.String())
	require.Equal(t, "beta", utils.Beta.String())
	require.Equal(t, "local", utils.Local.String())
	require.Equal(t, "test", utils.Test.String())
}
