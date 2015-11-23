package sentences

import "fmt"

// Interface used by the Tokenize function, can be extended to correct sentence
// boundaries that punkt misses.
type SentenceTokenizer interface {
	AnnotateTokens([]*Token, ...AnnotateTokens) []*Token
	Tokenize(string) []*Sentence
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

// Sane defaults for the sentence tokenizer
func NewSentenceTokenizer(s *Storage) *DefaultSentenceTokenizer {
	lang := NewPunctStrings()
	word := NewWordTokenizer(lang)

	annotations := NewAnnotations(s, lang, word)

	tokenizer := &DefaultSentenceTokenizer{
		Storage:       s,
		PunctStrings:  lang,
		WordTokenizer: word,
		Annotations:   annotations,
	}

	return tokenizer

}

// Wraps around DST doing the work for customizing the tokenizer
func NewTokenizer(s *Storage, word WordTokenizer, lang PunctStrings) *DefaultSentenceTokenizer {
	annotations := NewAnnotations(s, lang, word)

	tokenizer := &DefaultSentenceTokenizer{
		Storage:       s,
		PunctStrings:  lang,
		WordTokenizer: word,
		Annotations:   annotations,
	}

	return tokenizer
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

/*
Fully annotated word tokens.  This allows for adhoc adjustments to the tokens
*/
func (s *DefaultSentenceTokenizer) AnnotatedTokens(text string) []*Token {
	// Use the default word tokenizer but only grab the tokens that
	// relate to a sentence ending punctuation.  This means grab the word
	// before and after the punctuation.
	tokens := s.WordTokenizer.Tokenize(text, true)

	if len(tokens) == 0 {
		return nil
	}

	return s.AnnotateTokens(tokens, s.Annotations...)
}

/*
Instead of returning an array of sentences, this method returns only the positions
for the sentence boundaries.
*/
func (s *DefaultSentenceTokenizer) SentencePositions(text string) []int {
	annotatedTokens := s.AnnotatedTokens(text)

	positions := make([]int, 0, len(annotatedTokens))
	for _, token := range annotatedTokens {
		if !token.SentBreak {
			continue
		}

		positions = append(positions, token.Position)
	}

	lastChar := len(text)
	positions = append(positions, lastChar)

	return positions
}

/*
Container to hold sentences, provides the character positions
as well as the text for that sentence.
*/
type Sentence struct {
	Start int    `json:"start"`
	End   int    `json:"end"`
	Text  string `json:"text"`
}

func (s Sentence) String() string {
	return fmt.Sprintf("<Sentence [%d:%d]>", s.Start, s.End)
}

func (s *DefaultSentenceTokenizer) Tokenize(text string) []*Sentence {
	annotatedTokens := s.AnnotatedTokens(text)

	lastBreak := 0
	sentences := make([]*Sentence, 0, len(annotatedTokens))
	for _, token := range annotatedTokens {
		if !token.SentBreak {
			continue
		}

		sentence := &Sentence{lastBreak, token.Position, text[lastBreak:token.Position]}
		sentences = append(sentences, sentence)

		lastBreak = token.Position
	}

	lastChar := len(text)
	sentence := &Sentence{lastBreak, lastChar, text[lastBreak:lastChar]}
	sentences = append(sentences, sentence)

	return sentences
}
