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
		TxHash       string  `json:"tx_hash"`
		Amount       big.Int `json:"amount"`
		ReceiverAddr string  `json:"receiver_addr"`
		AgentName    string  `json:"agent_name"`
		AgentBtcAddr string  `json:"agent_btc_addr"`
		ChainId      uint32  `json:"chain_id"`
	}

	MintRecordValue struct {
		TxHash       string `json:"tx_hash"`
		Amount       string `json:"amount"`
		ReceiverAddr string `json:"receiver_addr,omitempty"`
		AgentName    string `json:"agent_name,omitempty"`
		AgentBtcAddr string `json:"agent_btc_addr,omitempty"`
		ChainId      uint32 `json:"chain_id,omitempty"`
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

	receiverAddrBytes, err := base64.StdEncoding.DecodeString(value.ReceiverAddr)
	if err != nil {
		return nil, err
	}
	receiverAddr := ethereum.BytesToAddress(receiverAddrBytes)

	return &MintEvent{
		TxHash:       txHash.String(),
		Amount:       *new(big.Int).Mul(amount, big.NewInt(1e10)),
		ReceiverAddr: receiverAddr.String(),
		AgentName:    value.AgentName,
		AgentBtcAddr: value.AgentBtcAddr,
		ChainId:      value.ChainId,
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
