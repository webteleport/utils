package utils

import (
	"encoding/base64"
	"testing"
)

func TestParseBasicAuth(t *testing.T) {
	encode := func(user, pass string) string {
		return "Basic " + base64.StdEncoding.EncodeToString([]byte(user+":"+pass))
	}

	tests := []struct {
		name         string
		auth         string
		wantUser     string
		wantPass     string
		wantOK       bool
	}{
		{"valid credentials", encode("alice", "secret"), "alice", "secret", true},
		{"empty password", encode("bob", ""), "bob", "", true},
		{"colon in password", encode("user", "pa:ss"), "user", "pa:ss", true},
		{"missing Basic prefix", "Bearer token123", "", "", false},
		{"empty string", "", "", "", false},
		{"invalid base64", "Basic !!!invalid!!!", "", "", false},
		{"no colon in decoded", "Basic " + base64.StdEncoding.EncodeToString([]byte("nodivider")), "", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, pass, ok := ParseBasicAuth(tt.auth)
			if ok != tt.wantOK {
				t.Errorf("ParseBasicAuth(%q) ok = %v, want %v", tt.auth, ok, tt.wantOK)
			}
			if user != tt.wantUser {
				t.Errorf("ParseBasicAuth(%q) user = %q, want %q", tt.auth, user, tt.wantUser)
			}
			if pass != tt.wantPass {
				t.Errorf("ParseBasicAuth(%q) pass = %q, want %q", tt.auth, pass, tt.wantPass)
			}
		})
	}
}
