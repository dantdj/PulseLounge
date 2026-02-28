package main

import (
	"io/fs"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path"
	"strings"
)

func spaHandler(uiFS fs.FS) http.HandlerFunc {
	fileServer := http.FileServer(http.FS(uiFS))

	return func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") {
			http.NotFound(w, r)
			return
		}

		cleanPath := strings.TrimPrefix(path.Clean(r.URL.Path), "/")
		if cleanPath == "." {
			cleanPath = ""
		}

		if cleanPath != "" {
			if _, err := fs.Stat(uiFS, cleanPath); err == nil {
				fileServer.ServeHTTP(w, r)
				return
			}
		}

		indexHTML, err := fs.ReadFile(uiFS, "index.html")
		if err != nil {
			http.Error(w, "index file not found", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(indexHTML)
	}
}

// spaProxyHandler returns an http.HandlerFunc that proxies requests to the specified target URL.
// It is used to forward requests to a development server when the UI_DEV_SERVER environment variable is set.
func spaProxyHandler(target string) (http.HandlerFunc, error) {
	u, err := url.Parse(target)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(u)

	return func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") {
			http.NotFound(w, r)
			return
		}
		proxy.ServeHTTP(w, r)
	}, nil
}
