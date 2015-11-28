package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/neurosnap/sentences/english"
)

func main() {
	if len(os.Args) < 2 {
		panic("Please supply data file to test sentence tokenizer")
	}

	file, _ := ioutil.ReadFile(os.Args[1])
	text := string(file)

	expected_sentences := strings.Split(text, "\n")

	tokenizer, err := english.NewSentenceTokenizer(nil)
	if err != nil {
		panic(err)
	}

	actual_sentences := tokenizer.Tokenize(text)

	a_len := len(actual_sentences)
	e_len := len(expected_sentences)
	perc := (float64(a_len) / float64(e_len)) * 100

	fmt.Printf("Actual Sentences: %d, Expected Sentences: %d, Percent: %f%%\n", a_len, e_len, perc)
}
