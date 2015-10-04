package punkt

import (
	"strings"
	"unicode"
)

type WordTokenizer interface {
	Tokenize(string, bool) []*Token
}

type DefaultWordTokenizer struct {
	PunctStrings
}

func (p *DefaultWordTokenizer) Tokenize(text string, onlyPeriodContext bool) []*Token {
	tokens := make([]*Token, 0, 50)
	lastSpace := 0
	lineStart := false
	paragraphStart := false
	getNextWord := false

	for i := 0; i < len(text); i++ {
		char := rune(text[i])
		if !unicode.IsSpace(char) {
			continue
		}

		if char == '\n' {
			if lineStart {
				paragraphStart = true
			}
			lineStart = true
		}

		word := strings.TrimSpace(text[lastSpace:i])
		if word == "" {
			continue
		}

		hasSentencePunct := p.HasSentencePunct(word)

		if onlyPeriodContext && !hasSentencePunct && !getNextWord {
			lastSpace = i
			continue
		}

		token := NewToken(word, p.PunctStrings)
		token.Position = i
		token.ParaStart = paragraphStart
		token.LineStart = lineStart
		tokens = append(tokens, token)

		lastSpace = i
		lineStart = false
		paragraphStart = false

		if hasSentencePunct {
			getNextWord = true
		} else {
			getNextWord = false
		}
	}

	return tokens
}
