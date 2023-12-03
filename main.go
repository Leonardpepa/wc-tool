package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"unicode/utf8"
)

type FileInfo struct {
	fileName           []string
	fileBytes          []byte
	numberOfBytes      bool
	numberOfLines      bool
	numberOfWords      bool
	numberOfCharacters bool
}

const MaxNumberOfArguments = 5

var ProgramName = filepath.Base(os.Args[0])

func main() {

	arguments := os.Args[1:]

	fileInfo, err := parseArguments(arguments)

	if err != nil {
		log.Fatal(err.Error())
	}

	for index, value := range fileInfo.fileName {

		fileContentsBytes, err := readFile(value)

		if err != nil {
			log.Println(err.Error())
			continue
		}

		fileInfo.fileBytes = fileContentsBytes
		printResults(&fileInfo, index)
	}
}

func printResults(fileInfo *FileInfo, i int) {
	output := " "

	// -l
	if fileInfo.numberOfLines {
		lineCounter := getNumberOfLines(fileInfo.fileBytes)
		output += fmt.Sprintf("%d  ", lineCounter)
	}

	// -w
	if fileInfo.numberOfWords {
		wordCount := getNumberOfWords(fileInfo.fileBytes)
		output += fmt.Sprintf("%d  ", wordCount)
	}

	// -c
	if fileInfo.numberOfBytes {
		output += fmt.Sprintf("%d  ", len(fileInfo.fileBytes))
	}

	// -m
	if fileInfo.numberOfCharacters {
		output += fmt.Sprintf("%d  ", getNumberOfCharacters(fileInfo.fileBytes))
	}

	if fileInfo.fileName[i] != "" && fileInfo.fileName[i] != "-" {
		output += fmt.Sprint(fileInfo.fileName[i])
	}

	fmt.Println(output)
}

func parseArguments(arguments []string) (FileInfo, error) {
	numOfArguments := len(arguments)

	if numOfArguments > MaxNumberOfArguments {
		return FileInfo{}, fmt.Errorf(usageMessage(ProgramName))
	}

	var fileInfo FileInfo

	for _, value := range arguments {
		switch value {
		case "--help":
			return FileInfo{}, fmt.Errorf(usageMessage(ProgramName))
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
				return FileInfo{}, fmt.Errorf(wrongArgumentMessage(value, ProgramName))
			}

			fileInfo.fileName = append(fileInfo.fileName, value)
		}
	}

	if !fileInfo.numberOfBytes && !fileInfo.numberOfLines && !fileInfo.numberOfWords && !fileInfo.numberOfCharacters {
		fileInfo.numberOfBytes, fileInfo.numberOfLines, fileInfo.numberOfWords = true, true, true
	}

	return fileInfo, nil
}

func readFile(fileName string) ([]byte, error) {

	if fileName == "" || fileName == "-" {
		return io.ReadAll(os.Stdin)
	}

	return os.ReadFile(fileName)
}

func getNumberOfCharacters(fileContent []byte) int {
	return utf8.RuneCount(fileContent)
}

func getNumberOfLines(fileContent []byte) int {
	return bytes.Count(fileContent, []byte("\n"))
}

func getNumberOfWords(fileContent []byte) int {
	return len(bytes.Fields(fileContent))
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
