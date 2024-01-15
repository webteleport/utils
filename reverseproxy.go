package utils

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// ReverseProxy
//
//   - upstream: http://user:pass@example.com
//     behavior: use upstream credential by default
//   - upstream: http://example.com
//     behavior: passthrough client credential if any
//   - upstream: http://-@example.com
//     behavior: don't pass any credential to upstream
func ReverseProxy(addr string) http.Handler {
	// google.com will be parsed as URL{Path: google.com} without an explicit protocol
	// hence the hack
	if !strings.Contains(addr, "://") {
		addr = "http://" + addr
	}
	upstream, err := url.Parse(addr)
	if err != nil {
		panic(err)
	}
	rewrite := func(r *httputil.ProxyRequest) {
		r.SetURL(upstream)
		r.SetXForwarded()

		// ignore credential when default upstream user set to -
		if upstream.User.String() == "-" {
			r.Out.Header.Del("Authorization")
			return
		}

		// use client credentials if available
		if r.In.Header.Get("Authorization") != "" {
			return
		}

		// otherwise use credentials from upstream
		if upstream.User != nil {
			user := upstream.User.Username()
			pass, _ := upstream.User.Password()
			r.Out.SetBasicAuth(user, pass)
		}
	}
	modify := func(r *http.Response) error {
		r.Header.Del("Content-Security-Policy")
		return nil
	}
	rp := &httputil.ReverseProxy{
		Rewrite:        rewrite,
		ModifyResponse: modify,
	}
	return rp
}
