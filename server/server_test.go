package server

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServerSetup(t *testing.T) {
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

	s, err = NewServer(EnableTranslations())
	assert.Equal(t, 1, len(s.Options))
	assert.True(t, s.TranslationsEnabled)
	assert.Nil(t, err)
}

func testRoute(w http.ResponseWriter, r *http.Request) {
}
