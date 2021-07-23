package http

import (
	"context"
	"encoding/json"
	"net/http"

	kitlog "github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/opentracing/opentracing-go"
	"github.com/wosai/go-web-scaffold/internal/application/user"
	"github.com/wosai/go-web-scaffold/internal/application/user/command"
)

func commonDecode(ctx context.Context, r *http.Request) (interface{}, error) {
	v := new(command.CreateUserRequest)
	err := json.NewDecoder(r.Body).Decode(v)
	return v, err
}

func commonEncode(ctx context.Context, w http.ResponseWriter, v interface{}) error {
	return json.NewEncoder(w).Encode(v)
}

func BuildRouter(app *user.Application, tracer opentracing.Tracer, logger kitlog.Logger) http.Handler {
	router := http.NewServeMux()

	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
	}

	router.Handle("/", httptransport.NewServer(
		app.CreateUser,
		commonDecode,
		commonEncode,
		append(options, httptransport.ServerBefore(JaegerHTTPToContext(tracer, "create user", logger)))...,
	))
	return router
}
