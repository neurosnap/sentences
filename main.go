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

type englishAnnotation struct{}

func (e *englishAnnotation) Annotate(tokens []*punkt.Token) []*punkt.Token {
	fmt.Println("MADE IT!")
	return tokens
}

func main() {
	b, err := data.Asset("data/english.json")
	if err != nil {
		panic(err)
	}

	training, err := punkt.LoadTraining(b)
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(os.Stdin)
	text, _ := ioutil.ReadAll(reader)

	tokenizer := punkt.NewTokenizer(training)
	tokenizer.Annotations = append(tokenizer.Annotations, &englishAnnotation{})

	sentences := tokenizer.Tokenize(string(text))
	for _, s := range sentences {
		fmt.Printf("%q\n", s)
	}
}
