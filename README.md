Sentences - Punkt Sentence Tokenizer
====================================

An unsupervised multilingual sentence boundary detection library for golang.
The goal of this library is to be able to break up any text into a list of sentences
in multiple languages without creating special heuristics for any language in particular.
The way the punkt system accomplishes this goal is through training the tokenizer
with text in that given language.  Once the likelyhoods of abbreviations, collocations,
and sentence starters are determined, finding sentence boundaries becomes easier.

There are many problems that arise when tokenizing text into sentences, the primary
issue being abbreviations.  The punkt system attempts to determine whether a  word
is an abbrevation, an end to a sentence, or even both through training the system with text
in the given language.  The punkt system incorporates both token- and type-based
analysis on the text through two different phases of annotation.

Original Research Article
-------------------------

[Unsupervised multilingual sentence boundary detection](http://citeseerx.ist.psu.edu/viewdoc/download;jsessionid=BAE5C34E5C3B9DC60DFC4D93B85D8BB1?doi=10.1.1.85.5017&rep=rep1&type=pdf)

Notice
------

I have not tested this tokenizer in any other language besides English.  I
welcome anyone willing to test the other languages to submit updates as needed.

This library is a port of the [nltk's](http://www.nltk.org) punkt tokenizer.

Command line
------------

Currently very simple, it takes `stdin` and outputs a sentence on each line.
The default command line utility pre-loads the english training data, so loading
it is not necessary.

```
go install github.com/neurosnap/sentences
```


Get it
------

```
go get github.com/neurosnap/sentences
```

Use it
------

```
import (
    "fmt"
    "io/ioutil"
    "github.com/neurosnap/sentences/punkt"
)

func main() {
    text := `A perennial also-ran, Stallings won his seat when longtime lawmaker David Holmes
    died 11 days after the filing deadline. Suddenly, Stallings was a shoo-in, not
    the long shot. In short order, the Legislature attempted to pass a law allowing
    former U.S. Rep. Carolyn Cheeks Kilpatrick to file; Stallings challenged the
    law in court and won. Kilpatrick mounted a write-in campaign, but Stallings won.`

    // English data is loaded in by default
    // Compiling language specific data into a binary file can be accomplished
    // by using `make <lang>` and then using the same method below
    b, _ := data.Asset("data/english.json");

    // load the training data
    training, _ := punkt.LoadTraining(data)

    // create the default sentence tokenizer
    tokenizer := punkt.NewSentenceTokenizer(training)
    sentences := punkt.Tokenize(text, tokenizer)

    for _, s := range sentences {
        fmt.Println(s)
    }
}
```

Got some abbreviations you want to add to the list?
```
tokenizer := punkt.NewSentenceTokenizer(storage)

tokenizer.AbbrevTypes.Add("al")
tokenizer.AbbrevTypes.Add("etc")

sentences := Tokenize(text, tokenizer)
```

Want to extend the tokenizer?  In my mind there is one method in particular
that will yield a ton of extendability: `AnnotateTokens`

This method conducts both the first and second pass annotations that determine
abbreviations, period context, and sentence boundaries.

```
type CustomSentenceTokenizer struct {
    *punkt.DefaultSentenceTokenizer
}

func (s *CustomSentenceTokenizer) AnnotateTokens(tokens []*punkt.DefaultToken) {
    tokens = s.AnnotateFirstPass(tokens)
    tokens = s.AnnotateSecondPass(tokens)

    // Do a third pass and find any sentence boundaries that were missed by punkt

    return tokens
}

tokenizer := &CustomSentenceTokenize{
    &punkt.DefaultSentenceTokenizer{
        Base: punkt.NewBase(),
        Punctuation: punkt.Punctuation,
    },
}

tokenizer.Storage = training
tokenizer.SentenceTokenizer = tokenizer

punkt.Tokenize(text, tokenizer)
```
