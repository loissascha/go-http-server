package server

import (
	"fmt"
	"net/http"
	"strings"
)

func getRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed!", http.StatusMethodNotAllowed)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func postRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed!", http.StatusMethodNotAllowed)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func putRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "Method Not Allowed!", http.StatusMethodNotAllowed)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func deleteRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Method Not Allowed!", http.StatusMethodNotAllowed)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Server) redirectToTranslatedUrl(w http.ResponseWriter, r *http.Request) {
	if s.AutoDetectLanguageEnabled {
		langHeader := r.Header.Get("Accept-Language")
		fmt.Println("langHeader:", langHeader)

		languages := strings.SplitSeq(langHeader, ",")
		for lang := range languages {
			lang := strings.TrimSpace(strings.Split(lang, ";")[0])
			lang = strings.ToLower(lang)

			for short := range s.Languages {
				if strings.HasPrefix(lang, short) {
					http.Redirect(w, r, fmt.Sprintf("/%s%s", short, r.URL.Path), http.StatusFound)
					return
				}
			}
		}
	}
	http.Redirect(w, r, fmt.Sprintf("/%s%s", s.DefaultLanguage, r.URL.Path), http.StatusFound)
}
