package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServerSetup(t *testing.T) {
	s := NewServer(EnableTranslations())
	assert.Equal(t, 1, len(s.Options))
}
