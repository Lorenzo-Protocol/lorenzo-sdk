package config_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Lorenzo-Protocol/lorenzo-sdk/config"
)

// TestLorenzoQueryConfig ensures that the default Lorenzo query config is valid
func TestLorenzoQueryConfig(t *testing.T) {
	defaultConfig := config.DefaultLorenzoQueryConfig()
	err := defaultConfig.Validate()
	require.NoError(t, err)
}
