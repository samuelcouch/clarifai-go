package main

import (
	"fmt"
	"html"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"

	"github.com/go-kit/kit/endpoint"
	//httptransport "github.com/go-kit/kit/transport/http"
	"golang.org/x/net/context"
)

/*
type proxymw struct {
	context.Context
	ClarifaiAPIService // Serve most stuff with this embedded service.
	ProxyEndpoint      endpoint.Endpoint
}

func proxyMiddleware(proxyURL string, ctx context.Context) ServiceMiddleware {
	return func(next ClarifaiAPIService) ClarifaiAPIService {
		return proxymw{ctx, next, makeProxyEndpoint(ctx, proxyURL)}
	}
}

func makeProxyEndpoint(ctx context.Context, proxyURL string) endpoint.Endpoint {
	u, err := url.Parse(proxyURL)
	if err != nil {
		panic(err)
	}
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return httptransport.NewClient(
			"GET",
			u,
			encodeRequestToAPIv1,
			decodeAPIv1Response,
		).Endpoint()
	}
}

// FIXME: can this be generic, or do we need per-endpoint codec?
func encodeRequestToAPIV1() {} // FIXME
*/

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
//
// This differs from httptransport.NewClient, which does a similar thing, but expects
// the proxied-to service to have a response schema.  This version passes through
// requests and responses without trying to parse them into a schema, it just plays back
// the bytes using an httptest.ResponseRecorder.
func makePassthroughProxyEndpoint(proxyHandler http.Handler) endpoint.Endpoint {
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
func decodePassthroughProxyRequest(r *http.Request) (interface{}, error) {
	return r, nil
}

// Pipe the recorded response from a ResponseRecorder to a new ResponseWriter.
func encodeFromRecordedResponse(w http.ResponseWriter, response interface{}) error {
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
