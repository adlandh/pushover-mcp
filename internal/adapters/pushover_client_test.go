package adapters

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/adlandh/pushover-mcp/internal/domain"
)

func TestNewClient_Validation(t *testing.T) {
	_, err := NewPushoverClient(Config{}, &http.Client{})
	if err == nil || !strings.Contains(err.Error(), "missing APIToken") {
		t.Fatalf("expected missing APIToken error, got: %v", err)
	}

	_, err = NewPushoverClient(Config{APIToken: "token"}, &http.Client{})
	if err == nil || !strings.Contains(err.Error(), "missing UserKey") {
		t.Fatalf("expected missing UserKey error, got: %v", err)
	}

	_, err = NewPushoverClient(Config{APIToken: "token", UserKey: "user"}, nil)
	if err == nil || !strings.Contains(err.Error(), "http client is required") {
		t.Fatalf("expected http client is required error, got: %v", err)
	}
}

func TestSend_SuccessAndFormPayload(t *testing.T) {
	var received url.Values

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s, want POST", r.Method)
		}
		if got := r.Header.Get("Content-Type"); got != "application/x-www-form-urlencoded" {
			t.Fatalf("content-type = %q, want %q", got, "application/x-www-form-urlencoded")
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body error: %v", err)
		}

		received, err = url.ParseQuery(string(body))
		if err != nil {
			t.Fatalf("parse query error: %v", err)
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":1}`))
	}))
	defer ts.Close()

	client, err := NewPushoverClient(Config{
		APIToken: "token-1",
		UserKey:  "user-1",
		APIURL:   ts.URL,
	}, ts.Client())
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	priority := 2
	title := "Build"
	n := domain.Notification{
		Message:  "deployed",
		Title:    &title,
		Priority: &priority,
	}

	err = client.Send(context.Background(), n)
	if err != nil {
		t.Fatalf("Send() error = %v", err)
	}

	if received.Get("token") != "token-1" {
		t.Fatalf("token = %q, want %q", received.Get("token"), "token-1")
	}
	if received.Get("user") != "user-1" {
		t.Fatalf("user = %q, want %q", received.Get("user"), "user-1")
	}
	if received.Get("message") != "deployed" {
		t.Fatalf("message = %q, want %q", received.Get("message"), "deployed")
	}
	if received.Get("title") != "Build" {
		t.Fatalf("title = %q, want %q", received.Get("title"), "Build")
	}
	if received.Get("priority") != "2" {
		t.Fatalf("priority = %q, want %q", received.Get("priority"), "2")
	}
	if received.Get("retry") != "60" {
		t.Fatalf("retry = %q, want %q", received.Get("retry"), "60")
	}
	if received.Get("expire") != "3600" {
		t.Fatalf("expire = %q, want %q", received.Get("expire"), "3600")
	}
}

func TestSend_Non2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"status":0,"errors":["invalid token"]}`))
	}))
	defer ts.Close()

	client, err := NewPushoverClient(Config{
		APIToken: "token-1",
		UserKey:  "user-1",
		APIURL:   ts.URL,
	}, ts.Client())
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	err = client.Send(context.Background(), domain.Notification{Message: "hello"})
	if err == nil {
		t.Fatal("Send() error = nil, want non-nil")
	}
	if !strings.Contains(err.Error(), "pushover returned 400 Bad Request") {
		t.Fatalf("unexpected error: %v", err)
	}
}
