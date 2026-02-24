package utils

import (
	"net/http"
	"testing"
)

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

func TestRealIP(t *testing.T) {
	tests := []struct {
		name       string
		headers    map[string]string
		remoteAddr string
		want       string
	}{
		{
			name:       "X-Forwarded-For single IP",
			headers:    map[string]string{"X-Forwarded-For": "1.2.3.4"},
			remoteAddr: "10.0.0.1:1234",
			want:       "1.2.3.4",
		},
		{
			name:       "X-Forwarded-For multiple IPs returns first",
			headers:    map[string]string{"X-Forwarded-For": "1.2.3.4, 5.6.7.8, 9.10.11.12"},
			remoteAddr: "10.0.0.1:1234",
			want:       "1.2.3.4",
		},
		{
			name:       "X-Real-IP takes precedence over X-Forwarded-For",
			headers:    map[string]string{"X-Real-IP": "2.2.2.2", "X-Forwarded-For": "1.1.1.1"},
			remoteAddr: "10.0.0.1:1234",
			want:       "2.2.2.2",
		},
		{
			name:       "falls back to RemoteAddr when no headers",
			headers:    map[string]string{},
			remoteAddr: "10.0.0.1:1234",
			want:       "10.0.0.1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &http.Request{
				Header:     make(http.Header),
				RemoteAddr: tt.remoteAddr,
			}
			for k, v := range tt.headers {
				r.Header.Set(k, v)
			}
			got := RealIP(r)
			if got != tt.want {
				t.Errorf("RealIP() = %q, want %q", got, tt.want)
			}
		})
	}
}
