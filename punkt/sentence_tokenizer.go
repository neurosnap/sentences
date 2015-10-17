package punkt

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

// Sane defaults for the sentence tokenizer
func NewSentenceTokenizer(s *Storage) *DefaultSentenceTokenizer {
	lang := NewLanguage()
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

func (s *DefaultSentenceTokenizer) Tokenize(text string) []string {
	// Use the default word tokenizer but only grab the tokens that
	// relate to a sentence ending punctuation.  This means grab the word
	// before and after the punctuation.
	tokens := s.WordTokenizer.Tokenize(text, true)

	lastBreak := 0
	// Think of AnnotateTokens as a pipeline that we send our tokens through to process.
	annotatedTokens := s.AnnotateTokens(tokens, s.Annotations...)
	sentences := make([]string, 0, len(annotatedTokens))
	for _, token := range annotatedTokens {
		if !token.SentBreak {
			continue
		}

		sentence := text[lastBreak:token.Position]
		sentences = append(sentences, sentence)

		lastBreak = token.Position
	}

	sentences = append(sentences, text[lastBreak:])
	return sentences
}
