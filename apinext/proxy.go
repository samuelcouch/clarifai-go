package main

import (
	"fmt"
	"html"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"

	"github.com/go-kit/kit/endpoint"
	"golang.org/x/net/context"
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

// Returns an Endpoint that proxies to the specified uri.
//
// The goal is to wrap the proxy handler in an endpoint, so we can use
// standard middleware for logging, alerting, throttling, etc.
func NewProxyEndpoint(proxyHandler http.Handler) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		r := request.(*http.Request)

		// Remove the proxy Host header.
		// FIXME: Need this?
		//r.Host = r.URL.Host

		w := httptest.NewRecorder()
		proxyHandler.ServeHTTP(w, r)
		return w, nil
	}
}

// No-op:  a proxy request is just the original http.Request.
func decodeProxyRequest(r *http.Request) (interface{}, error) {
	return r, nil
}

// Pipe the recorded response from a ResponseRecorder to a new ResponseWriter.
func encodeRecordedResponse(w http.ResponseWriter, response interface{}) error {
	rec := response.(*httptest.ResponseRecorder) // FIXME pointer?
	for k, v := range rec.Header() {
		w.Header()[k] = v
	}
	// Just testing ability to mess with the response when proxying:
	w.Header().Set("X-Clarifai-Proxied", "yes")
	w.WriteHeader(rec.Code)
	w.Write(rec.Body.Bytes())
	return nil
}
