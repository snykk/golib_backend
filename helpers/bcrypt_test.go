package helpers_test

import (
	"testing"

	"github.com/snykk/golib_backend/helpers"
	"github.com/stretchr/testify/assert"
)

func TestGenerateHash(t *testing.T) {
	pass := "111111"
	hash, err := helpers.GenerateHash(pass)
	if err != nil {
		t.Error(err)
	}

	assert.NotNil(t, hash)
	assert.NotEqual(t, pass, hash)
	assert.Equal(t, 60, len(hash))
}
