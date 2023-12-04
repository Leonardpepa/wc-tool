package main

import (
	"bufio"
	"strings"
	"testing"
)

func TestCalculate(t *testing.T) {
	cases := []struct {
		description string
		input       string
		expected    Result
	}{
		{"Empty", "", Result{numberOfBytes: 0, numberOfWords: 0, numberOfLines: 0, numberOfCharacters: 0}},
		{"Single char", "s", Result{numberOfBytes: 1, numberOfWords: 1, numberOfLines: 0, numberOfCharacters: 1}},
		{"Multibyte chars", "s⌘ f", Result{numberOfBytes: 6, numberOfWords: 2, numberOfLines: 0, numberOfCharacters: 4}},
		{"Trailing newline", "this is a sentence\n\nacross multiple\nlines\n", Result{numberOfBytes: 42, numberOfWords: 7, numberOfLines: 4, numberOfCharacters: 42}},
		{"No trailing newline", "this is a sentence\n\nacross multiple\nlines", Result{numberOfBytes: 41, numberOfWords: 7, numberOfLines: 3, numberOfCharacters: 41}},
	}

	for _, test := range cases {
		t.Run(test.description, func(t *testing.T) {
			options := programOptions{
				numberOfCharacters: true,
				numberOfLines:      true,
				numberOfWords:      true,
				numberOfBytes:      true,
			}

			reader := bufio.NewReader(strings.NewReader(test.input))
			result := Result{}
			calculate(reader, &result, &options)
			if result != test.expected {
				t.Errorf("result %v, expected %v", result, test.expected)
			}

		})
	}
}

func TestSingleFile(t *testing.T) {

	t.Run("Calculate single file", func(t *testing.T) {
		fileNames := []string{"./tests/test.txt"}
		expected := Result{
			fileName:           "total",
			numberOfBytes:      342190,
			numberOfLines:      7145,
			numberOfWords:      58164,
			numberOfCharacters: 339292,
		}

		total, _ := processMultipleFiles(&programOptions{
			fileNames:          fileNames,
			numberOfBytes:      true,
			numberOfWords:      true,
			numberOfLines:      true,
			numberOfCharacters: true,
		})

		if total != expected {
			t.Errorf("result %v, expected %v", total, expected)
		}

	})

}

func TestCalculateTotals(t *testing.T) {
	t.Run("Calculate multiple files", func(t *testing.T) {
		fileNames := []string{"./tests/test.txt", "./tests/test.txt"}
		expected := Result{
			fileName:           "total",
			numberOfBytes:      684380,
			numberOfLines:      14290,
			numberOfWords:      116328,
			numberOfCharacters: 678584,
		}

		total, _ := processMultipleFiles(&programOptions{
			fileNames:          fileNames,
			numberOfBytes:      true,
			numberOfWords:      true,
			numberOfLines:      true,
			numberOfCharacters: true,
		})

		if total != expected {
			t.Errorf("result %v, expected %v", total, expected)
		}

	})

}

func TestOutput(t *testing.T) {
	cases := []struct {
		description string
		options     programOptions
		input       string
		fileName    string
		expected    string
	}{
		{"Empty file", programOptions{numberOfCharacters: true, numberOfLines: true, numberOfWords: true, numberOfBytes: true}, "", "fileEmpty", "0\t0\t0\t0\tfileEmpty"},
		{"Default", programOptions{numberOfCharacters: false, numberOfLines: true, numberOfWords: true, numberOfBytes: true}, "s", "singleCharFile", "0\t1\t1\tsingleCharFile"},
		{"Single argument", programOptions{numberOfCharacters: true, numberOfLines: false, numberOfWords: false, numberOfBytes: false}, "s", "singleCharFile", "1\tsingleCharFile"},
	}

	for _, test := range cases {
		t.Run(test.description, func(t *testing.T) {
			reader := bufio.NewReader(strings.NewReader(test.input))
			result := Result{fileName: test.fileName}
			calculate(reader, &result, &test.options)
			out := formatOutput(&result, &test.options)

			if out != test.expected {
				t.Errorf("result %v, expected %v", out, test.expected)
			}
		})
	}
}
