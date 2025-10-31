package clinput

// Making a small dependancy is a bad idea, but I wanted to test out the vendor

import "strings"

func Clean_Input(input string) []string {
	return strings.Fields(strings.Trim(strings.ToLower(input), " "))
}
