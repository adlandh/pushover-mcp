package driven

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/adlandh/pushover-mcp/internal/domain"
)

const defaultAPIBaseURL = "https://api.pushover.net/1/messages.json"

type Config struct {
	APIToken string
	UserKey  string
	APIURL   string
}

type PushoverClient struct {
	httpClient *http.Client
	apiToken   string
	userKey    string
	apiURL     string
}

func NewPushoverClient(cfg Config, httpClient *http.Client) (*PushoverClient, error) {
	if cfg.APIToken == "" {
		return nil, errors.New("missing APIToken")
	}

	if cfg.UserKey == "" {
		return nil, errors.New("missing UserKey")
	}

	if httpClient == nil {
		return nil, errors.New("http client is required")
	}

	apiURL := cfg.APIURL
	if strings.TrimSpace(apiURL) == "" {
		apiURL = defaultAPIBaseURL
	}

	return &PushoverClient{
		apiToken:   cfg.APIToken,
		userKey:    cfg.UserKey,
		apiURL:     apiURL,
		httpClient: httpClient,
	}, nil
}

func (c *PushoverClient) Send(ctx context.Context, notification domain.Notification) error {
	form := buildFormValues(c.apiToken, c.userKey, notification)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.apiURL, strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	//nolint:gosec // API URL is controlled by explicit runtime configuration.
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request pushover: %w", err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if err := validateResponse(resp); err != nil {
		return err
	}

	return nil
}

func buildFormValues(apiToken, userKey string, notification domain.Notification) url.Values {
	form := url.Values{}
	form.Set("token", apiToken)
	form.Set("user", userKey)
	form.Set("message", notification.Message)
	setOptionalString(form, "title", notification.Title)
	setPriority(form, notification)
	setOptionalString(form, "sound", notification.Sound)
	setOptionalString(form, "url", notification.URL)
	setOptionalString(form, "url_title", notification.URLTitle)
	setOptionalString(form, "device", notification.Device)

	return form
}

const (
	emergencyPriority      = 2
	emergencyRetrySeconds  = 60
	emergencyExpireSeconds = 3600
)

func setOptionalString(form url.Values, key, value string) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return
	}

	form.Set(key, trimmed)
}

func setPriority(form url.Values, notification domain.Notification) {
	if notification.Priority == nil {
		return
	}

	form.Set("priority", strconv.Itoa(*notification.Priority))

	if *notification.Priority == emergencyPriority {
		retry := emergencyRetrySeconds
		if notification.Retry != nil {
			retry = *notification.Retry
		}

		expire := emergencyExpireSeconds
		if notification.Expire != nil {
			expire = *notification.Expire
		}

		form.Set("retry", strconv.Itoa(retry))
		form.Set("expire", strconv.Itoa(expire))
	}
}

func validateResponse(resp *http.Response) error {
	body, err := io.ReadAll(io.LimitReader(resp.Body, 4096))
	if err != nil {
		return fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("pushover returned %s: %s", resp.Status, strings.TrimSpace(string(body)))
	}

	return nil
}
