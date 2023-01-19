package helpers_test

import (
	"testing"

	"github.com/snykk/golib_backend/helpers"
	"github.com/stretchr/testify/assert"
)

func TestIsGenderValid(t *testing.T) {
	t.Run("When Success", func(t *testing.T) {
		t.Run("Test 1", func(t *testing.T) {
			err := helpers.IsGenderValid("male")

			assert.Nil(t, err)
		})
		t.Run("Test 2", func(t *testing.T) {
			err := helpers.IsGenderValid("female")

			assert.Nil(t, err)
		})
	})
	t.Run("When Failure", func(t *testing.T) {
		t.Run("Test 1", func(t *testing.T) {
			err := helpers.IsGenderValid("malee")

			assert.NotNil(t, err)
		})
		t.Run("Test 2", func(t *testing.T) {
			err := helpers.IsGenderValid("")

			assert.NotNil(t, err)
		})
	})
}

func TestIsRatingValid(t *testing.T) {
	t.Run("When Success", func(t *testing.T) {
		t.Run("Test 1", func(t *testing.T) {
			err := helpers.IsRatingValid(9)

			assert.Nil(t, err)
		})
		t.Run("Test 2", func(t *testing.T) {
			err := helpers.IsRatingValid(3)

			assert.Nil(t, err)
		})
		t.Run("Test 3", func(t *testing.T) {
			err := helpers.IsRatingValid(1)

			assert.Nil(t, err)
		})
		t.Run("Test 4", func(t *testing.T) {
			err := helpers.IsRatingValid(0)

			assert.Nil(t, err)
		})
	})
	t.Run("When Failure", func(t *testing.T) {
		t.Run("Test 1", func(t *testing.T) {
			err := helpers.IsRatingValid(11)

			assert.NotNil(t, err)
		})
		t.Run("Test 2", func(t *testing.T) {
			err := helpers.IsRatingValid(-2)

			assert.NotNil(t, err)
		})
	})
}
