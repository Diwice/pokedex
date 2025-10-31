package main

import (
	"testing"
	"dep/clinput"
)

func Test_clean_input(t *testing.T) {
	cases := []struct{
		input string
		expected []string
	}{
		{
			input: "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input: "DiFfErEnT sTrInG tHiS tImE   ",
			expected: []string{"different", "string", "this", "time"},
		},
		{
			input: "Test suite #3. Shouldn't fail",
			expected: []string{"test", "suite", "#3.", "shouldn't", "fail"},
		},
	}

	for _, c := range cases {
		actual := clinput.Clean_Input(c.input)

		if len(actual) != len(c.expected) {
			t.Errorf("Expected length doesn't match the actual length : %d (%d expected)", len(actual), len(c.expected))
		}

		for i := range actual {
			word := actual[i]
			expected_word := c.expected[i]

			if word != expected_word {
				t.Errorf("Expected : %s ; Got : %s ; index : %d", expected_word, word, i)
			}
		}
	}
}
