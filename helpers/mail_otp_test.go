package helpers_test

import (
	"testing"

	"github.com/snykk/golib_backend/helpers"
)

func TestGenerateCode(t *testing.T) {
	length := 6

	code1, err := helpers.GenerateCode(length)
	if err != nil {
		t.Error("error occurred while generating code 1")
	}

	code2, err := helpers.GenerateCode(length)
	if err != nil {
		t.Error("error occurred while generating code 2")
	}

	if code1 == code2 {
		t.Error("function have to generate difference code when execute twice ")
	}

	if len(code1) != length {
		t.Error("invalid length")
	}
}
