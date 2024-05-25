package webserver

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type HandlerFuncMethod struct {
	HandleFunc http.HandlerFunc
	Method     string
}

type WebServer struct {
	Handlers      map[string]http.HandlerFunc
	WebServerPort string
	Tracer        trace.Tracer
}

func NewWebServer(serverPort string, tracer trace.Tracer) *WebServer {
	return &WebServer{
		Handlers:      make(map[string]http.HandlerFunc),
		WebServerPort: serverPort,
		Tracer:        tracer,
	}
}

func (s *WebServer) AddRoute(path string, handler http.HandlerFunc) {
	s.Handlers[path] = handler
}

func (s *WebServer) Start() error {
	mux := http.NewServeMux()
	for path, handler := range s.Handlers {
		mux.Handle(path, otelhttp.WithRouteTag(path, s.TelemetryMiddleware(http.HandlerFunc(handler))))
	}
	mux.Handle("GET /metrics", otelhttp.WithRouteTag("GET /metrics", promhttp.Handler()))

	return http.ListenAndServe(s.WebServerPort, mux)
}

func (s *WebServer) TelemetryMiddleware(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		carrier := propagation.HeaderCarrier(r.Header)
		ctx := r.Context()
		ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)
		ctx, span := s.Tracer.Start(ctx, r.URL.Path)
		defer span.End()
		*r = *r.WithContext(ctx)
		handler.ServeHTTP(w, r)
	}
}
