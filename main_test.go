package main

import (
	"os"
	"testing"
)

func TestRun_MissingConfig(t *testing.T) {
	os.Unsetenv("PUSHOVER_API_TOKEN")
	os.Unsetenv("PUSHOVER_USER_KEY")

	err := run()
	if err == nil {
		t.Fatal("expected error for missing config")
	}
}
