package english

import (
	"regexp"
	"strings"

	"gopkg.in/neurosnap/sentences.v1"
	"gopkg.in/neurosnap/sentences.v1/data"
)

type WordTokenizer struct {
	sentences.DefaultWordTokenizer
}

var reAbbr = regexp.MustCompile(`((?:[\w]\.)+[\w]*\.)`)

// English customized sentence tokenizer.
func NewSentenceTokenizer(s *sentences.Storage) (*sentences.DefaultSentenceTokenizer, error) {
	training := s

	if training == nil {
		b, err := data.Asset("data/english.json")
		if err != nil {
			return nil, err
		}

		training, err = sentences.LoadTraining(b)
		if err != nil {
			return nil, err
		}
	}

	// supervisor abbreviations
	abbrevs := []string{"sgt", "gov", "no", "mt"}
	for _, abbr := range abbrevs {
		training.AbbrevTypes.Add(abbr)
	}

	lang := sentences.NewPunctStrings()
	word := NewWordTokenizer(lang)
	annotations := sentences.NewAnnotations(training, lang, word)

	ortho := &sentences.OrthoContext{
		Storage:      training,
		PunctStrings: lang,
		TokenType:    word,
		TokenFirst:   word,
	}

	multiPunct := &MultiPunctWordAnnotation{
		Storage:      training,
		TokenParser:  word,
		TokenGrouper: &sentences.DefaultTokenGrouper{},
		Ortho:        ortho,
	}

	annotations = append(annotations, multiPunct)

	tokenizer := &sentences.DefaultSentenceTokenizer{
		Storage:       training,
		PunctStrings:  lang,
		WordTokenizer: word,
		Annotations:   annotations,
	}

	return tokenizer, nil
}

func NewWordTokenizer(p sentences.PunctStrings) *WordTokenizer {
	word := &WordTokenizer{}
	word.PunctStrings = p

	return word
}

// Find any punctuation excluding the period final
func (e *WordTokenizer) HasSentEndChars(t *sentences.Token) bool {
	enders := []string{
		`."`, `.)`, `.’`, `.”`,
		`?`, `?"`, `?'`, `?)`, `?’`, `?”`,
		`!`, `!"`, `!'`, `!)`, `!’`, `!”`,
	}

	for _, ender := range enders {
		if strings.HasSuffix(t.Tok, ender) {
			return true
		}
	}

	parens := []string{
		`.[`, `.(`, `."`,
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

// Attempts to tease out custom Abbreviations, e.g. F.B.I.
type MultiPunctWordAnnotation struct {
	*sentences.Storage
	sentences.TokenParser
	sentences.TokenGrouper
	sentences.Ortho
}

func (a *MultiPunctWordAnnotation) Annotate(tokens []*sentences.Token) []*sentences.Token {
	for _, tokPair := range a.TokenGrouper.Group(tokens) {
		if len(tokPair) < 2 || tokPair[1] == nil {
			continue
		}

		a.tokenAnnotation(tokPair[0], tokPair[1])
	}

	return tokens
}

// looksInternal determines if tok's punctuation could appear
// sentence-internally (i.e., parentheses or quotations).
func looksInternal(tok string) bool {
	internal := []string{")", `’`, `”`, `"`, `'`}
	for _, punc := range internal {
		if strings.HasSuffix(tok, punc) {
			return true
		}
	}
	return false
}

func (a *MultiPunctWordAnnotation) tokenAnnotation(tokOne, tokTwo *sentences.Token) {
	nextTyp := a.TokenParser.TypeNoSentPeriod(tokTwo)

	/*
		If the tokOne's sentence-breaking punctuation looks like it could occur
		sentence-internally, ensure that the following word is either
		capitalized or a frequent sentence starter.
	*/
	if tokOne.SentBreak && looksInternal(tokOne.Tok) {
		if a.TokenParser.FirstLower(tokTwo) && a.SentStarters[nextTyp] == 0 {
			tokOne.SentBreak = false
			return
		}
	}

	if len(reAbbr.FindAllString(tokOne.Tok, 1)) == 0 {
		return
	}

	if a.IsInitial(tokOne) {
		return
	}

	tokOne.Abbr = true
	tokOne.SentBreak = false

	/*
		[4.1.1. Orthographic Heuristic] Check if there's
		orthogrpahic evidence about whether the next word
		starts a sentence or not.
	*/
	isSentStarter := a.Ortho.Heuristic(tokTwo)
	if isSentStarter == 1 {
		tokOne.SentBreak = true
		return
	}

	/*
		[4.1.3. Frequent Sentence Starter Heruistic] If the
		next word is capitalized, and is a member of the
		frequent-sentence-starters list, then label tok as a
		sentence break.
	*/
	if a.TokenParser.FirstUpper(tokTwo) && a.SentStarters[nextTyp] != 0 {
		tokOne.SentBreak = true
		return
	}
}
