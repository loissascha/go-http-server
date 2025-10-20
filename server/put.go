package server

import (
	"fmt"
	"net/http"
)

// PUT route (ignores translation rules)
func (s *Server) PUTI(route string, h func(w http.ResponseWriter, r *http.Request), opts ...RouteOption) {
	routeInfo := initRouteInfo()
	for _, opt := range opts {
		opt(&routeInfo)
	}

	s.addPath(route, ServerPath{
		Route:   route,
		Method:  METHOD_PUT,
		Info:    routeInfo,
		Handler: h,
	})
}

// PUT route
func (s *Server) PUT(route string, h func(w http.ResponseWriter, r *http.Request), opts ...RouteOption) {
	routeInfo := initRouteInfo()
	for _, opt := range opts {
		opt(&routeInfo)
	}

	if !s.TranslationsEnabled {
		s.addPath(route, ServerPath{
			Route:   route,
			Method:  METHOD_PUT,
			Info:    routeInfo,
			Handler: h,
		})
	} else {
		s.addPath(route, ServerPath{
			Route:   route,
			Method:  METHOD_PUT,
			Info:    routeInfo,
			Handler: s.redirectToTranslatedUrl,
		})

		for short := range s.Languages {
			r := fmt.Sprintf("/%s%s", short, route)
			s.addPath(r, ServerPath{
				Route:   r,
				Method:  METHOD_PUT,
				Info:    routeInfo,
				Handler: h,
			})
		}
	}
}
