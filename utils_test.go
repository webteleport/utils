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
		{"", ""},
	}
	for _, tc := range tests {
		got := ExtractPort(tc.input)
		if got != tc.want {
			t.Errorf("ExtractPort(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestStripPort(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"example.com:8080", "example.com"},
		{"localhost:443", "localhost"},
		{"example.com", "example.com"},
		{"[::1]:80", "::1"},
		{"[::1]", "[::1]"},
	}
	for _, tc := range tests {
		got := StripPort(tc.input)
		if got != tc.want {
			t.Errorf("StripPort(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestParseBasicAuth(t *testing.T) {
	tests := []struct {
		auth     string
		wantUser string
		wantPass string
		wantOk   bool
	}{
		{"Basic dXNlcjpwYXNz", "user", "pass", true},
		{"Basic dXNlcjo=", "user", "", true},
		{"Bearer token", "", "", false},
		{"", "", "", false},
		{"Basic !!!invalid-base64!!!", "", "", false},
	}
	for _, tc := range tests {
		user, pass, ok := ParseBasicAuth(tc.auth)
		if ok != tc.wantOk || user != tc.wantUser || pass != tc.wantPass {
			t.Errorf("ParseBasicAuth(%q) = (%q, %q, %v), want (%q, %q, %v)",
				tc.auth, user, pass, ok, tc.wantUser, tc.wantPass, tc.wantOk)
		}
	}
}
