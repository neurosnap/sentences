package punkt

import (
	"bytes"
	"regexp"
	"strings"
	"text/template"
)

type PunctStrings interface {
	SentEndChars() string
	NonWordChars() string
	PeriodContext() string
	NonPunct() string
	Punctuation() string
	MultiCharPunct() string
	HasSentencePunct(string) bool
}

/*
Format of a regular expression to find contexts including possible
sentence boundaries. Matches token which the possible sentence boundary
ends, and matches the following token within a lookahead expression
*/
type periodContextStruct struct {
	SentEndChars string
	NonWord      string
}

// Language holds language specific regular expressions to help determine
// information about the text that is being parsed.
type Language struct{}

// Creates a default set of properties for the Language struct
func NewLanguage() *Language {
	return &Language{}
}

// Characters that are candidates for sentence boundaries
func (p *Language) SentEndChars() string {
	return `.?!".".'?".)`
}

// Characters that cannot appear within words
func (p *Language) NonWordChars() string {
	return `(?:[?!)’”"';}\]\*:@\'\({\[])`
}

func (p *Language) NonPunct() string {
	return `[^\W\d]`
}

func (p *Language) Punctuation() string {
	return ";:,.!?"
}

func (p *Language) HasSentencePunct(text string) bool {
	endPunct := `.!?`
	for _, char := range endPunct {
		for _, achar := range text {
			if char == achar {
				return true
			}
		}
	}

	return false
}

func (p *Language) MultiCharPunct() string {
	return `(?:\-{2,}|\.{2,}|(?:\.\s){2,}\.)|(\.\S)`
}

// Compile the context of a period context using a regular expression.
// To determine a sentence boundary, punkt must have information about the
// context in which a period is used.
func (p *Language) PeriodContext() string {
	periodContextFmt := `\S*{{.SentEndChars}}(?P<after_tok>{{.NonWord}}|\s+(?P<next_tok>\S+))`
	sentEndChars := regexp.QuoteMeta(p.SentEndChars())

	t := template.Must(template.New("periodContext").Parse(periodContextFmt))
	r := new(bytes.Buffer)

	t.Execute(r, periodContextStruct{
		SentEndChars: strings.Join([]string{`[`, sentEndChars, `][’”"']?`}, ""),
		NonWord:      p.NonWordChars(),
	})

	return strings.Trim(r.String(), " ")
}
