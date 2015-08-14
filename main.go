package main

import (
	//	"bufio"
	"fmt"
	"github.com/neurosnap/go-sentences/punkt"
	"io/ioutil"
	//	"os"
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
	text := `Here... are some initials E.R.B. and also an etc. in the middle.
Periods that form part of an abbreviation but are taken to be end-of-sentence markers
or vice versa do not only introduce errors in the determination of sentence boundaries.
Segmentation errors propagate into further components which rely on accurate
sentence segmentation and subsequent analyses are most likely affected negatively.
Walker et al. (2001), for example, stress the importance of correct sentence boundary
disambiguation for machine translation and Kiss and Strunk (2002b) show that errors
in sentence boundary detection lead to a higher error rate in part-of-speech tagging.
In this paper, we present an approach to sentence boundary detection that builds
on language-independent methods and determines sentence boundaries with high accuracy.
It does not make use of additional annotations, part-of-speech tagging, or precompiled
lists to support sentence boundary detection but extracts all necessary data
from the corpus to be segmented. Also, it does not use orthographic information as primary
evidence and is thus suited to process single-case text. It focuses on robustness
and flexibility in that it can be applied with good results to a variety of languages without
any further adjustments. At the same time, the modular structure of the proposed
system makes it possible in principle to integrate language-specific methods and clues
to further improve its accuracy. The basic algorithm has been determined experimentally
on the basis of an unannotated development corpus of English. We have applied
the resulting system to further corpora of English text as well as to corpora from ten
other languages: Brazilian Portuguese, Dutch, Estonian, French, German, Italian, Norwegian,
Spanish, Swedish, and Turkish. Without further additions or amendments to
the system produced through experimentation on the development corpus, the mean
accuracy of sentence boundary detection on newspaper corpora in eleven languages is
98.74 %.`

	b, err := ioutil.ReadFile("data/english.json")
	if err != nil {
		panic(err)
	}

	params, err := Load(b)
	if err != nil {
		panic(err)
	}

	tokenizer := punkt.NewSentenceTokenizer(params)
	//reader := bufio.NewReader(os.Stdin)
	//contents, _ := ioutil.ReadAll(reader)
	sentences := tokenizer.Tokenize(text)
	for _, s := range sentences {
		fmt.Println(s)
		fmt.Println("------")
	}
}
