package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/loissascha/go-http-server/respond"
	"github.com/loissascha/go-http-server/server"
)

// Example Implementation

type loginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResult struct {
	Method  string `json:"method"`
	Success bool   `json:"success"`
	Jwt     string `json:"jwt"`
}

var s *server.Server

func main() {
	se, err := server.NewServer(
		server.EnableTranslations(),
		server.EnableAutoDetectLanguage(),
		server.AddTranslationFile("en", "en.json"),
		server.AddTranslationFile("de", "de.json"),
		server.SetDefaultLanguage("de"),
		server.SetExportTypesLocation("./out/types.ts"),
	)
	if err != nil {
		panic(err)
	}
	s = se

	s.GET("/", homeHandler)
	s.GET("/test", homeHandler)

	s.GET("/login", loginGet)
	s.POST("/login", loginPost, server.WithExportType[loginInput](), server.WithExportType[loginResult]())

	fmt.Println("server:", s)

	err = s.Serve(":4422")
	if err != nil {
		panic(err)
	}
}

func loginPost(w http.ResponseWriter, r *http.Request) {
	res := loginResult{
		Method:  "POST",
		Success: true,
		Jwt:     "TEST",
	}
	respond.JSON(w, http.StatusOK, res)
}

func loginGet(w http.ResponseWriter, r *http.Request) {
	respond.JSON(w, http.StatusOK, map[string]string{
		"method": "GET",
	})
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	c, err := json.Marshal(s.GetTMap(r))
	if err != nil {
		panic(err)
	}
	respond.JSON(w, http.StatusOK, map[string]string{
		"status":      "ok",
		"test_str":    s.GetTString(r, "test_str"),
		"unknown_key": s.GetTString(r, "unknown_key"),
		"map":         string(c),
	})
}
