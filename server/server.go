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

func (s *Server) GetLanguages() []string {
	res := []string{}
	for l := range s.Languages {
		res = append(res, l)
	}
	return res
}

func (s *Server) GetActiveLanguage(r *http.Request) string {
	if !s.TranslationsEnabled {
		panic("Translations not enabled!")
	}
	url := strings.TrimSpace(r.URL.Path)
	url = strings.TrimLeft(url, "/")
	urlSplit := strings.Split(url, "/")
	activeLanguage := urlSplit[0]
	return activeLanguage
}

func (s *Server) GetTMap(r *http.Request) map[string]string {
	if !s.TranslationsEnabled {
		panic("Translations not enabled!")
	}
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

func (s *Server) GetTString(r *http.Request, key string) string {
	if !s.TranslationsEnabled {
		panic("Translations not enabled!")
	}
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
		getChain := []func(http.Handler) http.Handler{}
		var getHandler *func(http.ResponseWriter, *http.Request)
		postChain := []func(http.Handler) http.Handler{}
		var postHandler *func(http.ResponseWriter, *http.Request)
		putChain := []func(http.Handler) http.Handler{}
		var putHandler *func(http.ResponseWriter, *http.Request)
		deleteChain := []func(http.Handler) http.Handler{}
		var deleteHandler *func(http.ResponseWriter, *http.Request)
		for _, path := range serverPaths {
			switch path.Method {
			case METHOD_GET:
				if getHandler != nil {
					panic("double GET route!")
				}
				allMiddlewares := append(path.Info.Middlewares, getRequest)
				getChain = allMiddlewares
				getHandler = &path.Handler
			case METHOD_POST:
				if postHandler != nil {
					panic("double POST route!")
				}
				allMiddlewares := append(path.Info.Middlewares, postRequest)
				postChain = allMiddlewares
				postHandler = &path.Handler
			case METHOD_PUT:
				if putHandler != nil {
					panic("double PUT route!")
				}
				allMiddlewares := append(path.Info.Middlewares, putRequest)
				putChain = allMiddlewares
				putHandler = &path.Handler
			case METHOD_DELETE:
				if deleteHandler != nil {
					panic("double DELETE route!")
				}
				allMiddlewares := append(path.Info.Middlewares, deleteRequest)
				deleteChain = allMiddlewares
				deleteHandler = &path.Handler
			default:
				panic("method not implemented .")
			}
		}
		s.mux.Handle(route, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var chain []func(http.Handler) http.Handler
			var h *func(http.ResponseWriter, *http.Request)

			switch r.Method {
			case http.MethodGet:
				chain = getChain
				h = getHandler
			case http.MethodPost:
				chain = postChain
				h = postHandler
			case http.MethodPut:
				chain = putChain
				h = putHandler
			case http.MethodDelete:
				chain = deleteChain
				h = deleteHandler
			default:
				panic("method not implemented")
			}

			if h == nil {
				http.Error(w, "Not found", http.StatusNotFound)
				return
			}

			finalHandler := chainMiddleware(http.HandlerFunc(*h), chain...)
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
