Sentences - A command line sentence tokenizer written in Golang
===============================================================

This command line utility will convert a blob of text into a list of sentences.

* [Demo](http://sentences.erock.io)
* [Docs](https://godoc.org/gopkg.in/neurosnap/sentences.v1)

Notice
------

I have not tested this tokenizer in any other language besides English.  I
welcome anyone willing to test the other languages to submit updates as needed.

This library is a port of the [nltk's](http://www.nltk.org) punkt tokenizer.

Install
-------

```
go get gopkg.in/neurosnap/sentences.v1
go install gopkg.in/neurosnap/sentences.v1/cmd/sentences
```

### Binaries

#### Linux

* [Linux 386](https://s3-us-west-2.amazonaws.com/sentence-binaries/sentences_linux-386.tar.gz)
* [Linux AMD64](https://s3-us-west-2.amazonaws.com/sentence-binaries/sentences_linux-amd64.tar.gz)

#### Mac

* [Darwin 386](https://s3-us-west-2.amazonaws.com/sentence-binaries/sentences_darwin-386.tar.gz)
* [Darwin AMD64](https://s3-us-west-2.amazonaws.com/sentence-binaries/sentences_darwin-amd64.tar.gz)

#### Windows

* [Windows 386](https://s3-us-west-2.amazonaws.com/sentence-binaries/sentences_windows-386.tar.gz)
* [Windows AMD64](https://s3-us-west-2.amazonaws.com/sentence-binaries/sentences_windows-amd64.tar.gz)

Command-Line
------------

![Command line](sentences.gif?raw=true)

Get it
------

```
go get gopkg.in/neurosnap/sentences.v1
```

Use it
------

```
import (
    "fmt"

    "gopkg.in/neurosnap/sentences.v1"
    "gopkg.in/neurosnap/sentences.v1/data"
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
    training, _ := sentences.LoadTraining(data)

    // create the default sentence tokenizer
    tokenizer := sentences.NewSentenceTokenizer(training)
    sentences := tokenizer.Tokenize(text)
    for _, s := range sentences {
        fmt.Println(s.Text)
    }
}
```

English
-------

This package attempts to fix some problems I noticed for english.

```
import (
    "fmt"

    "gopkg.in/neurosnap/sentences.v1"
)

func main() {
    text := "Hi there. Cool."

    tokenizer := english.NewSentenceTokenizer(nil)
    sentences := tokenizer.Tokenize(text)
    for _, s := range sentences {
        fmt.Println(s.Text)
    }
}
```

A Punkt Tokenizer
-----------------

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

[Unsupervised multilingual sentence boundary detection](http://citeseerx.ist.psu.edu/viewdoc/download;jsessionid=BAE5C34E5C3B9DC60DFC4D93B85D8BB1?doi=10.1.1.85.5017&rep=rep1&type=pdf)

