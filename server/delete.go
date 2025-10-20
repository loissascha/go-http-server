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

	s.addPath(route, ServerPath{
		Route:   route,
		Method:  METHOD_DELETE,
		Info:    routeInfo,
		Handler: h,
	})
}

// DELETE route
func (s *Server) DELETE(route string, h func(w http.ResponseWriter, r *http.Request), opts ...RouteOption) {
	routeInfo := initRouteInfo()
	for _, opt := range opts {
		opt(&routeInfo)
	}

	if !s.TranslationsEnabled {
		s.addPath(route, ServerPath{
			Route:   route,
			Method:  METHOD_DELETE,
			Info:    routeInfo,
			Handler: h,
		})
	} else {
		s.addPath(route, ServerPath{
			Route:   route,
			Method:  METHOD_DELETE,
			Info:    routeInfo,
			Handler: s.redirectToTranslatedUrl,
		})

		for short := range s.Languages {
			r := fmt.Sprintf("/%s%s", short, route)
			s.addPath(r, ServerPath{
				Route:   r,
				Method:  METHOD_DELETE,
				Info:    routeInfo,
				Handler: h,
			})
		}
	}
}
