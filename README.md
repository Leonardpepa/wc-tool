# wc tool clone

## Purpose
This project is a solution for [Build your own wc tool](https://codingchallenges.fyi/challenges/challenge-wc)
build for my personal educational purposes

## Description
wc is a command line tool, read the [original specification](https://www.gnu.org/software/coreutils/manual/html_node/wc-invocation.html#wc-invocation) for more

## Features
* Implemented
```terminal 
-l, --lines Number of lines
-w, --words Number of words
-c, --bytes Nymber of bytes
-m, --chars Number of characters
-h, --help  Usage
  ```
* Not implemented
```terminal
--files0-from=file
--total=when
--max-line-length
```
## Usage
```temrinal
Usage: ./wc-tool [OPTIONS]... [FILE]...
If no [OPTIONS] specified then -l, -w, -c = true
If no [FILE] specified then input = stdin
[OPTIONS]:
  -l, --lines Number of lines
  -w, --words Number of words
  -c, --bytes Nymber of bytes
  -m, --chars Number of characters
  -h, --help Usage
```
## How to run
1. Clone the repo ```git clone https://github.com/Leonardpepa/wc-tool```
2. Build ```go build```
3. run on windows```wc-tool.exe [OPTIONS] [FILE]```
4. run on linux ```./wc-tool [OPTIONS] [FILE]```
5. run tests ```go test```
