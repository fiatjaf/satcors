package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

type corsTransport struct {
	referer     string
	origin      string
	credentials string
}

func (t corsTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	// Put in the Referer if specified
	if t.referer != "" {
		r.Header.Add("Referer", t.referer)
	}

	// Do the actual request
	res, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		return nil, err
	}

	res.Header.Set("Access-Control-Allow-Origin", t.origin)
	res.Header.Set("Access-Control-Allow-Credentials", t.credentials)

	return res, nil
}

func handleProxy(w http.ResponseWriter, r *http.Request, origin string, credentials string) {
	// Check for the User-Agent header
	if r.Header.Get("User-Agent") == "" {
		http.Error(w, "Missing User-Agent header", http.StatusBadRequest)
		return
	}

	// Get the optional Referer header
	referer := r.URL.Query().Get("referer")
	if referer == "" {
		referer = r.Header.Get("Referer")
	}

	// Get the URL
	urlParam := r.URL.Query().Get("url")
	// Validate the URL
	urlParsed, err := url.Parse(urlParam)
	if err != nil {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	// Check if HTTP(S)
	if urlParsed.Scheme != "http" && urlParsed.Scheme != "https" {
		http.Error(w, "The URL scheme is neither HTTP nor HTTPS", http.StatusBadRequest)
		return
	}

	// Setup for the proxy
	proxy := httputil.ReverseProxy{
		Director: func(r *http.Request) {
			r.URL = urlParsed
			r.Host = urlParsed.Host
		},
		Transport: corsTransport{referer, origin, credentials},
	}

	// Execute the request
	proxy.ServeHTTP(w, r)
}

func HandleProxy(w http.ResponseWriter, r *http.Request) {
	shouldBeAllowed := checkRequest(r.Header.Get("Referer"))

	if !shouldBeAllowed {
		w.WriteHeader(402)
		return
	}

	handleProxy(w, r, "*", "true")
}
