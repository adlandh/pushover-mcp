package config

import (
	"strings"
	"testing"
)

func TestFromEnv_Success(t *testing.T) {
	t.Setenv("PUSHOVER_API_TOKEN", "token-123")
	t.Setenv("PUSHOVER_USER_KEY", "user-456")
	t.Setenv("PUSHOVER_API_URL", "https://example.com/messages")

	cfg, err := FromEnv()
	if err != nil {
		t.Fatalf("FromEnv() error = %v", err)
	}

	if cfg.Pushover.APIToken != "token-123" {
		t.Fatalf("APIToken = %q, want %q", cfg.Pushover.APIToken, "token-123")
	}
	if cfg.Pushover.UserKey != "user-456" {
		t.Fatalf("UserKey = %q, want %q", cfg.Pushover.UserKey, "user-456")
	}
	if cfg.Pushover.APIURL != "https://example.com/messages" {
		t.Fatalf("APIURL = %q, want %q", cfg.Pushover.APIURL, "https://example.com/messages")
	}
}

func TestFromEnv_MissingAPIToken(t *testing.T) {
	t.Setenv("PUSHOVER_API_TOKEN", "")
	t.Setenv("PUSHOVER_USER_KEY", "user-456")
	t.Setenv("PUSHOVER_API_URL", "")

	_, err := FromEnv()
	if err == nil {
		t.Fatal("FromEnv() error = nil, want non-nil")
	}

	errText := err.Error()
	if !strings.Contains(errText, "parse env") || !strings.Contains(errText, "PUSHOVER_API_TOKEN") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestFromEnv_MissingUserKey(t *testing.T) {
	t.Setenv("PUSHOVER_API_TOKEN", "token-123")
	t.Setenv("PUSHOVER_USER_KEY", "")
	t.Setenv("PUSHOVER_API_URL", "")

	_, err := FromEnv()
	if err == nil {
		t.Fatal("FromEnv() error = nil, want non-nil")
	}

	errText := err.Error()
	if !strings.Contains(errText, "parse env") || !strings.Contains(errText, "PUSHOVER_USER_KEY") {
		t.Fatalf("unexpected error: %v", err)
	}
}
