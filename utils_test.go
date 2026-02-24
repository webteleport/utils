package utils

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
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

func TestUnwrapInnermost(t *testing.T) {
	base := fmt.Errorf("base error")
	wrapped := fmt.Errorf("wrap1: %w", base)
	doubleWrapped := fmt.Errorf("wrap2: %w", wrapped)

	if got := UnwrapInnermost(base); got != base {
		t.Errorf("UnwrapInnermost(base) = %v, want %v", got, base)
	}
	if got := UnwrapInnermost(wrapped); got != base {
		t.Errorf("UnwrapInnermost(wrapped) = %v, want %v", got, base)
	}
	if got := UnwrapInnermost(doubleWrapped); got != base {
		t.Errorf("UnwrapInnermost(doubleWrapped) = %v, want %v", got, base)
	}
	if got := UnwrapInnermost(nil); got != nil {
		t.Errorf("UnwrapInnermost(nil) = %v, want nil", got)
	}
}

func TestGraft(t *testing.T) {
	tests := []struct {
		base string
		alt  string
		want string
	}{
		{"example.com", ":8080", "example.com:8080"},
		{"example.com:443", ":8080", "example.com:8080"},
		{"example.com", "other.com:9090", "other.com:9090"},
		{"example.com:443", "noporthere", "example.com:443"},
	}
	for _, tt := range tests {
		got := Graft(tt.base, tt.alt)
		if got != tt.want {
			t.Errorf("Graft(%q, %q) = %q, want %q", tt.base, tt.alt, got, tt.want)
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
		{"[::1]:9090", "::1"},
		{"[::1]", "[::1]"},
	}
	for _, tt := range tests {
		got := StripPort(tt.input)
		if got != tt.want {
			t.Errorf("StripPort(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestParseDomainCandidates(t *testing.T) {
	tests := []struct {
		input string
		want  []string
	}{
		{"/a/b/c", []string{"a", "b", "c"}},
		{"/single", []string{"single"}},
		{"", nil},
		{"/", nil},
		{"/ spaced /trimmed ", []string{"spaced", "trimmed"}},
	}
	for _, tt := range tests {
		got := ParseDomainCandidates(tt.input)
		if len(got) != len(tt.want) {
			t.Errorf("ParseDomainCandidates(%q) = %v, want %v", tt.input, got, tt.want)
			continue
		}
		for i := range got {
			if got[i] != tt.want[i] {
				t.Errorf("ParseDomainCandidates(%q)[%d] = %q, want %q", tt.input, i, got[i], tt.want[i])
			}
		}
	}
}

func TestExtractURLPort(t *testing.T) {
	tests := []struct {
		input *url.URL
		want  string
	}{
		{nil, ""},
		{&url.URL{Host: "example.com:8080"}, ":8080"},
		{&url.URL{Host: "example.com"}, ""},
	}
	for _, tt := range tests {
		got := ExtractURLPort(tt.input)
		if got != tt.want {
			t.Errorf("ExtractURLPort(%v) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

type customErr struct{ msg string }

func (e *customErr) Error() string { return e.msg }

// Verify errors.As still works with UnwrapInnermost result
func TestUnwrapInnermostPreservesType(t *testing.T) {
	ce := &customErr{msg: "custom"}
	wrapped := fmt.Errorf("outer: %w", ce)

	got := UnwrapInnermost(wrapped)
	var target *customErr
	if !errors.As(got, &target) {
		t.Errorf("UnwrapInnermost did not preserve error type")
	}
}
