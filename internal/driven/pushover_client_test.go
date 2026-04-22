package driven

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

const (
	errNewClient = "NewClient() error = %v"
	errSend      = "Send() error = %v"
	errSendNil   = "Send() error = nil, want non-nil"
	testURL      = "https://example.com"

	testAPIToken = "token-1"
	testUserKey  = "user-1"
)

func testConfig(apiURL string) Config {
	return Config{
		APIToken: testAPIToken,
		UserKey:  testUserKey,
		APIURL:   apiURL,
	}
}

func newTestClient(t *testing.T, ts *httptest.Server) *PushoverClient {
	t.Helper()

	client, err := NewPushoverClient(testConfig(ts.URL), ts.Client())
	if err != nil {
		t.Fatalf(errNewClient, err)
	}

	return client
}

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

func TestNewClient_DefaultAPIURL(t *testing.T) {
	client, err := NewPushoverClient(testConfig(""), &http.Client{})
	if err != nil {
		t.Fatalf(errNewClient, err)
	}

	if client == nil {
		t.Fatal("client is nil")
	}
}

func setupTestServer(t *testing.T, received *url.Values) *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertRequestMethodAndContentType(t, r)

		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body error: %v", err)
		}

		*received, err = url.ParseQuery(string(body))
		if err != nil {
			t.Fatalf("parse query error: %v", err)
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":1}`))
	}))

	return ts
}

func assertRequestMethodAndContentType(t *testing.T, r *http.Request) {
	t.Helper()

	if r.Method != http.MethodPost {
		t.Fatalf("method = %s, want POST", r.Method)
	}

	if got := r.Header.Get("Content-Type"); got != "application/x-www-form-urlencoded" {
		t.Fatalf("content-type = %q, want %q", got, "application/x-www-form-urlencoded")
	}
}

func assertFormValues(t *testing.T, received url.Values, expected map[string]string) {
	t.Helper()

	for key, want := range expected {
		if got := received.Get(key); got != want {
			t.Fatalf("%s = %q, want %q", key, got, want)
		}
	}
}

func assertFormValueEmpty(t *testing.T, received url.Values, key string) {
	t.Helper()

	if got := received.Get(key); got != "" {
		t.Fatalf("%s = %q, want empty", key, got)
	}
}

func TestSend_SuccessAndFormPayload(t *testing.T) {
	var received url.Values

	ts := setupTestServer(t, &received)
	defer ts.Close()

	client := newTestClient(t, ts)

	priority := 2
	title := "Build"
	sound := "pushover"
	testURLValue := testURL
	urlTitle := "Link"
	device := "iphone"

	n := domain.Notification{
		Message:  "deployed",
		Title:    &title,
		Priority: &priority,
		Sound:    &sound,
		URL:      &testURLValue,
		URLTitle: &urlTitle,
		Device:   &device,
	}

	if err := client.Send(context.Background(), n); err != nil {
		t.Fatalf(errSend, err)
	}

	assertFormValues(t, received, map[string]string{
		"token":     testAPIToken,
		"user":      testUserKey,
		"message":   "deployed",
		"title":     "Build",
		"priority":  "2",
		"retry":     "60",
		"expire":    "3600",
		"sound":     "pushover",
		"url":       testURL,
		"url_title": "Link",
		"device":    "iphone",
	})
}

func TestSend_PriorityNotEmergency_NoRetryExpire(t *testing.T) {
	var received url.Values

	ts := setupTestServer(t, &received)
	defer ts.Close()

	client := newTestClient(t, ts)

	priority := 1
	n := domain.Notification{
		Message:  "test",
		Priority: &priority,
	}

	if err := client.Send(context.Background(), n); err != nil {
		t.Fatalf(errSend, err)
	}

	assertFormValueEmpty(t, received, "retry")
	assertFormValueEmpty(t, received, "expire")
}

func TestSend_NoPriority(t *testing.T) {
	var received url.Values

	ts := setupTestServer(t, &received)
	defer ts.Close()

	client := newTestClient(t, ts)

	n := domain.Notification{
		Message: "test without priority",
	}

	if err := client.Send(context.Background(), n); err != nil {
		t.Fatalf(errSend, err)
	}

	assertFormValueEmpty(t, received, "priority")
}

func TestSend_EmptyOptionalFields(t *testing.T) {
	var received url.Values

	ts := setupTestServer(t, &received)
	defer ts.Close()

	client := newTestClient(t, ts)

	emptyTitle := ""
	emptySound := "   "
	n := domain.Notification{
		Message: "test",
		Title:   &emptyTitle,
		Sound:   &emptySound,
	}

	if err := client.Send(context.Background(), n); err != nil {
		t.Fatalf(errSend, err)
	}

	assertFormValueEmpty(t, received, "title")
	assertFormValueEmpty(t, received, "sound")
}

func TestSend_RequestError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":1}`))
	}))
	defer ts.Close()

	client := newTestClient(t, ts)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := client.Send(ctx, domain.Notification{Message: "test"})
	if err == nil {
		t.Fatal(errSendNil)
	}
}

func TestSend_InvalidURL(t *testing.T) {
	client, err := NewPushoverClient(Config{
		APIToken: testAPIToken,
		UserKey:  testUserKey,
		APIURL:   "http://\x00invalid",
	}, &http.Client{})
	if err != nil {
		t.Fatalf(errNewClient, err)
	}

	err = client.Send(context.Background(), domain.Notification{Message: "test"})
	if err == nil {
		t.Fatal(errSendNil)
	}

	if !strings.Contains(err.Error(), "create request") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSend_Non2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"status":0,"errors":["invalid token"]}`))
	}))
	defer ts.Close()

	client := newTestClient(t, ts)

	err := client.Send(context.Background(), domain.Notification{Message: "hello"})
	if err == nil {
		t.Fatal(errSendNil)
	}

	if !strings.Contains(err.Error(), "pushover returned 400 Bad Request") {
		t.Fatalf("unexpected error: %v", err)
	}
}
