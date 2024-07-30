package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Lorenzo-Protocol/lorenzo-sdk/v2/client"
	"github.com/Lorenzo-Protocol/lorenzo-sdk/v2/config"
)

func main() {
	key := os.Getenv("KEY")
	keyDir := os.Getenv("KEY_DIR")
	chainId := os.Getenv("CHAIN_ID")
	conf := createConfig(keyDir, key, chainId)
	lorenzoClient, err := client.New(&conf, nil)
	if err != nil {
		panic(err)
	}

	fmt.Println(lorenzoClient.MustGetAddr())
}

func createConfig(keyDir string, key string, chainId string) config.LorenzoConfig {
	return config.LorenzoConfig{
		Key:            key,
		ChainID:        chainId,
		RPCAddr:        "http://localhost:26657",
		AccountPrefix:  "lrz",
		KeyringBackend: "test",
		GasAdjustment:  1.2,
		GasPrices:      "2ulrz",
		KeyDirectory:   keyDir,
		Debug:          true,
		Timeout:        20 * time.Second,
		OutputFormat:   "json",
		SignModeStr:    "direct",
	}
}
