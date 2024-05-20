package webserver

import (
	"fmt"
	"net/http"
	"time"
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

		// mux.Handle(path, otelhttp.WithRouteTag(path, http.HandlerFunc(handler)))
		// mux.Handle(path, otelhttp.WithRouteTag(path, http.HandlerFunc(OtelMiddleware(handler))))
		mux.HandleFunc(path, Middleware(handler))
	}
	return http.ListenAndServe(s.WebServerPort, mux)
}

func Middleware(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Path
		delay := time.Now()
		fmt.Printf("path %s start\n", name)
		handler.ServeHTTP(w, r)
		fmt.Printf("path %s end at %d\n", name, time.Now().Sub(delay).Milliseconds())
	}
}
