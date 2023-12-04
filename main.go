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
	options            programOptions
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

	if numberOfFiles > 0 {
		handleFiles(options, numberOfFiles > 1)
		return
	}

	handleStdin(options)
}

func handleStdin(options programOptions) {
	reader := bufio.NewReader(os.Stdin)
	res := fileResults{options: options}

	calculate(reader, &res)
	printResults(&res)
}

func handleFiles(options programOptions, multipleFiles bool) {
	result := fileResults{options: options}
	var storeTotalResults *fileResults
	var reader *bufio.Reader

	if multipleFiles {
		total := "total"
		storeTotalResults = &fileResults{options: options, fileName: &total}
	}

	for _, fileName := range options.fileNames {
		var file *os.File
		var err error

		// recalculate only when file changes
		// first time file is nil
		if file == nil || *result.fileName != fileName {
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
			calculate(reader, &result)
		}

		printResults(&result)

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
		printResults(storeTotalResults)
	}
}

func printResults(results *fileResults) {
	output := " "

	// -l
	if results.options.numberOfLines {
		output += fmt.Sprintf("%v  ", results.numberOfLines)
	}

	// -w
	if results.options.numberOfWords {
		output += fmt.Sprintf("%v  ", results.numberOfWords)
	}

	// -c
	if results.options.numberOfBytes {
		output += fmt.Sprintf("%v  ", results.numberOfBytes)
	}

	// -m
	if results.options.numberOfCharacters {
		output += fmt.Sprintf("%v  ", results.numberOfCharacters)
	}

	if results.fileName != nil && *results.fileName != "-" {
		output += fmt.Sprint(*results.fileName)
	}

	fmt.Println(output)
}

func calculate(fileReader *bufio.Reader, results *fileResults) {

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

		if results.options.numberOfBytes {
			results.numberOfBytes += uint64(runeSize)
		}

		if results.options.numberOfLines && runeRead == '\n' {
			results.numberOfLines++
		}

		if results.options.numberOfCharacters {
			results.numberOfCharacters++
		}

		if results.options.numberOfWords {
			if unicode.IsSpace(runeRead) && unicode.IsSpace(prevRune) == false {
				results.numberOfWords++
			}
		}
		prevRune = runeRead
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
