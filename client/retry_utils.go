package client

import (
	"strings"
	"time"

	"cosmossdk.io/errors"
	"github.com/avast/retry-go/v4"
)

// Variables used for retries
var (
	rtyAttNum = uint(5)
	rtyAtt    = retry.Attempts(rtyAttNum)
	rtyDel    = retry.Delay(time.Millisecond * 400)
	rtyErr    = retry.LastErrorOnly(true)
)

func (c *Client) SetRetryAttempts(retryNumber uint) {
	rtyAttNum = retryNumber
	rtyAtt = retry.Attempts(rtyAttNum)
}

func (c *Client) SetRetryDelay(milliseconds int) {
	rtyDel = retry.Delay(time.Millisecond * time.Duration(milliseconds))
}

func (c *Client) SetIsReturnLatestErrorOnly(isReturnLatestErrorOnly bool) {
	rtyErr = retry.LastErrorOnly(isReturnLatestErrorOnly)
}

func errorContained(err error, errList []*errors.Error) bool {
	for _, e := range errList {
		if strings.Contains(err.Error(), e.Error()) {
			return true
		}
	}

	return false
}
