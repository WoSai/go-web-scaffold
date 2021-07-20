package tracing

import (
	"io"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
	jeagerlog "github.com/uber/jaeger-client-go/log/zap"
	"go.uber.org/zap"
)

var (
	Tracer opentracing.Tracer
	Closer io.Closer
)

func BuildOpenTracing(log *zap.Logger, opts ...config.Option) {
	cfg, err := config.FromEnv()
	if err != nil {
		panic(err)
	}

	tracer, closer, err := cfg.NewTracer(opts...)
	if err != nil {
		panic(err)
	}

	tracer = jeagerlog.NewLoggingTracer(log, tracer)
	opentracing.SetGlobalTracer(tracer)

	Tracer = tracer
	Closer = closer
}
