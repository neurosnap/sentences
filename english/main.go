package english

import (
	"strings"

	"github.com/neurosnap/sentences/punkt"
)

type WordTokenizer struct {
	punkt.DefaultWordTokenizer
}

func NewSentenceTokenizer(s *punkt.Storage) *punkt.DefaultSentenceTokenizer {
	lang := punkt.NewLanguage()
	wordTokenizer := NewWordTokenizer(lang)
	return punkt.NewTokenizer(s, wordTokenizer, lang)
}

func NewWordTokenizer(p punkt.PunctStrings) *WordTokenizer {
	word := &WordTokenizer{}
	word.PunctStrings = p

	return word
}

// Find any punctuation excluding the period final
func (e *WordTokenizer) HasSentEndChars(t *punkt.Token) bool {
	println("I GOT HIT!")
	enders := []string{
		`."`, `.'`, `.)`, `.’`, `.”`,
		`?`, `?"`, `?'`, `?)`, `?’`, `?”`,
		`!`, `!"`, `!'`, `!)`, `!’`, `!”`,
	}

	for _, ender := range enders {
		if strings.HasSuffix(t.Tok, ender) {
			return true
		}
	}

	parens := []string{
		`.[`, `.(`, `."`, `.'`,
		`?[`, `?(`,
		`![`, `!(`,
	}

	for _, paren := range parens {
		if strings.Index(t.Tok, paren) != -1 {
			return true
		}
	}

	return false
}
