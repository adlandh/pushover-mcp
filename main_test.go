package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/adlandh/pushover-mcp/internal/config"
	"github.com/adlandh/pushover-mcp/internal/driven"
)

func TestRun_MissingConfig(t *testing.T) {
	t.Setenv("PUSHOVER_API_TOKEN", "")
	t.Setenv("PUSHOVER_USER_KEY", "")

	err := run()
	if err == nil {
		t.Fatal("expected error for missing config")
	}
}

func TestBuildServer_SendTool_HappyPath(t *testing.T) {
	var (
		gotToken   string
		gotUser    string
		gotMessage string
	)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		form, _ := url.ParseQuery(string(body))
		gotToken = form.Get("token")
		gotUser = form.Get("user")
		gotMessage = form.Get("message")

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":1}`))
	}))
	defer ts.Close()

	env := config.EnvConfig{
		Pushover: driven.Config{
			APIToken: "tok",
			UserKey:  "usr",
			APIURL:   ts.URL,
		},
		Timeout: 5 * time.Second,
	}

	s, err := buildServer(env)
	if err != nil {
		t.Fatalf("buildServer() error = %v", err)
	}

	tool := s.GetTool("send")
	if tool == nil {
		t.Fatal("send tool not registered")
	}

	result, err := tool.Handler(context.Background(), mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      "send",
			Arguments: map[string]any{"message": "hello"},
		},
	})
	if err != nil {
		t.Fatalf("handler error = %v", err)
	}
	if result.IsError {
		t.Fatalf("result is error: %v", mcp.GetTextFromContent(result.Content[0]))
	}

	if gotToken != "tok" || gotUser != "usr" || gotMessage != "hello" {
		t.Fatalf("unexpected request: token=%q user=%q message=%q", gotToken, gotUser, gotMessage)
	}
}
