package utils

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/miekg/dns"
	"golang.org/x/net/idna"
)

func UnwrapInnermost(err error) error {
	for {
		unwrapped := errors.Unwrap(err)
		if unwrapped == nil {
			return err
		}
		err = unwrapped
	}
}

func LookupHostTXT(domain, server string) ([]string, error) {
	answers := []string{}

	// Create a new DNS client
	client := dns.Client{}

	// Create a new DNS message
	message := dns.Msg{}
	message.SetQuestion(dns.Fqdn(domain), dns.TypeTXT)

	// Send the DNS query
	response, _, err := client.Exchange(&message, server)
	if err != nil {
		return answers, fmt.Errorf("DNS query failed: %w\n", err)
	}

	// Process the DNS response
	if response.Rcode != dns.RcodeSuccess {
		return answers, fmt.Errorf("DNS query failed with response code: %s\n", dns.RcodeToString[response.Rcode])
	}

	// Extract and print the TXT records
	for _, answer := range response.Answer {
		if txt, ok := answer.(*dns.TXT); ok {
			a := strings.ReplaceAll(txt.Txt[0], `\`, "")
			answers = append(answers, a)
		}
	}
	return answers, nil
}

// ExtractURLPort returns the :port part from URL.Host (host[:port])
//
// An empty string is returned if no port is found
func ExtractURLPort(u *url.URL) string {
	_, p, ok := strings.Cut(u.Host, ":")
	if ok {
		return ":" + p
	}
	return ""
}

// ToIdna converts a string to its idna form at best effort
// Should only be used on the hostname part without port
func ToIdna(s string) string {
	ascii, err := idna.ToASCII(s)
	if err != nil {
		log.Println(err)
		return s
	}
	return ascii
}

// Graft returns Host(base):Port(alt)
//
// assuming
// - base is host[:port]
// - alt is [host]:port
func Graft(base, alt string) string {
	althost, altport, _ := strings.Cut(alt, ":")
	if altport == "" {
		// altport not found
		// it should never happen
		return base
	}
	if althost != "" {
		// alt is host:port
		// it is rare
		return alt
	}
	basehost, _, _ := strings.Cut(base, ":")
	return basehost + ":" + altport
}

// Print Hyperlink via OSC 8 ansi sequence.
// The syntax is: 'OSC 8 ; params ; url ST text OSC 8 ; ; ST'
// for more info see https://gist.github.com/egmontkob/eb114294efbcd5adb1944c9f3cb5feda
func Hyperlink(name, url string) string {
	return fmt.Sprintf("\u001B]8;%s;%s\u001B\\%s\u001B]8;;\u001B\\", "", url, name)
}

// MaybeHyperlink turns input into ANSI hyperlink when stdin is a tty
func MaybeHyperlink(l string) string {
	if Isatty() {
		return Hyperlink(l, l)
	}
	return l
}

// ParseDomainCandidates splits a path string like /a/b/cd/😏
// into a list of subdomains: [a, b, cd, 😏]
//
// when result is empty, a random subdomain will be assigned by the server
func ParseDomainCandidates(p string) []string {
	var list []string
	parts := strings.Split(p, "/")
	for _, part := range parts {
		dom := strings.Trim(part, " ")
		if dom == "" {
			continue
		}
		list = append(list, dom)
	}
	return list
}

func StripPort(hostport string) string {
	// use net.SplitHostPort instead of strings.Split
	// because it can handle ipv6 addresses
	host, _, err := net.SplitHostPort(hostport)
	if err != nil {
		// if there is no port, just return the input
		return hostport
	}
	return host
}

func RealIP(r *http.Request) (clientIP string) {
	// Retrieve the client IP address from the request headers
	for _, x := range []string{
		r.Header.Get("X-Envoy-External-Address"),
		r.Header.Get("X-Real-IP"),
		r.Header.Get("X-Forwarded-For"),
		StripPort(r.RemoteAddr),
	} {
		if x != "" {
			clientIP = x
			break
		}
	}
	return
}
