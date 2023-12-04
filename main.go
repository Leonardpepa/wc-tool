package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"unicode"
)

type programOptions struct {
	fileNames          []string
	numberOfBytes      bool
	numberOfLines      bool
	numberOfWords      bool
	numberOfCharacters bool
}

type fileResults struct {
	fileName           *string
	numberOfBytes      uint64
	numberOfLines      uint64
	numberOfWords      uint64
	numberOfCharacters uint64
}

var ProgramName = filepath.Base(os.Args[0])

func main() {

	arguments := os.Args[1:]

	options, err := parseArguments(arguments)

	if err != nil {
		log.Fatal(err.Error())
	}

	numberOfFiles := len(options.fileNames)

	if numberOfFiles == 0 {
		handleStdin(&options)
	}

	handleFiles(&options, numberOfFiles > 1)
}

func handleStdin(options *programOptions) {
	reader := bufio.NewReader(os.Stdin)
	res := fileResults{}

	calculate(reader, &res, options)
	printResults(&res, options)
}

func handleFiles(options *programOptions, multipleFiles bool) {
	result := fileResults{}
	var storeTotalResults *fileResults
	var reader *bufio.Reader

	if multipleFiles {
		total := "total"
		storeTotalResults = &fileResults{fileName: &total}
	}

	for index, fileName := range options.fileNames {
		var file *os.File
		var err error

		// recalculate only when file changes
		// first time file is nil
		if index == 0 || *result.fileName != fileName {
			file, err = os.Open(fileName)
			if err != nil {
				log.Println(err.Error())
				continue
			}

			// create object once
			if reader == nil {
				reader = bufio.NewReader(file)
			} else {
				reader.Reset(file)
			}
			result.fileName = &fileName
			calculate(reader, &result, options)
		}

		printResults(&result, options)

		if multipleFiles {
			storeTotalResults.numberOfLines += result.numberOfLines
			storeTotalResults.numberOfWords += result.numberOfWords
			storeTotalResults.numberOfBytes += result.numberOfBytes
			storeTotalResults.numberOfCharacters += result.numberOfCharacters
		}
		// close the file if its open
		if file != nil {
			func(file *os.File) {
				err := file.Close()
				if err != nil {
					log.Println(err.Error())
				}
			}(file)
		}
	}

	if multipleFiles {
		printResults(storeTotalResults, options)
	}
}

func printResults(results *fileResults, options *programOptions) {
	output := " "

	// -l
	if options.numberOfLines {
		output += fmt.Sprintf("%v  ", results.numberOfLines)
	}

	// -w
	if options.numberOfWords {
		output += fmt.Sprintf("%v  ", results.numberOfWords)
	}

	// -c
	if options.numberOfBytes {
		output += fmt.Sprintf("%v  ", results.numberOfBytes)
	}

	// -m
	if options.numberOfCharacters {
		output += fmt.Sprintf("%v  ", results.numberOfCharacters)
	}

	if results.fileName != nil && *results.fileName != "-" {
		output += fmt.Sprint(*results.fileName)
	}

	fmt.Println(output)
}

func calculate(fileReader *bufio.Reader, results *fileResults, options *programOptions) {

	results.numberOfLines = 0
	results.numberOfWords = 0
	results.numberOfBytes = 0
	results.numberOfCharacters = 0

	var prevRune rune

	for {
		runeRead, runeSize, err := fileReader.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err.Error())
		}

		if options.numberOfBytes {
			results.numberOfBytes += uint64(runeSize)
		}

		if options.numberOfLines && runeRead == '\n' {
			results.numberOfLines++
		}

		if options.numberOfCharacters {
			results.numberOfCharacters++
		}

		if options.numberOfWords {
			if unicode.IsSpace(runeRead) && unicode.IsSpace(prevRune) == false {
				results.numberOfWords++
			}
		}
		prevRune = runeRead
	}
	if prevRune != rune(0) && unicode.IsSpace(prevRune) == false {
		results.numberOfWords++
	}
}

// custom parsing function
// you can do this better with the lib flag
func parseArguments(arguments []string) (programOptions, error) {
	var options programOptions

	for _, value := range arguments {
		switch value {
		case "-h", "--help":
			return programOptions{}, fmt.Errorf(usageMessage(ProgramName))
		case "-c", "--bytes":
			options.numberOfBytes = true
		case "-l", "--lines":
			options.numberOfLines = true
		case "-w", "--words":
			options.numberOfWords = true
		case "-m":
			options.numberOfCharacters = true
		default:
			// wrongs argument given
			if value[0] == '-' {
				return programOptions{}, fmt.Errorf(wrongArgumentMessage(value, ProgramName))
			}

			options.fileNames = append(options.fileNames, value)
		}
	}

	if !options.numberOfBytes && !options.numberOfLines && !options.numberOfWords && !options.numberOfCharacters {
		options.numberOfBytes, options.numberOfLines, options.numberOfWords = true, true, true
	}

	return options, nil
}

func usageMessage(programName string) string {
	return fmt.Sprintf(`
Usage: %s [OPTIONS]... [FILE]...
If no [OPTIONS] specified then -l, -w, -c = true
If no [FILE] specified then input = stdin
[OPTIONS]:
  -l, --lines Number of lines
  -w, --words Number of words
  -c, --bytes Nymber of bytes
  -m, --chars Number of characters`, programName)
}

func wrongArgumentMessage(argument string, programName string) string {
	return fmt.Sprintf("unknown option -- %s\nTry '%s --help' for more information.", argument, programName)
}
