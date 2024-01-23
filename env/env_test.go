package env

import (
	"os"
	"testing"
)

func TestGetWithDefault(t *testing.T) {
	if value := GetWithDefault("TEST_ENV_VAR", "default"); value != "default" {
		t.Errorf("expected value to be 'default', got '%s'", value)
	}
	if err := os.Setenv("TEST_ENV_VAR", "isset"); err != nil {
		t.Fatalf("failed to set TEST_ENV_VAR")
	}
	if value := GetWithDefault("TEST_ENV_VAR", "default"); value != "isset" {
		t.Errorf("expected value to be 'isset', got '%s'", value)
	}
}
