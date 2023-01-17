package helpers

import (
	"errors"
	"fmt"

	"github.com/snykk/golib_backend/constants"
)

func IsGenderValid(gender string) error {
	if !isArrayContains(constants.ListGender, gender) {
		var option string
		for index, g := range constants.ListGender {
			option += g
			if index != len(constants.ListGender)-1 {
				option += ", "
			}
		}

		return fmt.Errorf("gender must be one of [%s]", option)
	}

	return nil
}

func isArrayContains(arr []string, str string) bool {
	for _, item := range arr {
		if item == str {
			return true
		}
	}
	return false
}

func IsRatingValid(rating int) error {
	if rating < 1 || rating > 10 {
		return errors.New("the rating must be in the range 1 - 10")
	}
	return nil
}
