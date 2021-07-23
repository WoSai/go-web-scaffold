package main

import (
	"net/http"
	"os"

	"github.com/wosai/go-web-scaffold/internal/application/user"
	"github.com/wosai/go-web-scaffold/internal/log"
	"github.com/wosai/go-web-scaffold/internal/tracing"
	router "github.com/wosai/go-web-scaffold/internal/transport/http"
)

func main() {
	os.Setenv("JAEGER_SERVICE_NAME", "localhost")
	logger := log.BuildLogger(log.Option{Level: "info"})
	logger.Info("hello", "world")
	tracing.BuildOpenTracing(logger.(*log.KitZapLogger).Zap)
	app := user.BuildApplication(nil)

	h := router.BuildRouter(app, tracing.Tracer, logger.(*log.KitZapLogger).Kit)
	http.ListenAndServe(":8080", h)
	//http.ListenAndServe(":8080", h)
}
