package handler

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/laWiki/gateway/config"
)

// ReverseProxy is a handler that takes a target host and proxies requests to it
func ReverseProxy(target string, prefixToStrip string) func(http.ResponseWriter, *http.Request) {
	targetURL, err := url.Parse(target)
	if err != nil {
		config.App.Logger.Panic().Msg("Invalid proxy target URL")
	}

	config.App.Logger.Info().Interface("targetURL", targetURL)

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// Modify the request before sending it to the backend
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		config.App.Logger.Debug().
			Str("original_url", req.URL.String()).
			Str("original_host", req.Host).
			Msg("Original request details")
		// Set the scheme and host to the target's scheme and host
		req.URL.Scheme = targetURL.Scheme
		req.URL.Host = targetURL.Host

		// Strip the specified prefix from the request URL Path
		req.URL.Path = strings.TrimPrefix(req.URL.Path, prefixToStrip)

		// Ensure the path is not empty
		if req.URL.Path == "" {
			req.URL.Path = "/"
		}

		// Update the request Host header to the target host
		req.Host = targetURL.Host
		config.App.Logger.Debug().
			Str("modified_url", req.URL.String()).
			Str("modified_host", req.Host).
			Msg("Modified request details")
	}

	proxy.ModifyResponse = func(resp *http.Response) error {
		return nil
	}

	// Error handler
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		config.App.Logger.Error().Err(err).Msg("Proxy error")
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte("Bad Gateway"))
	}

	return func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	}
}
