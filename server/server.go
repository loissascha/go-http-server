package server

import (
	"net/http"
	"strings"
)

type Method int

const (
	METHOD_GET    Method = 0
	METHOD_POST   Method = 1
	METHOD_PUT    Method = 2
	METHOD_DELETE Method = 3
)

type ServerPath struct {
	Route   string
	Method  Method
	Info    RouteInfo
	Handler func(w http.ResponseWriter, r *http.Request)
}

type Server struct {
	mux                       *http.ServeMux
	Paths                     map[string][]ServerPath
	Options                   []ServerOption
	TranslationsEnabled       bool
	AutoDetectLanguageEnabled bool
	Languages                 map[string]map[string]string
	DefaultLanguage           string
}

func (s *Server) addPath(route string, p ServerPath) {
	sp, found := s.Paths[route]
	if !found {
		s.Paths[route] = []ServerPath{
			p,
		}
		return
	}
	s.Paths[route] = append(sp, p)
}

func NewServer(options ...ServerOption) (*Server, error) {
	s := Server{
		Paths:     map[string][]ServerPath{},
		Options:   options,
		Languages: map[string]map[string]string{},
	}
	s.initServerOptions()
	s.mux = http.NewServeMux()
	return &s, nil
}

func (s *Server) GetActiveLanguage(r *http.Request) string {
	url := strings.TrimSpace(r.URL.Path)
	url = strings.TrimLeft(url, "/")
	urlSplit := strings.Split(url, "/")
	activeLanguage := urlSplit[0]
	return activeLanguage
}

func (s *Server) GetLanguageStringMap(r *http.Request) map[string]string {
	// get currently active language (from url)
	url := strings.TrimSpace(r.URL.Path)
	url = strings.TrimLeft(url, "/")
	urlSplit := strings.Split(url, "/")
	activeLanguage := urlSplit[0]

	// find the value for the key
	l, ok := s.Languages[activeLanguage]
	if !ok {
		// take default
		l, ok = s.Languages[s.DefaultLanguage]
		if !ok {
			panic("Server default language configuration failed")
		}
		return l
	}

	return l
}

func (s *Server) GetLanguageString(r *http.Request, key string) string {
	// get currently active language (from url)
	url := strings.TrimSpace(r.URL.Path)
	url = strings.TrimLeft(url, "/")
	urlSplit := strings.Split(url, "/")
	activeLanguage := urlSplit[0]

	// find the value for the key
	l, ok := s.Languages[activeLanguage]
	if !ok {
		// take default
		l, ok = s.Languages[s.DefaultLanguage]
		if !ok {
			panic("Server default language configuration failed")
		}
		return l[key]
	}

	value, valueFound := l[key]

	// if not found -> take from default
	if !valueFound {
		l, ok = s.Languages[s.DefaultLanguage]
		if !ok {
			panic("Server default language configuration failed")
		}
		return l[key]
	}

	return value
}

func (s *Server) GetMux() *http.ServeMux {
	return s.mux
}

func (s *Server) Serve(addr string) error {

	for route, serverPaths := range s.Paths {
		getM := []func(http.Handler) http.Handler{}
		var getH *func(http.ResponseWriter, *http.Request)
		postM := []func(http.Handler) http.Handler{}
		var postH *func(http.ResponseWriter, *http.Request)
		putM := []func(http.Handler) http.Handler{}
		var putH *func(http.ResponseWriter, *http.Request)
		deleteM := []func(http.Handler) http.Handler{}
		var deleteH *func(http.ResponseWriter, *http.Request)
		for _, path := range serverPaths {
			switch path.Method {
			case METHOD_GET:
				if getH != nil {
					panic("double GET route!")
				}
				allMiddlewares := append(path.Info.Middlewares, getRequest)
				getM = allMiddlewares
				getH = &path.Handler
			case METHOD_POST:
				allMiddlewares := append(path.Info.Middlewares, postRequest)
				postM = allMiddlewares
				postH = &path.Handler
			case METHOD_PUT:
				allMiddlewares := append(path.Info.Middlewares, putRequest)
				putM = allMiddlewares
				putH = &path.Handler
			case METHOD_DELETE:
				allMiddlewares := append(path.Info.Middlewares, deleteRequest)
				deleteM = allMiddlewares
				deleteH = &path.Handler
			default:
				panic("method not implemented .")
			}
		}
		s.mux.Handle(route, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var chain []func(http.Handler) http.Handler
			var h func(http.ResponseWriter, *http.Request)

			switch r.Method {
			case http.MethodGet:
				chain = getM
				h = *getH
			case http.MethodPost:
				chain = postM
				h = *postH
			case http.MethodPut:
				chain = putM
				h = *putH
			case http.MethodDelete:
				chain = deleteM
				h = *deleteH
			default:
				panic("method not implemented")
			}

			finalHandler := chainMiddleware(http.HandlerFunc(h), chain...)
			finalHandler.ServeHTTP(w, r)
		}))
	}

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
