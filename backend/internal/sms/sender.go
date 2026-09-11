// Package sms delivers OTP codes to phone numbers. No real provider is wired in
// yet — swap one in (Twilio, MSG91, AWS SNS, ...) by implementing Sender and
// constructing it in cmd/api/main.go in place of the fallbacks in this package.
package sms

import "context"

type Sender interface {
	// Send delivers code to phone. The caller (AuthService) has already generated
	// and stored the code — this is purely the delivery step.
	Send(ctx context.Context, phone, code string) error
}
