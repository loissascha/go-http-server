package server

import (
	"fmt"
	"net/http"
)

// GET route (ignores translation rules)
func (s *Server) GETI(route string, h func(w http.ResponseWriter, r *http.Request), opts ...RouteOption) {
	routeInfo := getRouteInfos(opts...)

	s.addPath(route, ServerPath{
		Route:   route,
		Method:  METHOD_GET,
		Info:    routeInfo,
		Handler: h,
	})
}

// GET route
func (s *Server) GET(route string, h func(w http.ResponseWriter, r *http.Request), opts ...RouteOption) {
	if !s.TranslationsEnabled {
		s.GETI(route, h, opts...)
		return
	}

	routeInfo := getRouteInfos(opts...)
	s.addPath(route, ServerPath{
		Route:   route,
		Method:  METHOD_GET,
		Info:    routeInfo,
		Handler: s.redirectToTranslatedUrl,
	})

	for short := range s.Languages {
		r := fmt.Sprintf("/%s%s", short, route)
		s.addPath(r, ServerPath{
			Route:   r,
			Method:  METHOD_GET,
			Info:    routeInfo,
			Handler: h,
		})
	}
}
