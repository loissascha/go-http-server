package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/loissascha/go-http-server/respond"
	"github.com/loissascha/go-http-server/server"
)

// Example Implementation

var s *server.Server

func main() {
	server, err := server.NewServer(
		server.EnableTranslations(),
		server.EnableAutoDetectLanguage(),
		server.AddTranslationFile("en", "en_test.json"),
		server.AddTranslationFile("de", "de_test.json"),
		server.SetDefaultLanguage("de"),
	)
	if err != nil {
		panic(err)
	}
	s = server

	s.GET("/", homeHandler)
	s.GET("/test", homeHandler)

	fmt.Println("server:", s)

	err = s.Serve(":44444")
	if err != nil {
		panic(err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	c, err := json.Marshal(s.GetLanguageStringMap(r))
	if err != nil {
		panic(err)
	}
	respond.JSON(w, http.StatusOK, map[string]string{
		"status":      "ok",
		"test_str":    s.GetLanguageString(r, "test_str"),
		"unknown_key": s.GetLanguageString(r, "unknown_key"),
		"map":         string(c),
	})
}
