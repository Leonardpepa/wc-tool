package main

import (
	"bufio"
	"reflect"
	"strings"
	"testing"
)

func TestCalculate(t *testing.T) {
	cases := []struct {
		Description string
		input       string
		Want        fileResults
	}{
		{"Empty", "", fileResults{numberOfBytes: 0, numberOfWords: 0, numberOfLines: 0, numberOfCharacters: 0}},
		{"Single char", "s", fileResults{numberOfBytes: 1, numberOfWords: 1, numberOfLines: 0, numberOfCharacters: 1}},
		{"Multibyte chars", "s⌘ f", fileResults{numberOfBytes: 6, numberOfWords: 2, numberOfLines: 0, numberOfCharacters: 4}},
		{"Trailing newline", "this is a sentence\n\nacross multiple\nlines\n", fileResults{numberOfBytes: 42, numberOfWords: 7, numberOfLines: 4, numberOfCharacters: 42}},
		{"No trailing newline", "this is a sentence\n\nacross multiple\nlines", fileResults{numberOfBytes: 41, numberOfWords: 7, numberOfLines: 3, numberOfCharacters: 41}},
	}

	for _, test := range cases {
		t.Run(test.Description, func(t *testing.T) {
			options := programOptions{
				numberOfCharacters: true,
				numberOfLines:      true,
				numberOfWords:      true,
				numberOfBytes:      true,
			}

			reader := bufio.NewReader(strings.NewReader(test.input))
			result := fileResults{}
			calculate(reader, &result, &options)
			if reflect.DeepEqual(result, test.Want) == false {
				t.Errorf("got %v, want %v", result, test.Want)
			}

		})
	}
}
