package server

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServerSetup(t *testing.T) {
	// test simple server (without any translation and server settings)
	s, err := NewServer()
	s.GET("/test", testRoute)
	s.POST("/test", testRoute)
	s.PUT("/test", testRoute)
	s.DELETE("/test", testRoute)
	assert.Equal(t, 0, len(s.Options))
	assert.False(t, s.TranslationsEnabled)
	assert.False(t, s.AutoDetectLanguageEnabled)
	assert.Equal(t, s.DefaultLanguage, "")
	assert.Equal(t, 1, len(s.Paths)) // has only one path because "/test" always refers to the same path
	assert.Nil(t, err)

	// test if language settings work
	s, err = NewServer(
		EnableTranslations(),
		EnableAutoDetectLanguage(),
		SetDefaultLanguage("en"),
		AddTranslationFile("en", "en_test.json"),
		AddTranslationFile("de", "de_test.json"),
	)
	s.GET("/test", testRoute)
	s.POST("/test", testRoute)
	s.PUT("/test", testRoute)
	s.DELETE("/test", testRoute)
	s.POSTI("/test/no/langs", testRoute)
	assert.Equal(t, 5, len(s.Options))
	assert.True(t, s.TranslationsEnabled)
	assert.True(t, s.AutoDetectLanguageEnabled)
	assert.Equal(t, s.DefaultLanguage, "en")
	assert.Equal(t, 4, len(s.Paths)) // 3 routes for "/test" ("/test", "/en/test", "/de/test") and one for "/test/no/langs"
	assert.Contains(t, s.Paths, "/test")
	assert.Contains(t, s.Paths, "/de/test")
	assert.Contains(t, s.Paths, "/en/test")
	assert.NotContains(t, s.Paths, "fr/test")
	assert.Nil(t, err)
}

func testRoute(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Test was successful", http.StatusNotAcceptable)
}
