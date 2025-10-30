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

// Get all languages
func (s *Server) GetLanguages() []string {
	res := []string{}
	for l := range s.Languages {
		res = append(res, l)
	}
	return res
}

// Get currently active language
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

// Get all translations for the currently active language
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

// Get translated value for key
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

// get the http.ServeMux in case you need to access it. Preferably s.Method should be used..
func (s *Server) GetMux() *http.ServeMux {
	return s.mux
}

type chainHandler struct {
	chain   []func(http.Handler) http.Handler
	handler *func(http.ResponseWriter, *http.Request)
}

type routeMethodHandlers struct {
	get    chainHandler
	post   chainHandler
	put    chainHandler
	delete chainHandler
}

// create chain and handler for every method in this server paths array (per route)
func (s *Server) getRouteServerPathsChainAndHandler(serverPaths []ServerPath) routeMethodHandlers {
	result := routeMethodHandlers{}
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
	result.get = chainHandler{
		chain:   getChain,
		handler: getHandler,
	}
	result.post = chainHandler{
		chain:   postChain,
		handler: postHandler,
	}
	result.put = chainHandler{
		chain:   putChain,
		handler: putHandler,
	}
	result.delete = chainHandler{
		chain:   deleteChain,
		handler: deleteHandler,
	}
	return result
}

func (s *Server) setupHandlers() error {
	for route, serverPaths := range s.Paths {
		data := s.getRouteServerPathsChainAndHandler(serverPaths)
		s.mux.Handle(route, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var chain []func(http.Handler) http.Handler
			var h *func(http.ResponseWriter, *http.Request)

			switch r.Method {
			case http.MethodGet:
				chain = data.get.chain
				h = data.get.handler
			case http.MethodPost:
				chain = data.post.chain
				h = data.post.handler
			case http.MethodPut:
				chain = data.put.chain
				h = data.put.handler
			case http.MethodDelete:
				chain = data.delete.chain
				h = data.delete.handler
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
	return nil
}

// Start the server on the given addr (or port)
func (s *Server) Serve(addr string) error {
	err := s.setupHandlers()
	if err != nil {
		return err
	}

	err = http.ListenAndServe(addr, s.mux)
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
