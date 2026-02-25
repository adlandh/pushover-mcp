package adapters

import (
	"context"
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
		return nil, fmt.Errorf("missing APIToken")
	}

	if cfg.UserKey == "" {
		return nil, fmt.Errorf("missing UserKey")
	}

	if httpClient == nil {
		return nil, fmt.Errorf("http client is required")
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
	setPriority(form, notification.Priority)
	setOptionalString(form, "sound", notification.Sound)
	setOptionalString(form, "url", notification.URL)
	setOptionalString(form, "url_title", notification.URLTitle)
	setOptionalString(form, "device", notification.Device)

	return form
}

func setOptionalString(form url.Values, key string, value *string) {
	if value == nil {
		return
	}

	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return
	}

	form.Set(key, *value)
}

func setPriority(form url.Values, priority *int) {
	if priority == nil {
		return
	}

	form.Set("priority", strconv.Itoa(*priority))

	if *priority == 2 {
		form.Set("retry", "60")
		form.Set("expire", "3600")
	}
}

func validateResponse(resp *http.Response) error {
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("pushover returned %s: %s", resp.Status, strings.TrimSpace(string(body)))
	}

	return nil
}
