package config

import (
	"fmt"
	"github.com/cosmos/relayer/v2/relayer/chains/cosmos"
	"net/url"
	"time"
)

// LorenzoConfig defines configuration for the Lorenzo client
// adapted from https://github.com/strangelove-ventures/lens/blob/v0.5.1/client/config.go
type LorenzoConfig struct {
	Key            string        `mapstructure:"key"`
	ChainID        string        `mapstructure:"chain-id"`
	RPCAddr        string        `mapstructure:"rpc-addr"`
	AccountPrefix  string        `mapstructure:"account-prefix"`
	KeyringBackend string        `mapstructure:"keyring-backend"`
	GasAdjustment  float64       `mapstructure:"gas-adjustment"`
	GasPrices      string        `mapstructure:"gas-prices"`
	KeyDirectory   string        `mapstructure:"key-directory"`
	Debug          bool          `mapstructure:"debug"`
	Timeout        time.Duration `mapstructure:"timeout"`
	BlockTimeout   time.Duration `mapstructure:"block-timeout"`
	OutputFormat   string        `mapstructure:"output-format"`
	SignModeStr    string        `mapstructure:"sign-mode"`
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
