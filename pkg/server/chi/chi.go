package gchi

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"mwz.com/data"
	"mwz.com/service"
)

func Server(port int) {
	r := chi.NewRouter()
	// r.Use(middleware.Logger)

	// A good base middleware stack
	// r.Use(middleware.RequestID)
	// r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(1 * time.Second))
	// r.Use(middleware.Timeout(5 * time.Millisecond))

	r.Use(func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	d, _, _ := data.NewData()
	testService := service.NewTestService(d)
	r.Get("/ca", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(testService.GetRedis(r.Context())))
	})

	http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}
