package smsru

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"palomnik/internal/domain"
	"palomnik/internal/ports"
)

const (
	defaultBaseURL   = "https://sms.ru"
	defaultExpiresIn = 5 * 60
)

type Verifier struct {
	apiID      string
	baseURL    string
	httpClient *http.Client
}

func New(apiID string, timeout time.Duration) *Verifier {
	if timeout <= 0 {
		timeout = 15 * time.Second
	}
	return &Verifier{
		apiID:   strings.TrimSpace(apiID),
		baseURL: defaultBaseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (v *Verifier) Available() bool {
	return v != nil && strings.TrimSpace(v.apiID) != ""
}

type addResponse struct {
	Status          string `json:"status"`
	StatusCode      int    `json:"status_code"`
	StatusText      string `json:"status_text"`
	CheckID         string `json:"check_id"`
	CallPhone       string `json:"call_phone"`
	CallPhonePretty string `json:"call_phone_pretty"`
}

type statusResponse struct {
	Status          string `json:"status"`
	StatusCode      int    `json:"status_code"`
	StatusText      string `json:"status_text"`
	CheckStatus     string `json:"check_status"`
	CheckStatusText string `json:"check_status_text"`
}

func (v *Verifier) Start(ctx context.Context, phone string) (ports.PhoneChallenge, error) {
	if !v.Available() {
		return ports.PhoneChallenge{}, ports.ErrPhoneVerifierNotConfigured
	}

	apiPhone, err := toAPIPhone(phone)
	if err != nil {
		return ports.PhoneChallenge{}, err
	}

	form := url.Values{}
	form.Set("api_id", v.apiID)
	form.Set("phone", apiPhone)
	form.Set("json", "1")

	var resp addResponse
	if err := v.postForm(ctx, "/callcheck/add", form, &resp); err != nil {
		return ports.PhoneChallenge{}, err
	}
	if !strings.EqualFold(resp.Status, "OK") || resp.StatusCode != 100 {
		return ports.PhoneChallenge{}, fmt.Errorf("%w: %s (%d)", ports.ErrPhoneChallengeFailed, resp.StatusText, resp.StatusCode)
	}
	if strings.TrimSpace(resp.CheckID) == "" || strings.TrimSpace(resp.CallPhone) == "" {
		return ports.PhoneChallenge{}, fmt.Errorf("%w: empty check payload", ports.ErrPhoneChallengeFailed)
	}

	pretty := strings.TrimSpace(resp.CallPhonePretty)
	if pretty == "" {
		pretty = resp.CallPhone
	}

	return ports.PhoneChallenge{
		CheckID:         resp.CheckID,
		CallPhone:       resp.CallPhone,
		CallPhonePretty: pretty,
		ExpiresIn:       defaultExpiresIn,
	}, nil
}

func (v *Verifier) Status(ctx context.Context, checkID string) (ports.PhoneCheckStatus, error) {
	if !v.Available() {
		return "", ports.ErrPhoneVerifierNotConfigured
	}
	checkID = strings.TrimSpace(checkID)
	if checkID == "" {
		return "", ports.ErrPhoneChallengeFailed
	}

	form := url.Values{}
	form.Set("api_id", v.apiID)
	form.Set("check_id", checkID)
	form.Set("json", "1")

	var resp statusResponse
	if err := v.postForm(ctx, "/callcheck/status", form, &resp); err != nil {
		return "", err
	}
	if !strings.EqualFold(resp.Status, "OK") || resp.StatusCode != 100 {
		return "", fmt.Errorf("%w: %s (%d)", ports.ErrPhoneChallengeFailed, resp.StatusText, resp.StatusCode)
	}

	switch strings.TrimSpace(resp.CheckStatus) {
	case "400":
		return ports.PhoneCheckPending, nil
	case "401":
		return ports.PhoneCheckConfirmed, nil
	case "402":
		return ports.PhoneCheckExpired, nil
	default:
		return "", fmt.Errorf("%w: unknown check_status %q", ports.ErrPhoneChallengeFailed, resp.CheckStatus)
	}
}

func (v *Verifier) postForm(ctx context.Context, path string, form url.Values, dest any) error {
	endpoint := strings.TrimRight(v.baseURL, "/") + path
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := v.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return err
	}
	if resp.StatusCode >= 300 {
		return fmt.Errorf("%w: http %d: %s", ports.ErrPhoneChallengeFailed, resp.StatusCode, strings.TrimSpace(string(body)))
	}
	if err := json.Unmarshal(body, dest); err != nil {
		return fmt.Errorf("%w: decode: %v", ports.ErrPhoneChallengeFailed, err)
	}
	return nil
}

func toAPIPhone(phone string) (string, error) {
	normalized := domain.NormalizePhone(phone)
	if normalized == "" {
		return "", domain.ErrInvalidPhone
	}
	digits := strings.Builder{}
	for _, r := range normalized {
		if r >= '0' && r <= '9' {
			digits.WriteRune(r)
		}
	}
	value := digits.String()
	if len(value) < 10 {
		return "", domain.ErrInvalidPhone
	}
	return value, nil
}
