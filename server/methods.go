package server

import "net/http"

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
