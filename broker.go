package nsl

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type EnrollmentToken struct {
	TokenID         string    `json:"token_id"`
	NodeName        string    `json:"node_name"`
	EnrollmentToken string    `json:"enrollment_token"`
	ExpiresAt       time.Time `json:"expires_at"`
}

func IssueEnrollmentToken(ctx context.Context, brokerURL, adminToken, nodeName string, ttl time.Duration) (EnrollmentToken, error) {
	input := map[string]any{"node_name": nodeName, "ttl_seconds": int(ttl.Seconds())}
	body, err := json.Marshal(input)
	if err != nil {
		return EnrollmentToken{}, err
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimRight(brokerURL, "/")+"/v1/admin/enrollment-tokens", bytes.NewReader(body))
	if err != nil {
		return EnrollmentToken{}, err
	}
	request.Header.Set("Authorization", "Bearer "+adminToken)
	request.Header.Set("Content-Type", "application/json")
	response, err := (&http.Client{Timeout: 30 * time.Second}).Do(request)
	if err != nil {
		return EnrollmentToken{}, fmt.Errorf("broker unavailable: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusCreated {
		message, _ := io.ReadAll(io.LimitReader(response.Body, 4096))
		return EnrollmentToken{}, fmt.Errorf("broker returned %d: %s", response.StatusCode, strings.TrimSpace(string(message)))
	}
	var token EnrollmentToken
	if err := json.NewDecoder(response.Body).Decode(&token); err != nil {
		return EnrollmentToken{}, err
	}
	return token, nil
}
