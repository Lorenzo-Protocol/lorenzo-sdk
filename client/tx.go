package client

import (
	"context"
	"fmt"
	"sync"

	"cosmossdk.io/errors"
	agenttypes "github.com/Lorenzo-Protocol/lorenzo/v2/x/agent/types"
	btclctypes "github.com/Lorenzo-Protocol/lorenzo/v2/x/btclightclient/types"
	btcstakingtypes "github.com/Lorenzo-Protocol/lorenzo/v2/x/btcstaking/types"
	plantypes "github.com/Lorenzo-Protocol/lorenzo/v2/x/plan/types"
	"github.com/avast/retry-go/v4"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/relayer/v2/relayer/chains/cosmos"
	pv "github.com/cosmos/relayer/v2/relayer/provider"
	"go.uber.org/zap"
)

// ToProviderMsgs converts a list of sdk.Msg to a list of provider.RelayerMessage
func ToProviderMsgs(msgs []sdk.Msg) []pv.RelayerMessage {
	relayerMsgs := []pv.RelayerMessage{}
	for _, m := range msgs {
		relayerMsgs = append(relayerMsgs, cosmos.NewCosmosMessage(m, func(signer string) {}))
	}
	return relayerMsgs
}

// SendMsgToMempool sends a message to the mempool.
// It does not wait for the messages to be included.
func (c *Client) SendMsgToMempool(ctx context.Context, msg sdk.Msg) error {
	return c.SendMsgsToMempool(ctx, []sdk.Msg{msg})
}

// SendMsgsToMempool sends a set of messages to the mempool.
// It does not wait for the messages to be included.
func (c *Client) SendMsgsToMempool(ctx context.Context, msgs []sdk.Msg) error {
	if len(msgs) == 0 {
		return fmt.Errorf("empty message set provided")
	}

	relayerMsgs := ToProviderMsgs(msgs)
	if err := retry.Do(func() error {
		var sendMsgErr error
		krErr := c.accessKeyWithLock(func() {
			sendMsgErr = c.provider.SendMessagesToMempool(ctx, relayerMsgs, "", ctx, []func(*pv.RelayerTxResponse, error){})
		})
		if krErr != nil {
			c.logger.Error("unrecoverable err when submitting the tx, skip retrying", zap.Error(krErr))
			return retry.Unrecoverable(krErr)
		}
		return sendMsgErr
	}, retry.Context(ctx), rtyAtt, rtyDel, rtyErr, retry.OnRetry(func(n uint, err error) {
		c.logger.Debug("retrying", zap.Uint("attemp", n+1), zap.Uint("max_attempts", rtyAttNum), zap.Error(err))
	})); err != nil {
		return err
	}

	return nil
}

// ReliablySendMsg reliable sends a message to the chain.
// It utilizes a file lock as well as a keyring lock to ensure atomic access.
// TODO: needs tests
func (c *Client) ReliablySendMsg(ctx context.Context, msg sdk.Msg, expectedErrors []*errors.Error, unrecoverableErrors []*errors.Error) (*pv.RelayerTxResponse, error) {
	return c.ReliablySendMsgs(ctx, []sdk.Msg{msg}, expectedErrors, unrecoverableErrors)
}

// ReliablySendMsgs reliably sends a list of messages to the chain.
// It utilizes a file lock as well as a keyring lock to ensure atomic access.
// TODO: needs tests
func (c *Client) ReliablySendMsgs(ctx context.Context, msgs []sdk.Msg, expectedErrors []*errors.Error, unrecoverableErrors []*errors.Error) (*pv.RelayerTxResponse, error) {
	var (
		rlyResp     *pv.RelayerTxResponse
		callbackErr error
		wg          sync.WaitGroup
	)

	callback := func(rtr *pv.RelayerTxResponse, err error) {
		rlyResp = rtr
		callbackErr = err
		wg.Done()
	}

	wg.Add(1)

	// convert message type
	relayerMsgs := ToProviderMsgs(msgs)

	// TODO: consider using Lorenzo's retry package
	if err := retry.Do(func() error {
		var sendMsgErr error
		krErr := c.accessKeyWithLock(func() {
			sendMsgErr = c.provider.SendMessagesToMempool(ctx, relayerMsgs, "", ctx, []func(*pv.RelayerTxResponse, error){callback})
		})
		if krErr != nil {
			c.logger.Error("unrecoverable err when submitting the tx, skip retrying", zap.Error(krErr))
			return retry.Unrecoverable(krErr)
		}
		if sendMsgErr != nil {
			if errorContained(sendMsgErr, unrecoverableErrors) {
				c.logger.Error("unrecoverable err when submitting the tx, skip retrying", zap.Error(sendMsgErr))
				return retry.Unrecoverable(sendMsgErr)
			}
			if errorContained(sendMsgErr, expectedErrors) {
				// this is necessary because if err is returned
				// the callback function will not be executed so
				// that the inside wg.Done will not be executed
				wg.Done()
				c.logger.Error("expected err when submitting the tx, skip retrying", zap.Error(sendMsgErr))
				return nil
			}
			return sendMsgErr
		}
		return nil
	}, retry.Context(ctx), rtyAtt, rtyDel, rtyErr, retry.OnRetry(func(n uint, err error) {
		c.logger.Debug("retrying", zap.Uint("attemp", n+1), zap.Uint("max_attempts", rtyAttNum), zap.Error(err))
	})); err != nil {
		return nil, err
	}

	wg.Wait()

	if callbackErr != nil {
		if errorContained(callbackErr, expectedErrors) {
			return nil, nil
		}
		return nil, callbackErr
	}

	if rlyResp == nil {
		// this case could happen if the error within the retry is an expected error
		return nil, nil
	}

	if rlyResp.Code != 0 {
		return rlyResp, fmt.Errorf("transaction failed with code: %d", rlyResp.Code)
	}

	return rlyResp, nil
}

func (c *Client) InsertHeaders(ctx context.Context, msg *btclctypes.MsgInsertHeaders) (*pv.RelayerTxResponse, error) {
	return c.ReliablySendMsg(ctx, msg, []*errors.Error{}, []*errors.Error{})
}

func (c *Client) CreateBTCStakingWithBTCProof(ctx context.Context, msg *btcstakingtypes.MsgCreateBTCStaking) (*pv.RelayerTxResponse, error) {
	return c.ReliablySendMsg(ctx, msg, []*errors.Error{}, []*errors.Error{})
}

func (c *Client) AddAgent(ctx context.Context, msg *agenttypes.MsgAddAgent) (*pv.RelayerTxResponse, error) {
	return c.ReliablySendMsg(ctx, msg, []*errors.Error{}, []*errors.Error{})
}

func (c *Client) EditAgent(ctx context.Context, msg *agenttypes.MsgEditAgent) (*pv.RelayerTxResponse, error) {
	return c.ReliablySendMsg(ctx, msg, []*errors.Error{}, []*errors.Error{})
}

func (c *Client) RemoveAgent(ctx context.Context, msg *agenttypes.MsgRemoveAgent) (*pv.RelayerTxResponse, error) {
	return c.ReliablySendMsg(ctx, msg, []*errors.Error{}, []*errors.Error{})
}

/**************************
*******	Plan Module ********
************************/

func (c *Client) UpgradePlan(ctx context.Context, msg *plantypes.MsgUpgradePlan) (*pv.RelayerTxResponse, error) {
	return c.ReliablySendMsg(ctx, msg, []*errors.Error{}, []*errors.Error{})
}

func (c *Client) CreatePlan(ctx context.Context, msg *plantypes.MsgCreatePlan) (*pv.RelayerTxResponse, error) {
	return c.ReliablySendMsg(ctx, msg, []*errors.Error{}, []*errors.Error{})
}

func (c *Client) SetPlanMerkleRoot(ctx context.Context, msg *plantypes.MsgSetMerkleRoot) (*pv.RelayerTxResponse, error) {
	return c.ReliablySendMsg(ctx, msg, []*errors.Error{}, []*errors.Error{})
}

func (c *Client) Claims(ctx context.Context, msg *plantypes.MsgClaims) (*pv.RelayerTxResponse, error) {
	return c.ReliablySendMsg(ctx, msg, []*errors.Error{}, []*errors.Error{})
}

func (c *Client) CreateYAT(ctx context.Context, msg *plantypes.MsgCreateYAT) (*pv.RelayerTxResponse, error) {
	return c.ReliablySendMsg(ctx, msg, []*errors.Error{}, []*errors.Error{})
}

func (c *Client) UpdatePlanStatus(ctx context.Context, msg *plantypes.MsgUpdatePlanStatus) (*pv.RelayerTxResponse, error) {
	return c.ReliablySendMsg(ctx, msg, []*errors.Error{}, []*errors.Error{})
}

func (c *Client) SetMinter(ctx context.Context, msg *plantypes.MsgSetMinter) (*pv.RelayerTxResponse, error) {
	return c.ReliablySendMsg(ctx, msg, []*errors.Error{}, []*errors.Error{})
}

func (c *Client) RemoveMinter(ctx context.Context, msg *plantypes.MsgRemoveMinter) (*pv.RelayerTxResponse, error) {
	return c.ReliablySendMsg(ctx, msg, []*errors.Error{}, []*errors.Error{})
}
