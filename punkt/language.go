package punkt

import (
	"bytes"
	"regexp"
	"strings"
	"text/template"
)

var ReNonPunct = regexp.MustCompile(`[^\W\d]`)

/*
Format of a regular expression to find contexts including possible
sentence boundaries. Matches token which the possible sentence boundary
ends, and matches the following token within a lookahead expression
*/
const periodContextFmt string = `\S*{{.SentEndChars}}(?P<after_tok>{{.NonWord}}|\s+(?P<next_tok>\S+))`

type periodContextStruct struct {
	SentEndChars string
	NonWord      string
}

type Language struct {
	sentEndChars        []byte // Characters that are candidates for sentence boundaries
	internalPunctuation string // Sentence internal punctuation, which indicates an abbreviation if preceded by a period-final token
	reWordStart         string // Excludes some characters from starting word tokens
	reNonWordChars      string // Characters that cannot appear within words
	periodContextFmt    string
}

func NewLanguage() *Language {
	return &Language{
		sentEndChars:        []byte{'.', '?', '!'},
		internalPunctuation: ",:;",
		reWordStart:         "[^\\(\"\\`{\\[:;&\\#\\*@\\)}\\]\\-,]",
		reNonWordChars:      `(?:[?!)’”"';}\]\*:@\'\({\[])`,
		periodContextFmt:    periodContextFmt,
	}
}

// Compile period context regexp
func (p *Language) RePeriodContext() *regexp.Regexp {
	t := template.Must(template.New("periodContext").Parse(p.periodContextFmt))
	r := new(bytes.Buffer)

	t.Execute(r, periodContextStruct{
		SentEndChars: strings.Join([]string{`[`, p.ReSentEndChars(), `][’”"']?`}, ""),
		NonWord:      p.reNonWordChars,
	})

	return regexp.MustCompile(strings.Trim(r.String(), " "))
}

// Compiles and returns a regular expression to find contexts including possible sentence boundaries.
func (p *Language) PeriodContext(s string) []string {
	return p.RePeriodContext().FindAllString(s, -1)
}

func (p *Language) ReSentEndChars() string {
	return regexp.QuoteMeta(string(p.sentEndChars))
}
