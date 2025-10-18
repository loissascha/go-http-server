package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServerSetup(t *testing.T) {
	s, err := NewServer(EnableTranslations())
	assert.Equal(t, 1, len(s.Options))
	assert.Nil(t, err)
}
