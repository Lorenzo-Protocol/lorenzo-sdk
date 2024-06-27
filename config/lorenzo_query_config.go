package config

import (
	"fmt"
	"net/url"
	"time"
)

// LorenzoConfig defines configuration for the Lorenzo query client
type LorenzoQueryConfig struct {
	RPCAddr string        `mapstructure:"rpc-addr" toml:"rpc_addr"`
	Timeout time.Duration `mapstructure:"timeout" toml:"rpc_timeout"`
}

func (cfg *LorenzoQueryConfig) Validate() error {
	if _, err := url.Parse(cfg.RPCAddr); err != nil {
		return fmt.Errorf("cfg.RPCAddr is not correctly formatted: %w", err)
	}
	if cfg.Timeout <= 0 {
		return fmt.Errorf("cfg.Timeout must be positive")
	}
	return nil
}

func DefaultLorenzoQueryConfig() LorenzoQueryConfig {
	return LorenzoQueryConfig{
		RPCAddr: "http://localhost:26657",
		Timeout: 20 * time.Second,
	}
}
