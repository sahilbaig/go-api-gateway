package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	customMiddleware "github.com/sahilbaig/go-api-gateway/middleware"
)

func main(){
	fmt.Println("Holla Holla get some Dolla")
	r:=chi.NewRouter()
	r.Use(middleware.DefaultLogger)
	r.Use(customMiddleware.RateLimiter(1,2))
	r.Get("/" , func (w http.ResponseWriter , r *http.Request)  {
		w.Write([]byte("Holla Holla "))
	})
	// 
	fmt.Println("Server starting on :8080...")

	http.ListenAndServe(":8080" , r)
}