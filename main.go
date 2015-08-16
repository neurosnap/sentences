package main

import (
	//	"bufio"
	"fmt"
	"github.com/neurosnap/sentences/punkt"
	"io/ioutil"
	//"os"
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
	/*text := `Here... are some initials E.R.B. and also an etc. in the middle.
	Periods that form part of an abbreviation but are taken to be end-of-sentence markers
	or vice versa do not only introduce errors in the determination of sentence boundaries.
	What is funny is I grew up in the U.S. Segmentation errors propagate into further components which rely on accurate
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
	 98.74%. However, Mr. T. Pain is a great lyricist, but not a good person.`*/
	text := `
	    Things are hopping in Lansing, as the details of an affair between state Reps.
	    Todd Courser and Cindy Gamrat continue to unfold. With a faux-smear
	    campaign/cover-up claiming Courser hired male prostitutes,
	    orchestrated by Courser himself to “inoculate the herd” against the truth of
	    his affair with Gamrat, a series of purported blackmail texts from
	    “the Lansing Mafia” and ongoing administrative (and possibly criminal)
	    investigations, it’s the most excitement the Capitol has seen in a dog’s age.

	    But Courser and Gamrat should take comfort: Theirs isn’t the only scandal in
	    Michigan political history — and it surely won’t be the last.

	    We here at A Better Michigan worked tirelessly to bring you this by-no-means
	    exhaustive list of Michigan political scandals, burning the midnight oil,
	    thumbing ink-stained accounts from days of yore. Your Kwame Kilpatricks
	    and other well-trod scandalous ground are not for us. We’re bringing you the
	    obtuse and the bizarre — and sometimes, the beneficial policy changes
	    prompted by scandalous behavior. It was painstaking work. Oh, who are
	    we kidding? We love this stuff. Sit back and enjoy.

	    A perennial also-ran, Stallings won his seat when longtime lawmaker David Holmes
	    died 11 days after the filing deadline. Suddenly, Stallings was a shoo-in, not
	    the long shot. In short order, the Legislature attempted to pass a law allowing
	    former U.S. Rep. Carolyn Cheeks Kilpatrick to file; Stallings challenged the
	    law in court and won. Kilpatrick mounted a write-in campaign, but Stallings won.
	    `
	fmt.Println(text)
	b, err := ioutil.ReadFile("data/english.json")
	if err != nil {
		panic(err)
	}

	storage, err := punkt.LoadStorage(b)
	if err != nil {
		panic(err)
	}

	tokenizer := punkt.NewSentenceTokenizer(storage)
	tokenizer.AbbrevTypes.Add("al")
	tokenizer.AbbrevTypes.Add("etc")

	/*reader := bufio.NewReader(os.Stdin)
	text, _ := ioutil.ReadAll(reader)*/

	sentences := tokenizer.Tokenize(string(text))
	for _, s := range sentences {
		fmt.Println(s)
		fmt.Println("------")
	}
}
