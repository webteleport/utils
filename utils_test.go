package utils

import (
	"net/http"
	"net/url"
	"testing"
)

func TestExtractPort(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"host with port", "example.com:8080", ":8080"},
		{"host without port", "example.com", ""},
		{"IPv4 with port", "127.0.0.1:9000", ":9000"},
		{"IPv4 without port", "127.0.0.1", ""},
		{"IPv6 with port", "[::1]:443", ":443"},
		{"empty string", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractPort(tt.input)
			if got != tt.expected {
				t.Errorf("ExtractPort(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestExtractURLPort(t *testing.T) {
	tests := []struct {
		name     string
		rawURL   string
		expected string
	}{
		{"URL with port", "http://example.com:8080/path", ":8080"},
		{"URL without port", "http://example.com/path", ""},
		{"nil URL", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var u *url.URL
			if tt.rawURL != "" {
				var err error
				u, err = url.Parse(tt.rawURL)
				if err != nil {
					t.Fatalf("url.Parse(%q) error: %v", tt.rawURL, err)
				}
			}
			got := ExtractURLPort(u)
			if got != tt.expected {
				t.Errorf("ExtractURLPort(%q) = %q, want %q", tt.rawURL, got, tt.expected)
			}
		})
	}
}

func TestStripPort(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"host with port", "example.com:8080", "example.com"},
		{"host without port", "example.com", "example.com"},
		{"IPv4 with port", "127.0.0.1:9000", "127.0.0.1"},
		{"IPv6 with port", "[::1]:443", "::1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StripPort(tt.input)
			if got != tt.expected {
				t.Errorf("StripPort(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestGraft(t *testing.T) {
	tests := []struct {
		name     string
		base     string
		alt      string
		expected string
	}{
		{"replace port", "example.com:80", ":443", "example.com:443"},
		{"base has no port", "example.com", ":443", "example.com:443"},
		{"alt has host and port", "example.com:80", "other.com:443", "other.com:443"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Graft(tt.base, tt.alt)
			if got != tt.expected {
				t.Errorf("Graft(%q, %q) = %q, want %q", tt.base, tt.alt, got, tt.expected)
			}
		})
	}
}

func TestParseDomainCandidates(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{"simple path", "/a/b/c", []string{"a", "b", "c"}},
		{"empty path", "/", nil},
		{"single segment", "/foo", []string{"foo"}},
		{"path with spaces", "/ a / b ", []string{"a", "b"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseDomainCandidates(tt.input)
			if len(got) != len(tt.expected) {
				t.Errorf("ParseDomainCandidates(%q) = %v, want %v", tt.input, got, tt.expected)
				return
			}
			for i := range got {
				if got[i] != tt.expected[i] {
					t.Errorf("ParseDomainCandidates(%q)[%d] = %q, want %q", tt.input, i, got[i], tt.expected[i])
				}
			}
		})
	}
}

func TestRealIP(t *testing.T) {
	tests := []struct {
		name       string
		remoteAddr string
		headers    map[string]string
		expected   string
	}{
		{
			name:       "X-Envoy-External-Address takes priority",
			remoteAddr: "1.2.3.4:5000",
			headers:    map[string]string{"X-Envoy-External-Address": "10.0.0.1", "X-Real-Ip": "10.0.0.2"},
			expected:   "10.0.0.1",
		},
		{
			name:       "X-Real-Ip fallback",
			remoteAddr: "1.2.3.4:5000",
			headers:    map[string]string{"X-Real-Ip": "10.0.0.2"},
			expected:   "10.0.0.2",
		},
		{
			name:       "RemoteAddr fallback",
			remoteAddr: "1.2.3.4:5000",
			headers:    map[string]string{},
			expected:   "1.2.3.4",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &http.Request{
				RemoteAddr: tt.remoteAddr,
				Header:     make(http.Header),
			}
			for k, v := range tt.headers {
				r.Header.Set(k, v)
			}
			got := RealIP(r)
			if got != tt.expected {
				t.Errorf("RealIP() = %q, want %q", got, tt.expected)
			}
		})
	}
}
