package sms

import (
	"context"
	"log/slog"
)

// ConsoleSender "delivers" a code by logging it. It exists purely so OTP login is
// testable without a real SMS account — it is NOT safe for production, since an
// OTP must never appear in logs. cmd/api/main.go only constructs this when
// OTP_DEV_LOG_CODES=true, which the production .env example leaves unset.
type ConsoleSender struct {
	logger *slog.Logger
}

func NewConsoleSender(logger *slog.Logger) *ConsoleSender {
	return &ConsoleSender{logger: logger}
}

func (s *ConsoleSender) Send(_ context.Context, phone, code string) error {
	s.logger.Warn("DEV ONLY — OTP code printed to the console, never do this in production",
		"phone", phone, "code", code)
	return nil
}
