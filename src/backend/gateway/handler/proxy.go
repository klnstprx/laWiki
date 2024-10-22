package handler

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/laWiki/gateway/config"
)

// ReverseProxy is a handler that takes a target host and proxies requests to it
func ReverseProxy(target string) func(http.ResponseWriter, *http.Request) {
	targetURL, err := url.Parse(target)
	if err != nil {
		config.App.Logger.Panic().Msg("Invalid proxy target URL")
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// Modify the request before sending it to the backend
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		// You can modify the request here if needed
	}

	// Optional: Modify the response from the backend
	proxy.ModifyResponse = func(resp *http.Response) error {
		// You can modify the response here if needed
		return nil
	}

	// Error handler
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte("Bad Gateway"))
	}

	return func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	}
}
