package config

import (
	"strings"
	"testing"
	"time"
)

const (
	testAPIToken = "token-123"
	testUserKey  = "user-456"
	testAPIURL   = "https://example.com/messages"
	errParseEnv  = "parse env"
)

func setPushoverEnv(t *testing.T, token, userKey, apiURL, timeout string) {
	t.Helper()
	t.Setenv("PUSHOVER_API_TOKEN", token)
	t.Setenv("PUSHOVER_USER_KEY", userKey)
	t.Setenv("PUSHOVER_API_URL", apiURL)
	t.Setenv("PUSHOVER_TIMEOUT", timeout)
}

func assertParseEnvError(t *testing.T, err error, wantField string) {
	t.Helper()

	if err == nil {
		t.Fatal("FromEnv() error = nil, want non-nil")
	}

	errText := err.Error()
	if !strings.Contains(errText, errParseEnv) || !strings.Contains(errText, wantField) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestFromEnv_Success(t *testing.T) {
	setPushoverEnv(t, testAPIToken, testUserKey, testAPIURL, "30s")

	cfg, err := FromEnv()
	if err != nil {
		t.Fatalf("FromEnv() error = %v", err)
	}

	if cfg.Pushover.APIToken != testAPIToken {
		t.Fatalf("APIToken = %q, want %q", cfg.Pushover.APIToken, testAPIToken)
	}

	if cfg.Pushover.UserKey != testUserKey {
		t.Fatalf("UserKey = %q, want %q", cfg.Pushover.UserKey, testUserKey)
	}

	if cfg.Pushover.APIURL != testAPIURL {
		t.Fatalf("APIURL = %q, want %q", cfg.Pushover.APIURL, testAPIURL)
	}

	if cfg.Timeout != 30*time.Second {
		t.Fatalf("Timeout = %v, want %v", cfg.Timeout, 30*time.Second)
	}
}

func TestFromEnv_DefaultTimeout(t *testing.T) {
	setPushoverEnv(t, testAPIToken, testUserKey, "", "")

	cfg, err := FromEnv()
	if err != nil {
		t.Fatalf("FromEnv() error = %v", err)
	}

	if cfg.Timeout != 15*time.Second {
		t.Fatalf("Timeout = %v, want %v", cfg.Timeout, 15*time.Second)
	}
}

func TestFromEnv_MissingAPIToken(t *testing.T) {
	setPushoverEnv(t, "", testUserKey, "", "")

	_, err := FromEnv()
	assertParseEnvError(t, err, "PUSHOVER_API_TOKEN")
}

func TestFromEnv_MissingUserKey(t *testing.T) {
	setPushoverEnv(t, testAPIToken, "", "", "")

	_, err := FromEnv()
	assertParseEnvError(t, err, "PUSHOVER_USER_KEY")
}
