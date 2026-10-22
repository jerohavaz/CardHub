package config

import (
	"testing"
	"time"
)

func TestEnv(t *testing.T) {
	t.Setenv("CARDHUB_CONFIG_TEST", "from-env")
	if got := Env("CARDHUB_CONFIG_TEST", "fallback"); got != "from-env" {
		t.Fatalf("Env()=%q", got)
	}
	if got := Env("CARDHUB_NONEXISTENT_TEST", "fallback"); got != "fallback" {
		t.Fatalf("fallback=%q", got)
	}
}
func TestRequired(t *testing.T) {
	t.Setenv("CARDHUB_REQUIRED_TEST", "")
	if _, err := Required("CARDHUB_REQUIRED_TEST"); err == nil {
		t.Fatal("expected missing variable error")
	}
}
func TestDurationEnv(t *testing.T) {
	t.Setenv("CARDHUB_DURATION_TEST", "2s")
	if got, err := DurationEnv("CARDHUB_DURATION_TEST", time.Second); err != nil || got != 2*time.Second {
		t.Fatalf("duration=%v err=%v", got, err)
	}
	t.Setenv("CARDHUB_DURATION_TEST", "-2s")
	if _, err := DurationEnv("CARDHUB_DURATION_TEST", time.Second); err == nil {
		t.Fatal("expected invalid duration error")
	}
}
