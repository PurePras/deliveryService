package sms

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const msg91FlowURL = "https://control.msg91.com/api/v5/flow/"

// MSG91Sender sends a code via MSG91's Flow API using a pre-approved DLT template —
// mandatory for commercial SMS to Indian numbers. TemplateID must reference a template
// registered with exactly one variable, named OTP, e.g.:
// "Your OTP for Shri Ram Sabji Service is ##OTP##. Valid for 5 minutes. Do not share this code."
type MSG91Sender struct {
	AuthKey    string
	TemplateID string
	httpClient *http.Client
}

func NewMSG91Sender(authKey, templateID string) *MSG91Sender {
	return &MSG91Sender{
		AuthKey:    authKey,
		TemplateID: templateID,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *MSG91Sender) Send(ctx context.Context, phone, code string) error {
	body := map[string]any{
		"template_id": s.TemplateID,
		"short_url":   "0",
		"recipients": []map[string]string{
			{
				// MSG91 expects the country code with no leading '+'; phone is
				// always a bare 10-digit Indian number (see requirePhone).
				"mobiles": "91" + phone,
				"OTP":     code,
			},
		},
	}
	payload, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, msg91FlowURL, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("authkey", s.AuthKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("msg91: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("msg91: unexpected status %d: %s", resp.StatusCode, string(respBody))
	}
	return nil
}
