package articles

import "unicode"

func ValidateSlug(slug string) bool {
	if len(slug) == 0 {
		return false
	}
	for _, char := range slug {
		if !unicode.IsLetter(char) && !unicode.IsDigit(char) && char != '-' {
			return false
		}
	}
	return true
}
