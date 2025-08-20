package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sahilbaig/go-api-gateway/discovery"
	customMiddleware "github.com/sahilbaig/go-api-gateway/middleware"
)

func main(){
	fmt.Println("Holla Holla get some Dolla")
	r:=chi.NewRouter()
	

	data :=discovery.ServiceDiscovery();
	fmt.Println(data)
	
	r.Use(middleware.DefaultLogger)
	r.Use(customMiddleware.RateLimiter(1,2))
	r.Use(func(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Execute the request first
        next.ServeHTTP(w, r)

        // Then log the matched route
        if ctx := chi.RouteContext(r.Context()); ctx != nil {
            fmt.Println("Matched route:", ctx.RoutePattern())
        }
    })
})
	
	
	r.Get("/" , func (w http.ResponseWriter , r *http.Request)  {
		w.Write([]byte("Holla Holla "))
	})
	
	
   
	
	fmt.Println("Server starting on :8080...")

	http.ListenAndServe(":8080" , r)
}