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
	fileName           string
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

	var storeTotalResults fileResults
	total := len(options.fileNames) > 1

	if total {
		storeTotalResults = fileResults{options: options, fileName: "total"}
	}

	res := fileResults{options: options}

	var reader *bufio.Reader

	for _, fileName := range options.fileNames {
		file, err := os.Open(fileName)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		// recalculate only when file changes
		if res.fileName != fileName {
			reader = bufio.NewReader(file)
			res.fileName = fileName
			calculate(reader, &res)
		}

		printResults(&res)

		if total {
			storeTotalResults.numberOfLines += res.numberOfLines
			storeTotalResults.numberOfWords += res.numberOfWords
			storeTotalResults.numberOfBytes += res.numberOfBytes
			storeTotalResults.numberOfCharacters += res.numberOfCharacters
		}

		func(file *os.File) {
			err := file.Close()
			if err != nil {
				log.Fatal(err.Error())
			}
		}(file)

	}

	if total {
		printResults(&storeTotalResults)
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

	if results.fileName != "" && results.fileName != "-" {
		output += fmt.Sprint(results.fileName)
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

func parseArguments(arguments []string) (programOptions, error) {
	numOfArguments := len(arguments)

	if numOfArguments == 0 {
		return programOptions{}, fmt.Errorf(usageMessage(ProgramName))
	}

	var fileInfo programOptions

	for _, value := range arguments {
		switch value {
		case "--help":
			return programOptions{}, fmt.Errorf(usageMessage(ProgramName))
		case "-c":
			fileInfo.numberOfBytes = true
		case "-l":
			fileInfo.numberOfLines = true
		case "-w":
			fileInfo.numberOfWords = true
		case "-m":
			fileInfo.numberOfCharacters = true
		default:
			// wrongs argument given
			if value[0] == '-' {
				return programOptions{}, fmt.Errorf(wrongArgumentMessage(value, ProgramName))
			}

			fileInfo.fileNames = append(fileInfo.fileNames, value)
		}
	}

	if !fileInfo.numberOfBytes && !fileInfo.numberOfLines && !fileInfo.numberOfWords && !fileInfo.numberOfCharacters {
		fileInfo.numberOfBytes, fileInfo.numberOfLines, fileInfo.numberOfWords = true, true, true
	}

	return fileInfo, nil
}

func usageMessage(programName string) string {
	return fmt.Sprintf("\nUsage: %s [OPTIONS]... [FILE]..."+
		"\nIf no [OPTIONS] specified then -l, -w, -c = true"+
		"\nIf no [FILE] specified then input = stdin"+
		"\n[OPTIONS]:"+
		"\n-l Number of lines"+
		"\n-w Number of words"+
		"\n-c Number of bytes"+
		"\n-m Number Of characters", programName)
}

func wrongArgumentMessage(argument string, programName string) string {
	return fmt.Sprintf("unknown option -- %s\nTry '%s --help' for more information.", argument, programName)
}
