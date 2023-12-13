package util

func IsSlug(s string) bool {
	for _, c := range s {
		if isUppercase(c) {
			return false
		}
		if !isLetter(c) && !isNumber(c) && c != '-' {
			return false
		}
	}
	return true
}

func isLetter(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func isNumber(c rune) bool {
	return (c >= '0' && c <= '9')
}

func isUppercase(c rune) bool {
	return (c >= 'A' && c <= 'Z')
}
