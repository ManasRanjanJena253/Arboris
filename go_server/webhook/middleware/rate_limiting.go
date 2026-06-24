package middleware

import (
	"net/http"
	"sync"

	"golang.org/x/time/rate"
)

func RateLimiter(limiters *sync.Map, rps rate.Limit, burst int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			installID, ok := ctx.Value(InstallationIDKey).(string)

			if !ok || installID == "" {
				http.Error(w, "Install ID not found", http.StatusInternalServerError)
				return
			}
			limiterIface, _ := limiters.LoadOrStore(installID, rate.NewLimiter(rps, burst))

			limiter := limiterIface.(*rate.Limiter)

			if !limiter.Allow() {
				w.Header().Set("Retry-After", "1")
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
