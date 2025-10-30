package server

import (
	"fmt"
	"net/http"
)

// POST route (ignores translation rules)
func (s *Server) POSTI(route string, h func(w http.ResponseWriter, r *http.Request), opts ...RouteOption) {
	routeInfo := getRouteInfos(opts...)

	s.addPath(route, ServerPath{
		Route:   route,
		Method:  METHOD_POST,
		Info:    routeInfo,
		Handler: h,
	})
}

// POST route
func (s *Server) POST(route string, h func(w http.ResponseWriter, r *http.Request), opts ...RouteOption) {
	if !s.TranslationsEnabled {
		s.POSTI(route, h, opts...)
		return
	}

	routeInfo := getRouteInfos(opts...)
	s.addPath(route, ServerPath{
		Route:   route,
		Method:  METHOD_POST,
		Info:    routeInfo,
		Handler: s.redirectToTranslatedUrl,
	})

	for short := range s.Languages {
		r := fmt.Sprintf("/%s%s", short, route)
		s.addPath(r, ServerPath{
			Route:   r,
			Method:  METHOD_POST,
			Info:    routeInfo,
			Handler: h,
		})
	}
}
