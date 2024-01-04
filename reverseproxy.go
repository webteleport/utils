package utils

import (
	"net/http"
	"net/http/httputil"
	"net/url"
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
	rp := &httputil.ReverseProxy{
		Rewrite: rewrite,
	}
	return rp
}
