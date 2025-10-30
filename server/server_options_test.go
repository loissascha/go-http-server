package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadTranslationFile(t *testing.T) {
	res := readTranslationFile("en_test.json")
	assert.Equal(t, 1, len(res))
	assert.Equal(t, "Test String", res["test_str"])
}
