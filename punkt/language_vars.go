package punkt

import (
	"bytes"
	"regexp"
	"text/template"
)

var ReNonPunct = regexp.MustCompile(`[^\W\d]`)

// Format of a regular expression to split punctuation from words, excluding period.
const wordTokenFmt string = `(
{{.MultiChar}}
|
(?={{.WordStart}})\S+?				  # Accept word characters until end is found
(?=									  # Sequences marking a word's end
\s|                                   # White-space
$|									  # End-of-string
{{.NonWord}}|{{.MultiChar}}|          # Punctuation
,(?=$|\s|{{.NonWord}}|{{.MultiChar}}) # Comma if at end of word
)
|
\S
)`

type wordTokenStruct struct {
	MultiChar, WordStart, NonWord string
}

/*
Format of a regular expression to find contexts including possible
sentence boundaries. Matches token which the possible sentence boundary
ends, and matches the following token within a lookahead expression
*/
const periodContextFmt string = `
\S*                         # some word material
{{.SentEndChars}}           # a potential sentence ending
(?=(?P<after_tok>
{{.NonWord}}                # either other punctuation
|
\s+(?P<next_tok>\S+)        # or whitespace and some other token
))`

type periodContextStruct struct {
	SentEndChars []byte
	NonWord      string
}

type PunktLanguageVars struct {
	sentEndChars          []byte         // Characters that are candidates for sentence boundaries
	internalPunctuation   string         // Sentence internal punctuation, which indicates an abbreviation if preceded by a period-final token
	reBoundaryRealignment *regexp.Regexp // Used to realign punctuation that should be included in a sentence although it follows the period (or ?, !)
	reWordStart           string         // Excludes some characters from starting word tokens
	reNonWordChars        string         // Characters that cannot appear within words
	reMultiCharPunct      string         // Hyphen and ellipsis are multi-character punctuation
	wordTokenizeFmt       string
	periodContextFmt      string
}

func NewPunktLanguageVars() *PunktLanguageVars {
	return &PunktLanguageVars{
		sentEndChars:          []byte{'.', '?', '!'},
		internalPunctuation:   ",:;",
		reBoundaryRealignment: regexp.MustCompile(`["\')\]}]+?(?:\s+|(?=--)|$)`),
		reWordStart:           "[^\\(\"\\`{\\[:;&\\#\\*@\\)}\\]\\-,]",
		reNonWordChars:        `(?:[?!)\";}\]\*:@\'\({\[])`,
		reMultiCharPunct:      `(?:\-{2,}|\.{2,}|(?:\.\s){2,}\.)`,
		wordTokenizeFmt:       wordTokenFmt,
	}
}

// Compile word tokenizer regexp
func (p *PunktLanguageVars) ReWordTokenizer() *regexp.Regexp {
	t := template.Must(template.New("wordTokenizer").Parse(p.wordTokenizeFmt))
	var r bytes.Buffer

	t.Execute(&r, wordTokenStruct{
		MultiChar: p.reMultiCharPunct,
		WordStart: p.reWordStart,
		NonWord:   p.reNonWordChars,
	})

	return regexp.MustCompile(r.String())
}

// Tokenize a string to split off punctuation other than periods
func (p *PunktLanguageVars) WordTokenize(s string) []string {
	return p.ReWordTokenizer().FindAllString(s, -1)
}

// Compile period context regexp
func (p *PunktLanguageVars) RePeriodContext() *regexp.Regexp {
	t := template.Must(template.New("periodContext").Parse(p.periodContextFmt))
	var r bytes.Buffer

	t.Execute(&r, periodContextStruct{
		SentEndChars: p.sentEndChars,
		NonWord:      p.reNonWordChars,
	})

	return regexp.MustCompile(r.String())
}

// Compiles and returns a regular expression to find contexts including possible sentence boundaries.
func (p *PunktLanguageVars) PeriodContext(s string) []string {
	return p.RePeriodContext().FindAllString(s, -1)
}

func (p *PunktLanguageVars) ReSentEndChars() string {
	return regexp.QuoteMeta(string(p.sentEndChars))
}
