package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/neurosnap/sentences/english"
)

// VERSION is the semantic version number
var VERSION string

// COMMITHASH is the git commit hash value
var COMMITHASH string

func run(fname string, delim string, debug bool) {
	if debug {
		fmt.Printf("file [%s], delim [%s]\n", fname, delim)
	}

	var text []byte
	var err error

	if fname != "" {
		text, err = ioutil.ReadFile(fname)
		if err != nil {
			panic(err)
		}
	} else {
		stat, err := os.Stdin.Stat()
		if err != nil {
			panic(err)
		}

		if (stat.Mode() & os.ModeCharDevice) != 0 {
			return
		}

		reader := bufio.NewReader(os.Stdin)
		text, err = ioutil.ReadAll(reader)
		if err != nil {
			panic(err)
		}
	}

	tokenizer, err := english.NewSentenceTokenizer(nil)
	if err != nil {
		panic(err)
	}

	sentences := tokenizer.Tokenize(string(text))

	if debug {
		for _, s := range sentences {
			fmt.Println(s)
		}
		fmt.Println("---")
	}

	for _, s := range sentences {
		text := strings.Join(strings.Fields(s.Text), " ")

		text = strings.Join([]string{text, delim}, "")
		fmt.Printf("%s", text)
	}
}

func main() {
	var ver bool
	verStr := "Get current version of sentences"
	flag.BoolVar(&ver, "version", false, verStr)
	flag.BoolVar(&ver, "v", false, fmt.Sprintf("%s (alias of --version)", verStr))

	var fname string
	fileStr := "Read file as source input instead of stdin"
	flag.StringVar(&fname, "file", "", fileStr)
	flag.StringVar(&fname, "f", "", fmt.Sprintf("%s (alias of --file)", fileStr))

	var delim string
	delimStr := "Delimiter used to demarcate sentence boundaries"
	flag.StringVar(&delim, "delimiter", "\n", delimStr)
	flag.StringVar(&delim, "d", "\n", fmt.Sprintf("%s (alias of --delimiter)", delimStr))

	var debug bool
	debugStr := "Debug mode"
	flag.BoolVar(&debug, "debug", false, debugStr)

	flag.Parse()

	if ver {
		fmt.Println(VERSION)
		fmt.Println(COMMITHASH)
		return
	}

	run(fname, delim, debug)
}
