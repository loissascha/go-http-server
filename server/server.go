package server

import "net/http"

type Method int

const (
	METHOD_GET    Method = 0
	METHOD_POST   Method = 1
	METHOD_PUT    Method = 2
	METHOD_DELETE Method = 3
)

type ServerPath struct {
	Route  string
	Method Method
	Info   RouteInfo
}

type Server struct {
	mux     *http.ServeMux
	Paths   []ServerPath
	Options []ServerOption
}

func NewServer(options ...ServerOption) *Server {
	s := Server{
		Paths:   []ServerPath{},
		Options: options,
	}
	s.mux = http.NewServeMux()
	return &s
}

func (s *Server) GetMux() *http.ServeMux {
	return s.mux
}

func (s *Server) Serve(addr string) error {
	err := http.ListenAndServe(addr, s.mux)
	return err
}

func (s *Server) Handle(route string, h func(w http.ResponseWriter, r *http.Request), middlewares ...func(http.Handler) http.Handler) {
	s.mux.Handle(route, chainMiddleware(http.HandlerFunc(h), middlewares...))
}

func initRouteInfo() RouteInfo {
	routeInfo := RouteInfo{
		Middlewares: []func(http.Handler) http.Handler{},
		Tags:        []string{},
		Params:      []OpenAPIParam{},
		Responses:   map[string]OpenAPIResponse{},
	}
	return routeInfo
}

func (s *Server) GET(route string, h func(w http.ResponseWriter, r *http.Request), opts ...RouteOption) {
	routeInfo := initRouteInfo()
	for _, opt := range opts {
		opt(&routeInfo)
	}

	allMiddlewares := append(routeInfo.Middlewares, getRequest)
	s.mux.Handle(route, chainMiddleware(http.HandlerFunc(h), allMiddlewares...))
	s.Paths = append(s.Paths, ServerPath{
		Route:  route,
		Method: METHOD_GET,
		Info:   routeInfo,
	})
}

func (s *Server) POST(route string, h func(w http.ResponseWriter, r *http.Request), opts ...RouteOption) {
	routeInfo := initRouteInfo()
	for _, opt := range opts {
		opt(&routeInfo)
	}

	allMiddlewares := append(routeInfo.Middlewares, postRequest)
	s.mux.Handle(route, chainMiddleware(http.HandlerFunc(h), allMiddlewares...))
	s.Paths = append(s.Paths, ServerPath{
		Route:  route,
		Method: METHOD_POST,
		Info:   routeInfo,
	})
}

func (s *Server) PUT(route string, h func(w http.ResponseWriter, r *http.Request), opts ...RouteOption) {
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

func (s *Server) DELETE(route string, h func(w http.ResponseWriter, r *http.Request), opts ...RouteOption) {
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
