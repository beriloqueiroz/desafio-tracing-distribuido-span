package webserver

import (
	"net/http"
	"time"

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
	OtelTracer    trace.Tracer
}

func NewWebServer(serverPort string, otelTracer trace.Tracer) *WebServer {
	return &WebServer{
		Handlers:      make(map[string]http.HandlerFunc),
		WebServerPort: serverPort,
		OtelTracer:    otelTracer,
	}
}

func (s *WebServer) AddRoute(path string, handler http.HandlerFunc) {
	s.Handlers[path] = handler
}

func (s *WebServer) Start() error {
	mux := http.NewServeMux()
	for path, handler := range s.Handlers {
		mux.Handle(path, otelhttp.WithRouteTag(path, s.OtelMiddleware(http.HandlerFunc(handler))))
	}
	mux.Handle("GET /metrics", otelhttp.WithRouteTag("GET /metrics", promhttp.Handler()))

	return http.ListenAndServe(s.WebServerPort, mux)
}

func (s *WebServer) OtelMiddleware(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		carrier := propagation.HeaderCarrier(r.Header)
		ctx := r.Context()
		ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)
		startTime := time.Now()
		name := r.URL.Path
		ctx, span := s.OtelTracer.Start(ctx, name, trace.WithTimestamp(startTime))
		defer span.End()
		otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(r.Header))
		handler.ServeHTTP(w, r)
		span.End(trace.WithTimestamp(time.Now()))
	}
}
