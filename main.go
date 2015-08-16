package main

import (
	"bufio"
	"fmt"
	"github.com/neurosnap/sentences/punkt"
	"io/ioutil"
	"os"
)

type CustomTokenizer struct {
	*punkt.SentenceTokenizer
}

func (c *CustomTokenizer) AnnotateTokens(tokens []*punkt.Token) {
	fmt.Println("HI MOM")
	tokens = c.SentenceTokenizer.AnnotateTokens(tokens)
}

func NewCustomTokenizer(trainedData *punkt.Storage) *CustomTokenizer {
	st := &CustomTokenizer{
		punkt.NewSentenceTokenizer(trainedData),
	}

	return st
}

func main() {

	b, err := ioutil.ReadFile("data/english.json")
	if err != nil {
		panic(err)
	}

	training, err := punkt.LoadTraining(b)
	if err != nil {
		panic(err)
	}

	tokenizer := NewCustomTokenizer(training)

	tokenizer.AbbrevTypes.Add("al")
	tokenizer.AbbrevTypes.Add("etc")

	reader := bufio.NewReader(os.Stdin)
	text, _ := ioutil.ReadAll(reader)

	sentences := tokenizer.Tokenize(string(text))
	for _, s := range sentences {
		fmt.Println(s)
		fmt.Println("------")
	}
}
