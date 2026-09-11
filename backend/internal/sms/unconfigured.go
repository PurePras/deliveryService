package sms

import (
	"context"
	"errors"
)

// Unconfigured is the default when no real provider and no dev fallback are set
// up — it fails loudly instead of silently pretending an OTP was sent.
type Unconfigured struct{}

func (Unconfigured) Send(context.Context, string, string) error {
	return errors.New("no SMS provider configured: set OTP_DEV_LOG_CODES=true for local dev, or wire a real sms.Sender into cmd/api/main.go")
}
