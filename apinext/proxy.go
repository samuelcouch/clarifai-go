package main

import (
	"fmt"
	"html"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// TODO(madadam): Make this into an Endpoint/service that we can wrap with middleware.

// Returns an http.Handler that proxies to the specified uri.
func NewProxy(uri string) (http.Handler, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxywrapper := func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Proxying %s\n", html.EscapeString(r.URL.Path)) // FIXME
		// Remove the proxy Host header.
		r.Host = r.URL.Host
		proxy.ServeHTTP(w, r)
	}
	return http.HandlerFunc(proxywrapper), nil
}
