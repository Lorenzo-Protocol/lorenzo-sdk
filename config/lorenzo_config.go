package config

import (
	"fmt"
	"net/url"
	"time"

	"github.com/cosmos/relayer/v2/relayer/chains/cosmos"
)

// LorenzoConfig defines configuration for the Lorenzo client
// adapted from https://github.com/strangelove-ventures/lens/blob/v0.5.1/client/config.go
type LorenzoConfig struct {
	Key            string        `mapstructure:"key" toml:"key"`
	ChainID        string        `mapstructure:"chain-id" toml:"chain-id"`
	RPCAddr        string        `mapstructure:"rpc-addr" toml:"rpc-addr"`
	AccountPrefix  string        `mapstructure:"account-prefix" toml:"account-prefix"`
	KeyringBackend string        `mapstructure:"keyring-backend" toml:"keyring-backend"`
	GasAdjustment  float64       `mapstructure:"gas-adjustment" toml:"gas-adjustment"`
	GasPrices      string        `mapstructure:"gas-prices" toml:"gas-prices"`
	KeyDirectory   string        `mapstructure:"key-directory" toml:"key-directory"`
	Debug          bool          `mapstructure:"debug" toml:"debug"`
	Timeout        time.Duration `mapstructure:"timeout" toml:"timeout"`
	BlockTimeout   time.Duration `mapstructure:"block-timeout" toml:"block-timeout"`
	OutputFormat   string        `mapstructure:"output-format" toml:"output-format"`
	SignModeStr    string        `mapstructure:"sign-mode" toml:"sign-mode"`
}

func (cfg *LorenzoConfig) Validate() error {
	if _, err := url.Parse(cfg.RPCAddr); err != nil {
		return fmt.Errorf("rpc-addr is not correctly formatted: %w", err)
	}
	if cfg.Timeout <= 0 {
		return fmt.Errorf("timeout must be positive")
	}
	if cfg.BlockTimeout < 0 {
		return fmt.Errorf("block-timeout can't be negative")
	}
	return nil
}

func (cfg *LorenzoConfig) ToCosmosProviderConfig() cosmos.CosmosProviderConfig {
	return cosmos.CosmosProviderConfig{
		Key:            cfg.Key,
		ChainID:        cfg.ChainID,
		RPCAddr:        cfg.RPCAddr,
		AccountPrefix:  cfg.AccountPrefix,
		KeyringBackend: cfg.KeyringBackend,
		GasAdjustment:  cfg.GasAdjustment,
		GasPrices:      cfg.GasPrices,
		KeyDirectory:   cfg.KeyDirectory,
		Debug:          cfg.Debug,
		Timeout:        cfg.Timeout.String(),
		BlockTimeout:   cfg.BlockTimeout.String(),
		OutputFormat:   cfg.OutputFormat,
		SignModeStr:    cfg.SignModeStr,
	}
}
