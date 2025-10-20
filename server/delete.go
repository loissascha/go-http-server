package server

import (
	"fmt"
	"net/http"
)

// DELETE route (ignores translation rules)
func (s *Server) DELETEI(route string, h func(w http.ResponseWriter, r *http.Request), opts ...RouteOption) {
	routeInfo := initRouteInfo()
	for _, opt := range opts {
		opt(&routeInfo)
	}
	allMiddlewares := append(routeInfo.Middlewares, deleteRequest)

	s.mux.Handle(route, chainMiddleware(http.HandlerFunc(h), allMiddlewares...))
	s.Paths = append(s.Paths, ServerPath{
		Route:  route,
		Method: METHOD_DELETE,
		Info:   routeInfo,
	})
}

// DELETE route
func (s *Server) DELETE(route string, h func(w http.ResponseWriter, r *http.Request), opts ...RouteOption) {
	routeInfo := initRouteInfo()
	for _, opt := range opts {
		opt(&routeInfo)
	}
	allMiddlewares := append(routeInfo.Middlewares, deleteRequest)

	if !s.TranslationsEnabled {
		s.mux.Handle(route, chainMiddleware(http.HandlerFunc(h), allMiddlewares...))
		s.Paths = append(s.Paths, ServerPath{
			Route:  route,
			Method: METHOD_DELETE,
			Info:   routeInfo,
		})
	} else {
		s.mux.Handle(route, chainMiddleware(http.HandlerFunc(s.redirectToTranslatedUrl), allMiddlewares...))
		s.Paths = append(s.Paths, ServerPath{
			Route:  route,
			Method: METHOD_DELETE,
			Info:   routeInfo,
		})

		for short, _ := range s.Languages {
			r := fmt.Sprintf("/%s%s", short, route)
			s.mux.Handle(r, chainMiddleware(http.HandlerFunc(h), allMiddlewares...))
			s.Paths = append(s.Paths, ServerPath{
				Route:  r,
				Method: METHOD_DELETE,
				Info:   routeInfo,
			})
		}
	}
}
