package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/neurosnap/sentences"
	"github.com/neurosnap/sentences/data"
	"github.com/neurosnap/sentences/english"
)

var VERSION string

func main() {
	b, err := data.Asset("data/english.json")
	if err != nil {
		panic(err)
	}

	training, err := sentences.LoadTraining(b)
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(os.Stdin)
	text, _ := ioutil.ReadAll(reader)

	tokenizer := english.NewSentenceTokenizer(training)
	sentences := tokenizer.Tokenize(string(text))
	for _, s := range sentences {
		fmt.Printf("%q\n", s)
	}
}
