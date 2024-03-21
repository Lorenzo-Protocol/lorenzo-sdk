package config_test

import (
	"testing"

	"github.com/Lorenzo-Protocol/rpc-client/config"
	"github.com/stretchr/testify/require"
)

// TestLorenzoQueryConfig ensures that the default Lorenzo query config is valid
func TestLorenzoQueryConfig(t *testing.T) {
	defaultConfig := config.DefaultLorenzoQueryConfig()
	err := defaultConfig.Validate()
	require.NoError(t, err)
}
