package server

import (
	"fmt"
	"net/http"
)

func (s *Server) PUTI(route string, h func(w http.ResponseWriter, r *http.Request), opts ...RouteOption) {
	routeInfo := initRouteInfo()
	for _, opt := range opts {
		opt(&routeInfo)
	}

	allMiddlewares := append(routeInfo.Middlewares, putRequest)

	s.mux.Handle(route, chainMiddleware(http.HandlerFunc(h), allMiddlewares...))
	s.Paths = append(s.Paths, ServerPath{
		Route:  route,
		Method: METHOD_PUT,
		Info:   routeInfo,
	})
}

func (s *Server) PUT(route string, h func(w http.ResponseWriter, r *http.Request), opts ...RouteOption) {
	routeInfo := initRouteInfo()
	for _, opt := range opts {
		opt(&routeInfo)
	}

	allMiddlewares := append(routeInfo.Middlewares, putRequest)

	if !s.TranslationsEnabled {
		s.mux.Handle(route, chainMiddleware(http.HandlerFunc(h), allMiddlewares...))
		s.Paths = append(s.Paths, ServerPath{
			Route:  route,
			Method: METHOD_PUT,
			Info:   routeInfo,
		})
	} else {
		s.mux.Handle(route, chainMiddleware(http.HandlerFunc(s.redirectToTranslatedUrl), allMiddlewares...))
		s.Paths = append(s.Paths, ServerPath{
			Route:  route,
			Method: METHOD_PUT,
			Info:   routeInfo,
		})

		for short, _ := range s.Languages {
			r := fmt.Sprintf("/%s%s", short, route)
			s.mux.Handle(r, chainMiddleware(http.HandlerFunc(h), allMiddlewares...))
			s.Paths = append(s.Paths, ServerPath{
				Route:  r,
				Method: METHOD_PUT,
				Info:   routeInfo,
			})
		}
	}
}
