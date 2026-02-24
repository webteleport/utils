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
		auth     string
		wantUser string
		wantPass string
		wantOk   bool
	}{
		{encode("alice", "secret"), "alice", "secret", true},
		{encode("user", ""), "user", "", true},
		{encode("user", "p:ass"), "user", "p:ass", true},
		{"Basic " + base64.StdEncoding.EncodeToString([]byte("noconn")), "", "", false},
		{"Bearer token", "", "", false},
		{"", "", "", false},
		{"Basic !!invalid!!", "", "", false},
	}
	for _, tt := range tests {
		gotUser, gotPass, gotOk := ParseBasicAuth(tt.auth)
		if gotOk != tt.wantOk || gotUser != tt.wantUser || gotPass != tt.wantPass {
			t.Errorf("ParseBasicAuth(%q) = (%q, %q, %v), want (%q, %q, %v)",
				tt.auth, gotUser, gotPass, gotOk, tt.wantUser, tt.wantPass, tt.wantOk)
		}
	}
}
