package tool

import (
	"math/rand"
)

func Random(min, max int) int {
	return rand.Intn(max-min) + min
}

func RandString(s []string, l int) string {
	if len(s) == 0 {
		return ""
	}

	if l <= 0 {
		return s[rand.Intn(len(s))]
	}

	var result string

	for i := 0; i < l; i++ {
		result += s[rand.Intn(len(s))]
	}

	return result
}
