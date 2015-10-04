package punkt

import (
	"regexp"
	"strings"
	"time"
)

type Tokenizer interface {
	WordTokenizer
	SentenceTokenizer
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

// Interface used by the Tokenize function, can be extended to correct sentence
// boundaries that punkt misses.
type SentenceTokenizer interface {
	PeriodCtxTokenizer(string, WordTokenizer) []*PeriodCtx
	HasSentBreak(string, WordTokenizer) bool
	AnnotateTokens([]*Token, ...AnnotateTokens) []*Token
}

type PeriodCtx struct {
	// Entire context of the period, including word before and after
	Context string
	// Last character in sentence
	End int
}

func (s *DefaultSentenceTokenizer) NTokenize(text string) []string {
	timeit := time.Now()
	wordTokens := s.WordTokenizer.Tokenize(text, true)
	elapsedTime := time.Now().Sub(timeit)

	logger.Println(elapsedTime)

	tokens := make([]*Token, 0, len(wordTokens))
	for _, token := range wordTokens {
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
		//logger.Println(token)
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
Breaks text into sentences using the SentenceTokenizer interface
*/
func Tokenize(text string, t Tokenizer) []string {
	matches := t.PeriodCtxTokenizer(text, t)

	sentences := make([]string, 0, len(matches))
	lastBreak := 0
	for _, match := range matches {
		//logger.Println(match.Context)
		//logger.Println(t.HasSentBreak(match.Context, t))
		if t.HasSentBreak(match.Context, t) {
			sentence := text[lastBreak:match.End]
			sentence = strings.TrimSpace(sentence)
			if sentence == "" {
				continue
			}

			sentences = append(sentences, sentence)
			lastBreak = match.End
		}
	}

	sentences = append(sentences, text[lastBreak:])
	return sentences
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

// Discovers all periods within a body of text, captures the context
// in which it is used, and determines if a period denotes a sentence break.
func (s *DefaultSentenceTokenizer) PeriodCtxTokenizer(text string, w WordTokenizer) []*PeriodCtx {
	rePeriodContext := regexp.MustCompile(s.PeriodContext())
	matches := rePeriodContext.FindAllStringSubmatchIndex(text, -1)
	periodMatches := make([]*PeriodCtx, 0, len(matches))

	/*
	 * match = [15, 23, 20, 23, 21, 23]
	 * entire match = 0:1
	 * second token = 2:3
	 * newlines + second token = 4:5
	 */
	for _, match := range matches {
		context := text[match[0]:match[1]]
		matchEnd := 0

		nextTok := ""
		if match[4] != -1 && match[5] != -1 {
			nextTok = text[match[4]:match[5]]
		}

		matchEnd = match[1]
		// we want the extra stuff for the actual sentence
		if match[4] >= 0 && (!s.HasSentBreak(nextTok, w) || s.HasSentBreak(text[match[0]:match[4]], w)) {
			matchEnd = match[4]
		}

		periodCtx := &PeriodCtx{
			Context: context,
			End:     matchEnd,
		}

		periodMatches = append(periodMatches, periodCtx)
	}

	return periodMatches
}

/*
Returns True if the given text includes a sentence break.
*/
func (s *DefaultSentenceTokenizer) HasSentBreak(text string, w WordTokenizer) bool {
	word_tokens := w.Tokenize(text, true)

	tokens := make([]*Token, 0, len(word_tokens))
	for _, token := range word_tokens {
		splitTokens := s.splitToken(token)
		if splitTokens == nil {
			continue
		}

		tokens = append(tokens, splitTokens...)
	}

	if len(tokens) == 0 {
		return false
	}

	for _, t := range s.AnnotateTokens(tokens, s.Annotations...) {
		if t.SentBreak {
			return true
		}
	}

	return false
}

func (s *DefaultSentenceTokenizer) splitToken(token *Token) []*Token {
	word := strings.TrimSpace(token.Tok)
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
