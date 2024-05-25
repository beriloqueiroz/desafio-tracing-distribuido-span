package webserver

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type HandlerFuncMethod struct {
	HandleFunc http.HandlerFunc
	Method     string
}

type WebServer struct {
	Handlers      map[string]http.HandlerFunc
	WebServerPort string
}

func NewWebServer(serverPort string) *WebServer {
	return &WebServer{
		Handlers:      make(map[string]http.HandlerFunc),
		WebServerPort: serverPort,
	}
}

func (s *WebServer) AddRoute(path string, handler http.HandlerFunc) {
	s.Handlers[path] = handler
}

func (s *WebServer) Start() error {
	mux := http.NewServeMux()
	for path, handler := range s.Handlers {
		mux.Handle(path, otelhttp.WithRouteTag(path, http.HandlerFunc(handler)))
	}
	mux.Handle("GET /metrics", otelhttp.WithRouteTag("GET /metrics", promhttp.Handler()))

	return http.ListenAndServe(s.WebServerPort, mux)
}
