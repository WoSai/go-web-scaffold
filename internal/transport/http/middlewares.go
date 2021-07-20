package http

import (
	"context"
	"net/http"
	"time"

	"github.com/uber/jaeger-client-go"

	"github.com/opentracing/opentracing-go/ext"

	"github.com/opentracing/opentracing-go"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

type (
	ZapLogger struct {
		logger    *zap.Logger
		requestID string
	}
)

func NewZapLogEntry(logger *zap.Logger, r *http.Request) *ZapLogger {
	entry := &ZapLogger{
		logger: logger,
	}
	entry.requestID = middleware.GetReqID(r.Context())

	logger.Info("request started",
		zap.String("method", r.Method),
		zap.String("uri", r.RequestURI),
		zap.String("remote_addr", r.RemoteAddr),
		zap.String("user_agent", r.UserAgent()),
		zap.String("request_id", entry.requestID),
	)
	return entry
}

func (log *ZapLogger) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	log.logger.Info("request complete",
		zap.Int("response_status_code", status),
		zap.Int("response_bytes_length", bytes),
		zap.String("elapsed", elapsed.String()),
		zap.String("request_id", log.requestID),
	)
}

func (log *ZapLogger) Panic(v interface{}, stack []byte) {
	log.logger.Error("broken request",
		zap.Any("panic", v),
		zap.ByteString("stack", stack),
		zap.String("request_id", log.requestID),
	)
}

func ZapLog(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			entry := NewZapLogEntry(logger, r)
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			t1 := time.Now()

			next.ServeHTTP(ww, middleware.WithLogEntry(r, entry))

			entry.Write(ww.Status(), ww.BytesWritten(), ww.Header(), time.Since(t1), nil)
		}
		return http.HandlerFunc(fn)
	}
}

func Trace(tracer opentracing.Tracer) func(http.Handler) http.Handler {
	extract := jaeger.NewHTTPHeaderPropagator(&jaeger.HeadersConfig{}, jaeger.Metrics{})

	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			t1 := time.Now()
			var span opentracing.Span
			var ctx context.Context

			serverSpanCtx, err := extract.Extract(opentracing.HTTPHeaders)
			if err != nil {
				span, ctx = opentracing.StartSpanFromContext(r.Context(), "http request")
			} else {
				span, ctx = opentracing.StartSpanFromContextWithTracer(r.Context(), tracer, "http request", ext.RPCServerOption(serverSpanCtx))
			}

			defer span.Finish()

			span.SetTag("service.name", "localhost")
			ext.HTTPMethod.Set(span, r.Method)
			ext.HTTPUrl.Set(span, r.URL.Path)

			// next
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(ww, r.WithContext(ctx))

			// after
			status := ww.Status()
			ext.HTTPStatusCode.Set(span, uint16(status))
			span.LogKV("latency", time.Since(t1).String())
		}
		return http.HandlerFunc(fn)
	}
}
