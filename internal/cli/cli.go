package cli

import "fmt"

func AskYesOrNo(question string) bool {
	fmt.Println(question + " (y/n)")

	var answer string

	for answer != "y" && answer != "n" {
		fmt.Scanln(&answer)

		if answer == "y" {
			return true
		}

		if answer == "n" {
			return false
		}

		fmt.Println("Please answer with y or n")
	}

	return false
}

func Ask(question string) string {
	fmt.Println(question)

	var answer string
	fmt.Scanln(&answer)

	return answer
}

func isLowercaseAlphaNumericString(s string) bool {
	for _, c := range s {
		if !isLowercaseAlphaNumericRune(c) {
			return false
		}
	}

	return true
}

func isLowercaseAlphaNumericRune(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9')
}
