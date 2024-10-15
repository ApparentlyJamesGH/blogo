package utils

import (
	"math"
	"strings"
)

func ReadTime(s string) int {
	words := strings.Fields(s)
	wordCount := len(words)

	readSpeed := 200

	// Calculate time in minutes, use math.Ceil to round up to nearest whole number.
	readTime := math.Ceil(float64(wordCount) / float64(readSpeed))
	return int(readTime)
}
