// Starter test for env() so `go test ./...` does real work in CI. No DB needed.
package main

import "testing"

func TestEnvReturnsFallbackWhenUnset(t *testing.T) {
	t.Setenv("DEVBOARD_TEST_KEY", "")
	if got := env("DEVBOARD_TEST_KEY", "fallback"); got != "fallback" {
		t.Errorf("env() = %q, want %q", got, "fallback")
	}
}

func TestEnvReturnsValueWhenSet(t *testing.T) {
	t.Setenv("DEVBOARD_TEST_KEY", "real")
	if got := env("DEVBOARD_TEST_KEY", "fallback"); got != "real" {
		t.Errorf("env() = %q, want %q", got, "real")
	}
}
