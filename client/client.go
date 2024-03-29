package client

import (
	"context"
	"time"

	"github.com/Lorenzo-Protocol/lorenzo-sdk/config"
	"github.com/Lorenzo-Protocol/lorenzo-sdk/query"
	lorenzo "github.com/Lorenzo-Protocol/lorenzo/app"
	"github.com/cosmos/relayer/v2/relayer/chains/cosmos"
	"go.uber.org/zap"
)

type Client struct {
	*query.QueryClient

	provider *cosmos.CosmosProvider
	timeout  time.Duration
	logger   *zap.Logger
	cfg      *config.LorenzoConfig
}

func New(cfg *config.LorenzoConfig, logger *zap.Logger) (*Client, error) {
	var (
		zapLogger *zap.Logger
		err       error
	)

	// ensure cfg is valid
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	// use the existing logger or create a new one if not given
	zapLogger = logger
	if zapLogger == nil {
		zapLogger, err = newRootLogger("console", true)
		if err != nil {
			return nil, err
		}
	}

	// Create tmp Lorenzo app to retrieve and register codecs
	encCfg := lorenzo.MakeEncodingConfig()

	cosmosConfig := cfg.ToCosmosProviderConfig()
	provider, err := cosmosConfig.NewProvider(
		zapLogger,
		"", // TODO: set home path
		true,
		"lorenzo",
	)
	if err != nil {
		return nil, err
	}

	cp := provider.(*cosmos.CosmosProvider)
	cp.PCfg.KeyDirectory = cfg.KeyDirectory
	// Need to override this manually as otherwise option from config is ignored
	cp.Cdc = cosmos.Codec{
		InterfaceRegistry: encCfg.InterfaceRegistry,
		Marshaler:         encCfg.Codec,
		TxConfig:          encCfg.TxConfig,
		Amino:             encCfg.Amino,
	}

	err = cp.Init(context.Background())
	if err != nil {
		return nil, err
	}

	// create a queryClient so that the Client inherits all query functions
	queryClient, err := query.NewWithClient(cp.RPCClient, cfg.Timeout)
	if err != nil {
		return nil, err
	}

	return &Client{
		queryClient,
		cp,
		cfg.Timeout,
		zapLogger,
		cfg,
	}, nil
}

func (c *Client) GetConfig() *config.LorenzoConfig {
	return c.cfg
}

func (c *Client) Stop() error {
	if !c.provider.RPCClient.IsRunning() {
		return nil
	}

	return c.provider.RPCClient.Stop()
}
