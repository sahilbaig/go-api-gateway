package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sahilbaig/go-api-gateway/discovery"
	customMiddleware "github.com/sahilbaig/go-api-gateway/middleware"
)

// newProxy creates a reverse proxy to the given target URL
func newProxy(target string) http.Handler {
	u, _ := url.Parse(target)
	return httputil.NewSingleHostReverseProxy(u)
}

func main() {
	fmt.Println("Holla Holla get some Dolla")
	r := chi.NewRouter()

	// Discover services and their pod IPs
	servicePods := discovery.ServiceDiscovery()

	r.Use(middleware.DefaultLogger)
	r.Use(customMiddleware.RateLimiter(1, 2))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Holla Holla "))
	})

	// Dynamically create proxies for each service
	for svc, pods := range servicePods {
		if len(pods) == 0 {
			continue
		}

		// For now, just pick the first pod
		target := fmt.Sprintf("http://%s:3000", pods[0]) // adjust port per service if needed
		path := fmt.Sprintf("/%s/*", svc)

		r.Handle(path, http.StripPrefix("/"+svc, newProxy(target)))
		fmt.Printf("Proxying /%s -> %s\n", svc, target)
	}

	fmt.Println("Server starting on :8080...")
	http.ListenAndServe(":8080", r)
}
