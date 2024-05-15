package event

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	abci_types "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethereum "github.com/ethereum/go-ethereum/common"
)

const (
	EventTypeMint = "lorenzo.btcstaking.v1.EventBTCStakingCreated"
	EventTypeBurn = "lorenzo.btcstaking.v1.EventBurnCreated"
)

type (
	BurnEvent struct {
		Amount           sdk.Coin `json:"amount"`
		BtcTargetAddress string   `json:"btc_target_address"`
		Signer           string   `json:"signer"`
	}

	MintEvent struct {
		TxHash     string `json:"tx_hash"`
		Amount     uint64 `json:"amount"`
		MintToAddr string `json:"mint_to_addr"`
	}

	MintRecordValue struct {
		TxHash     string `json:"tx_hash"`
		Amount     string `json:"amount"`
		MintToAddr string `json:"mint_to_addr"`
	}
)

func NewMintEvent(event abci_types.Event) (*MintEvent, error) {
	var record string
	for _, attr := range event.Attributes {
		if attr.Key == "record" {
			record = attr.Value
		}
	}
	if record == "" {
		return nil, errors.New("invalid mint event attributes, missing record key")
	}

	var value MintRecordValue
	err := json.Unmarshal([]byte(record), &value)
	if err != nil {
		return nil, err
	}
	txHashBytes, err := base64.StdEncoding.DecodeString(value.TxHash)
	if err != nil {
		return nil, err
	}
	txHash, err := chainhash.NewHash(txHashBytes)
	if err != nil {
		return nil, err
	}
	amount, err := strconv.ParseUint(value.Amount, 10, 64)
	if err != nil {
		return nil, err
	}
	mintToAddrBytes, err := base64.StdEncoding.DecodeString(value.MintToAddr)
	if err != nil {
		return nil, err
	}
	mintToAddr := ethereum.BytesToAddress(mintToAddrBytes)

	return &MintEvent{
		TxHash:     txHash.String(),
		Amount:     amount * 1e10,
		MintToAddr: mintToAddr.String(),
	}, nil
}

func NewBurnEvent(event abci_types.Event) (*BurnEvent, error) {
	var (
		amount           sdk.Coin
		btcTargetAddress string
		signer           string
		err              error
	)
	for _, attr := range event.Attributes {
		value := strings.Trim(attr.Value, "\"")
		if attr.Key == "amount" {
			amount, err = sdk.ParseCoinNormalized(value)
			if err != nil {
				return nil, err
			}
		}
		if attr.Key == "btc_target_address" {
			btcTargetAddress = value
		}
		if attr.Key == "signer" {
			address, err := sdk.AccAddressFromBech32(value)
			if err != nil {
				return nil, err
			}
			signer = ethereum.BytesToAddress(address.Bytes()).String()
		}
	}

	return &BurnEvent{
		Amount:           amount,
		BtcTargetAddress: btcTargetAddress,
		Signer:           signer,
	}, nil
}
