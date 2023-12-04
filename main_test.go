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

func TestOutput(t *testing.T) {
	cases := []struct {
		Description string
		options     programOptions
		input       string
		fileName    string
		Want        string
	}{
		{"Empty", programOptions{numberOfCharacters: true, numberOfLines: true, numberOfWords: true, numberOfBytes: true}, "", "fileEmpty", "0\t0\t0\t0\tfileEmpty"},
		{"Default", programOptions{numberOfCharacters: false, numberOfLines: true, numberOfWords: true, numberOfBytes: true}, "s", "singleCharFile", "0\t1\t1\tsingleCharFile"},
		{"Single argument", programOptions{numberOfCharacters: true, numberOfLines: false, numberOfWords: false, numberOfBytes: false}, "s", "singleCharFile", "1\tsingleCharFile"},
	}

	for _, test := range cases {
		t.Run(test.Description, func(t *testing.T) {
			reader := bufio.NewReader(strings.NewReader(test.input))
			result := fileResults{fileName: &test.fileName}
			calculate(reader, &result, &test.options)
			out := formatOutput(&result, &test.options)

			if out != test.Want {
				t.Errorf("got %v, want %v", out, test.Want)
			}
		})
	}
}
