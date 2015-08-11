package main

import (
	"fmt"
	"github.com/neurosnap/go-sentences/punkt"
	"io/ioutil"
)

func main() {
	fmt.Println("TRAIN SOME SHIT")

	//fmt.Println(text)
	fmt.Println("-----")

	/*base := punkt.NewPunktBase()
	words := base.TokenizeWords(text)
	for _, word := range words {
		fmt.Println(word.Tok, word.ParaStart, word.LineStart)
	}*/
	/*trainer := punkt.NewPunktTrainer("", os.Stdin)
	results := trainer.GetParams()
	fmt.Println(results.AbbrevTypes)
	fmt.Println(results.Collocations)
	fmt.Println(results.SentStarters)
	fmt.Println(results.OrthoContext)*/

	b, err := ioutil.ReadFile("data/english.json")
	if err != nil {
		panic(err)
	}

	params, err := Load(b)
	if err != nil {
		panic(err)
	}

	sentences := punkt.NewSentenceTokenizer(params)
	fmt.Println(sentences)
}
