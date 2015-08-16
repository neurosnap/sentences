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
in the given language.

Original Research Article
-------------------------

[Unsupervised multilingual sentence boundary detection](http://citeseerx.ist.psu.edu/viewdoc/download;jsessionid=BAE5C34E5C3B9DC60DFC4D93B85D8BB1?doi=10.1.1.85.5017&rep=rep1&type=pdf)

Notice
------

I have not tested this tokenizer in any other language besides English.  I
welcome anyone willing to test the other languages to submit updates as needed.

This library is a port of the [http://www.nltk.org](nltk's) punkt tokenizer.

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
    "github.com/neurosnape/sentences/punkt"
)

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
thumbing ink-stained accounts from days of yore. ... Your Kwame Kilpatricks
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
data, err := ioutil.ReadFile("data/english.json")
if err != nil {
    panic(err)
}

storage, err := punkt.LoadStorage(data)
tokenizer := punkt.NewSentenceTokenizer(storage)
sentences := tokenizer.Tokenize(text)

for _, s := range sentences {
    fmt.Println(s)
}
```

Got some abbreviations you want to add to the list?
```
tokenizer := punkt.NewSentenceTokenizer(storage)

tokenizer.AbbrevTypes.Add("al")
tokenizer.AbbrevTypes.Add("etc")

sentences := tokenizer.Tokenize(text)
```
