package server

import (
	"fmt"
	"net/http"
)

func (s *Server) GET(route string, h func(w http.ResponseWriter, r *http.Request), opts ...RouteOption) {
	routeInfo := initRouteInfo()
	for _, opt := range opts {
		opt(&routeInfo)
	}

	allMiddlewares := append(routeInfo.Middlewares, getRequest)

	if !s.TranslationsEnabled {
		s.mux.Handle(route, chainMiddleware(http.HandlerFunc(h), allMiddlewares...))
		s.Paths = append(s.Paths, ServerPath{
			Route:  route,
			Method: METHOD_GET,
			Info:   routeInfo,
		})
	} else {
		// add multiple routes
		// base route -> redirect to default language
		// one route per language (and somehow pass down the used language?)

		s.mux.Handle(route, chainMiddleware(http.HandlerFunc(s.redirectToTranslatedUrl), allMiddlewares...))
		s.Paths = append(s.Paths, ServerPath{
			Route:  route,
			Method: METHOD_GET,
			Info:   routeInfo,
		})

		for short, _ := range s.Languages {
			r := fmt.Sprintf("/%s%s", short, route)
			s.mux.Handle(r, chainMiddleware(http.HandlerFunc(h), allMiddlewares...))
			s.Paths = append(s.Paths, ServerPath{
				Route:  r,
				Method: METHOD_GET,
				Info:   routeInfo,
			})
		}
	}
}
