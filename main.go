package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/neurosnap/sentences/data"
	"github.com/neurosnap/sentences/punkt"
)

var VERSION string

func main() {

	b, err := data.Asset("data/english.json")
	//b, err := ioutil.ReadFile("data/english.json")
	if err != nil {
		panic(err)
	}

	training, err := punkt.LoadTraining(b)
	if err != nil {
		panic(err)
	}

	tokenizer := punkt.NewSentenceTokenizer(training)

	reader := bufio.NewReader(os.Stdin)
	text, _ := ioutil.ReadAll(reader)

	sentences := tokenizer.Tokenize(string(text))
	for _, s := range sentences {
		fmt.Println(s)
	}
}
