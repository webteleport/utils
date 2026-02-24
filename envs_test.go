package utils

import (
	"os"
	"testing"
)

func TestEnvOr(t *testing.T) {
	tests := []struct {
		key      string
		setValue string
		set      bool
		fallback string
		want     string
	}{
		{"ENVTEST_UNSET", "", false, "default", "default"},
		{"ENVTEST_SET", "value", true, "default", "value"},
		{"ENVTEST_EMPTY", "", true, "default", ""},
	}
	for _, tt := range tests {
		os.Unsetenv(tt.key)
		if tt.set {
			os.Setenv(tt.key, tt.setValue)
			t.Cleanup(func() { os.Unsetenv(tt.key) })
		}
		got := EnvOr(tt.key, tt.fallback)
		if got != tt.want {
			t.Errorf("EnvOr(%q, %q) = %q, want %q", tt.key, tt.fallback, got, tt.want)
		}
	}
}

func TestLookupEnv(t *testing.T) {
	// Test when the environment variable is not set
	os.Clearenv()
	result := LookupEnv("NON_EXISTING_VAR")
	if result != nil {
		t.Errorf("Expected nil, got %s", *result)
	}

	// Test when the environment variable is set
	expectedValue := "test_value"
	os.Setenv("EXISTING_VAR", expectedValue)
	result = LookupEnv("EXISTING_VAR")
	if result == nil || *result != expectedValue {
		t.Errorf("Expected %s, got %s", expectedValue, *result)
	}
}
