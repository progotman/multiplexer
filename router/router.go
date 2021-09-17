package router

import "net/http"

type Router struct {
	MaxNumberOfIncomingRequests int
	Routes                      []Route
}

func (s *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	guard := make(chan struct{}, s.MaxNumberOfIncomingRequests)
	guard <- struct{}{}
	defer func() { <-guard }()

	for _, routeHandler := range s.Routes {
		if routeHandler.Supports(r) {
			routeHandler.Handler.ServeHTTP(w, r)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
}

type Route struct {
	Url     string
	Method  string
	Handler http.Handler
}

func (h *Route) Supports(r *http.Request) bool {
	return r.Method == h.Method && r.URL.Path == h.Url
}
