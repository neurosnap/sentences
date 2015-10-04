package punkt

import (
	"regexp"
	"strings"
	//"time"
)

// Interface used by the Tokenize function, can be extended to correct sentence
// boundaries that punkt misses.
type SentenceTokenizer interface {
	AnnotateTokens([]*Token, ...AnnotateTokens) []*Token
	Tokenize(string) []string
}

// A sentence tokenizer which uses an unsupervised algorithm to build a model
// for abbreviation words, collocations, and words that start sentences
// and then uses that model to find sentence boundaries.
type DefaultSentenceTokenizer struct {
	*Storage
	WordTokenizer
	PunctStrings
	Annotations []AnnotateTokens
}

func NewTokenizer(s *Storage) *DefaultSentenceTokenizer {
	lang := NewLanguage()

	annotations := []AnnotateTokens{
		&TypeBasedAnnotation{s, lang},
		&TokenBasedAnnotation{s, lang, &DefaultTokenGrouper{}},
	}

	tokenizer := &DefaultSentenceTokenizer{
		Storage:       s,
		PunctStrings:  lang,
		WordTokenizer: &DefaultWordTokenizer{lang},
		Annotations:   annotations,
	}

	return tokenizer
}

func (s *DefaultSentenceTokenizer) Tokenize(text string) []string {
	// Use the default word tokenizer but only grab the tokens that
	// relate to a sentence ending punctuation.  This means grab the word
	// before and after the punctuation.
	wordTokens := s.WordTokenizer.Tokenize(text, true)

	tokens := make([]*Token, 0, len(wordTokens))
	for _, token := range wordTokens {
		// split a token by special punctuation
		// TODO: get rid of this bloated piece of shit
		splitTokens := s.splitToken(token)
		if splitTokens == nil {
			continue
		}

		tokens = append(tokens, splitTokens...)
	}

	lastBreak := 0
	annotatedTokens := s.AnnotateTokens(tokens, s.Annotations...)
	sentences := make([]string, 0, len(annotatedTokens))
	for _, token := range annotatedTokens {
		if token.SentBreak {
			sentence := text[lastBreak:token.Position]
			sentence = strings.TrimSpace(sentence)
			if sentence == "" {
				continue
			}

			sentences = append(sentences, sentence)
			lastBreak = token.Position
		}
	}

	sentences = append(sentences, text[lastBreak:])
	return sentences
}

/*
Given a set of tokens augmented with markers for line-start and
paragraph-start, returns an iterator through those tokens with full
annotation including predicted sentence breaks.
*/
func (s *DefaultSentenceTokenizer) AnnotateTokens(tokens []*Token, annotate ...AnnotateTokens) []*Token {
	for _, ann := range annotate {
		tokens = ann.Annotate(tokens)
	}

	return tokens
}

func (s *DefaultSentenceTokenizer) splitToken(token *Token) []*Token {
	word := token.Tok
	endPuncts := []string{":", ",", "?", `?"`, ".)"}
	nonword := regexp.MustCompile(strings.Join([]string{s.NonWordChars(), s.MultiCharPunct()}, "|"))
	multi := regexp.MustCompile(s.MultiCharPunct())

	if len(word) == 1 {
		return nil
	}

	chars := []rune(word)

	first := word
	second := ""
	for _, punct := range endPuncts {
		if strings.HasSuffix(word, punct) {
			if len(punct) > 1 {
				first = string(chars[:len(chars)-2])
				second = string(chars[len(chars)-2:])
			} else {
				first = string(chars[:len(chars)-1])
				second = string(chars[len(chars)-1:])
			}
		}
	}

	multipunct := multi.FindStringIndex(word)
	if multipunct != nil {
		if strings.HasSuffix(word, ".") && (multipunct[1] != len(word) ||
			multipunct[0]+multipunct[1] == len(word)) {
			first = word[:len(chars)-1]
			second = "."
		} else {
			if multipunct[1] == len(word) {
				first = word[:multipunct[0]]
				second = word[multipunct[0]:]
			} else {
				first = word[:multipunct[1]]
				second = word[multipunct[1]:]
			}
		}
	}

	tokens := make([]*Token, 0, 2)
	if nonword.MatchString(second) || strings.HasSuffix(second, ",") {
		token.Tok = first
		token.Typ = token.GetType(first)

		secondToken := NewToken(second, s.PunctStrings)
		secondToken.Position = token.Position

		token.Position = token.Position - len(second)

		tokens = append(tokens, token, secondToken)
	} else {
		token.Tok = word
		token.Typ = token.GetType(word)
		tokens = append(tokens, token)
	}

	return tokens
}
