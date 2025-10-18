package main

import (
	"fmt"

	"github.com/loissascha/go-http-server/server"
)

// Example Implementation

func main() {
	s := server.NewServer(
		server.EnableTranslations(),
		server.AddTranslationFile("en", "en_test.json"),
		server.AddTranslationFile("de", "de_test.json"),
		server.SetDefaultLanguage("en"),
	)

	fmt.Println("server:", s)

	err := s.Serve(":44444")
	if err != nil {
		panic(err)
	}
}
