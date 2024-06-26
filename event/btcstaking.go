package event

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"math/big"
	"strings"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	abci_types "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethereum "github.com/ethereum/go-ethereum/common"
)

const (
	EventTypeMint       = "lorenzo.btcstaking.v1.EventBTCStakingCreated"
	EventTypeBurn       = "lorenzo.btcstaking.v1.EventBurnCreated"
	Bech32PrefixAccAddr = "lrz"
)

type (
	BurnEvent struct {
		Amount           sdk.Coin `json:"amount"`
		BtcTargetAddress string   `json:"btc_target_address"`
		Signer           string   `json:"signer"`
	}

	MintEvent struct {
		TxHash          string  `json:"tx_hash"`
		Amount          big.Int `json:"amount"`
		MintToAddr      string  `json:"mint_to_addr"`
		BtcReceiverName string  `json:"btc_receiver_name"`
		BtcReceiverAddr string  `json:"btc_receiver_addr"`
	}

	MintRecordValue struct {
		TxHash          string `json:"tx_hash"`
		Amount          string `json:"amount"`
		MintToAddr      string `json:"mint_to_addr"`
		BtcReceiverName string `json:"btc_receiver_name"`
		BtcReceiverAddr string `json:"btc_receiver_addr"`
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

	amount, ok := new(big.Int).SetString(value.Amount, 10)
	if !ok {
		return nil, errors.New("parse mint event error: invalid amount")
	}

	mintToAddrBytes, err := base64.StdEncoding.DecodeString(value.MintToAddr)
	if err != nil {
		return nil, err
	}
	mintToAddr := ethereum.BytesToAddress(mintToAddrBytes)

	return &MintEvent{
		TxHash:          txHash.String(),
		Amount:          *new(big.Int).Mul(amount, big.NewInt(1e10)),
		MintToAddr:      mintToAddr.String(),
		BtcReceiverName: value.BtcReceiverName,
		BtcReceiverAddr: value.BtcReceiverAddr,
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
		if attr.Key == "amount" {
			err = json.Unmarshal([]byte(attr.Value), &amount)
			if err != nil {
				return nil, err
			}
		}

		value := strings.Trim(attr.Value, "\"")

		if attr.Key == "btc_target_address" {
			btcTargetAddress = value
		}

		if attr.Key == "signer" {
			address, err := sdk.GetFromBech32(value, Bech32PrefixAccAddr)
			if err != nil {
				return nil, err
			}
			signer = ethereum.BytesToAddress(address).String()
		}
	}

	return &BurnEvent{
		Amount:           amount,
		BtcTargetAddress: btcTargetAddress,
		Signer:           signer,
	}, nil
}
