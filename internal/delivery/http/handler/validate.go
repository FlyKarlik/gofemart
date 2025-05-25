package handler

import (
	"strconv"
	"unicode"
)

func isValidOrderNumber(number string) bool {
	var sum int
	var alt bool

	for i := len(number) - 1; i >= 0; i-- {
		r := rune(number[i])

		if !unicode.IsDigit(r) {
			return false
		}

		n, _ := strconv.Atoi(string(r))
		if alt {
			n *= 2
			if n > 9 {
				n -= 9
			}
		}
		sum += n
		alt = !alt
	}
	return sum%10 == 0
}
