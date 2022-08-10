[![release](https://github.com/neurosnap/sentences/actions/workflows/release.yml/badge.svg)](https://github.com/neurosnap/sentences/actions/workflows/release.yml)
[![GODOC](https://godoc.org/github.com/nathany/looper?status.svg)](https://godoc.org/github.com/neurosnap/sentences)
![MIT](https://img.shields.io/packagist/l/doctrine/orm.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/neurosnap/sentences)](https://goreportcard.com/report/github.com/neurosnap/sentences)

# Sentences - A command line sentence tokenizer

This command line utility will convert a blob of text into a list of sentences.

* [Demo](https://sentences-231000.appspot.com/)
* [Docs](https://godoc.org/github.com/neurosnap/sentences)

## Features

* Supports multiple languages (english, czech, dutch, estonian, finnish,
  german, greek, italian, norwegian, polish, portuguese, slovene, and turkish)
* Zero dependencies
* Extendable
* Fast

## Install

### arch

[aur](https://aur.archlinux.org/packages/sentences-bin)

### mac

```
brew tap neurosnap/sentences
brew install sentences
```

### other

Or you can find the pre-built binaries on [the github
releases page](https://github.com/neurosnap/sentences/releases).

### using golang

```
go get github.com/neurosnap/sentences
go install github.com/neurosnap/sentences/cmd/sentences
```

## Command

![Command line](sentences.gif?raw=true)

## Get it

```
go get github.com/neurosnap/sentences
```

## Use it

```Go
import (
    "fmt"
    "os"

    "github.com/neurosnap/sentences"
)

func main() {
    text := `A perennial also-ran, Stallings won his seat when longtime lawmaker David Holmes
    died 11 days after the filing deadline. Suddenly, Stallings was a shoo-in, not
    the long shot. In short order, the Legislature attempted to pass a law allowing
    former U.S. Rep. Carolyn Cheeks Kilpatrick to file; Stallings challenged the
    law in court and won. Kilpatrick mounted a write-in campaign, but Stallings won.`

    // download the training data from this repo (./data) and save it somewhere
    b, _ := os.ReadFile("./path/to/english.json")

    // load the training data
    training, _ := sentences.LoadTraining(b)

    // create the default sentence tokenizer
    tokenizer := sentences.NewSentenceTokenizer(training)
    sentences := tokenizer.Tokenize(text)

    for _, s := range sentences {
        fmt.Println(s.Text)
    }
}
```

## English

This package attempts to fix some problems I noticed for english.

```Go
import (
    "fmt"

    "github.com/neurosnap/sentences/english"
)

func main() {
    text := "Hi there. Does this really work?"

    tokenizer, err := english.NewSentenceTokenizer(nil)
    if err != nil {
        panic(err)
    }

    sentences := tokenizer.Tokenize(text)
    for _, s := range sentences {
        fmt.Println(s.Text)
    }
}
```

## Contributing

I need help maintaining this library.  If you are interested in contributing
to this library then please start by looking at the [golden-rules](https://github.com/neurosnap/sentences/tree/golden-rule) branch which
tests the [Golden Rules](https://github.com/diasks2/pragmatic_segmenter/blob/master/README.md#the-golden-rules)
for english sentence tokenization created by the [Pragmatic Segmenter](https://github.com/diasks2/pragmatic_segmenter)
library.

Create an issue for a particular failing test and submit an issue/PR.

I'm happy to help anyone willing to contribute.

## Customize

`sentences` was built around composability, most major components of this package
can be extended.

Eager to make ad-hoc changes but don't know how to start?
Have a look at `github.com/neurosnap/sentences/english` for a solid example.

## Notice

I have not tested this tokenizer in any other language besides English.  By default
the command line utility loads english. I welcome anyone willing to test the
other languages to submit updates as needed.

A primary goal for this package is to be multilingual so I'm willing to help in
any way possible.

This library is a port of the [nltk's](http://www.nltk.org) punkt tokenizer.

## A Punkt Tokenizer

An unsupervised multilingual sentence boundary detection library for golang.
The way the punkt system accomplishes this goal is through training the tokenizer
with text in that given language.  Once the likelihoods of abbreviations, collocations,
and sentence starters are determined, finding sentence boundaries becomes easier.

There are many problems that arise when tokenizing text into sentences, the primary
issue being abbreviations.  The punkt system attempts to determine whether a  word
is an abbreviation, an end to a sentence, or even both through training the system with text
in the given language.  The punkt system incorporates both token- and type-based
analysis on the text through two different phases of annotation.

[Unsupervised multilingual sentence boundary detection](http://citeseerx.ist.psu.edu/viewdoc/download;jsessionid=BAE5C34E5C3B9DC60DFC4D93B85D8BB1?doi=10.1.1.85.5017&rep=rep1&type=pdf)

## Performance

Using [Brown Corpus](http://www.hit.uib.no/icame/brown/bcm.html) which is annotated American English
text, we compare this package with other libraries across multiple programming languages.

|Library    | Avg Speed (s, 10 runs) | Accuracy (%)
|:----------|:----------------------:|:-----------:
| Sentences | 1.96                   | 98.95
| NLTK      | 5.22                   | 99.21
