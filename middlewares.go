package main

import (
	"net/http"
	"strings"
)

func Print(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/static/views/") {
			http.NotFound(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}
