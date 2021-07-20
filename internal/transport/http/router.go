package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/opentracing/opentracing-go"
)

func BuildRouter(tracer opentracing.Tracer) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.Recoverer)
	router.Use(Trace(tracer))

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hello wosai"))
	})
	return router
}
