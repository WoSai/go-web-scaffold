package main

import (
	"net/http"
	"os"

	"github.com/wosai/go-web-scaffold/internal/log"
	"github.com/wosai/go-web-scaffold/internal/tracing"
	transporthttp "github.com/wosai/go-web-scaffold/internal/transport/http"
)

func main() {
	os.Setenv("JAEGER_SERVICE_NAME", "localhost")
	logger := log.BuildLogger(log.Option{Level: "debug"})
	tracing.BuildOpenTracing(logger)
	h := transporthttp.BuildRouter(tracing.Tracer)

	http.ListenAndServe(":8080", h)
}
