package middleware

import (
	"net/http"

	"golang.org/x/time/rate"
)

func RateLimiter(rps int, burst int)func(http.Handler) http.Handler {
	limiter:= rate.NewLimiter(rate.Limit(rps) , burst)

	return  func (next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !limiter.Allow(){
				http.Error(w , "Rate Limit Exceeded" , http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w,r)
		})
	}
}