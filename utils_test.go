package utils

import (
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
		{"[::1]:9000", ":9000"},
		{"[::1]", ""},
	}
	for _, tt := range tests {
		got := ExtractPort(tt.input)
		if got != tt.want {
			t.Errorf("ExtractPort(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestExtractURLPort(t *testing.T) {
	parse := func(raw string) *url.URL {
		u, _ := url.Parse(raw)
		return u
	}
	tests := []struct {
		u    *url.URL
		want string
	}{
		{nil, ""},
		{parse("http://example.com"), ""},
		{parse("http://example.com:8080"), ":8080"},
		{parse("http://[::1]:9000"), ":9000"},
	}
	for _, tt := range tests {
		got := ExtractURLPort(tt.u)
		if got != tt.want {
			t.Errorf("ExtractURLPort(%v) = %q, want %q", tt.u, got, tt.want)
		}
	}
}

func TestStripPort(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"example.com:8080", "example.com"},
		{"example.com", "example.com"},
		{"[::1]:9000", "::1"},
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
		{"/a/b/cd", []string{"a", "b", "cd"}},
		{"/", nil},
		{"/foo", []string{"foo"}},
		{"", nil},
		{"/a/ b /c", []string{"a", "b", "c"}},
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

func TestGraft(t *testing.T) {
	tests := []struct {
		base string
		alt  string
		want string
	}{
		{"example.com:8080", ":9090", "example.com:9090"},
		{"example.com", ":9090", "example.com:9090"},
		{"example.com:8080", "other.com:9090", "other.com:9090"},
		{"example.com:8080", "nocolon", "example.com:8080"},
	}
	for _, tt := range tests {
		got := Graft(tt.base, tt.alt)
		if got != tt.want {
			t.Errorf("Graft(%q, %q) = %q, want %q", tt.base, tt.alt, got, tt.want)
		}
	}
}
