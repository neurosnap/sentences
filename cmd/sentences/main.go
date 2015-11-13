package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/neurosnap/sentences/english"
)

var VERSION string

func main() {
	reader := bufio.NewReader(os.Stdin)
	text, _ := ioutil.ReadAll(reader)

	tokenizer := english.NewSentenceTokenizer(nil)
	sentences := tokenizer.Tokenize(string(text))
	for _, s := range sentences {
		fmt.Printf("%q\n", s)
	}
}
