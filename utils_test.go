package utils

import "testing"

func TestExtractPort(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"example.com:8080", ":8080"},
		{"localhost:443", ":443"},
		{"example.com", ""},
		{"[::1]:9090", ":9090"},
		{"[::1]", ""},
	}
	for _, tt := range tests {
		got := ExtractPort(tt.input)
		if got != tt.want {
			t.Errorf("ExtractPort(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
