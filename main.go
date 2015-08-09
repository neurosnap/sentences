package main

import (
	"fmt"
	"github.com/neurosnap/go-sentences/punkt"
)

func main() {
	fmt.Println("TRAIN SOME SHIT")

	text := `This is an abbr.  Isn't it cool? Here is a list: one, two, three,
	 etc. but that's not all! I know my initials are E.R.B but I don't like
	 herb. I forgot to add a hypen-to-my-text. but1.) is one, 2,0 is two... . , ...high... there`

	//fmt.Println(text)
	fmt.Println("-----")

	lang := punkt.NewPunktLanguageVars()
	words := lang.WordTokenizer(text /*"... ...one two... ...three..."*/)
	for _, word := range words {
		fmt.Println(*word)
	}
	//trainer := punkt.NewPunktTrainer("This is a sentence.")
	//fmt.Println("%v", trainer.GetParams())
}
