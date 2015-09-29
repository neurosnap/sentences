package punkt

import (
	"strings"
	"unicode"
)

type WordTokenizer interface {
	Tokenize(string) []*Token
}

type DefaultWordTokenizer struct {
	PunctStrings
}

func (p *DefaultWordTokenizer) Tokenize(text string) []*Token {
	words := strings.Split(text, " ")
	tokens := make([]*Token, 0, len(words))

	paragraphStart := false
	lineStart := false
	count := 0
	for _, word := range words {
		if word == "" {
			count += 1
			continue
		}

		// check if this word starts with a newline
		if strings.HasPrefix(word, "\n") {
			if strings.Count(word, "\n") > 1 || lineStart {
				paragraphStart = true
			}

			lineStart = true
		}

		multWord := strings.Fields(word)
		for i, mult := range multWord {
			if i != 0 {
				lineStart = true
				for _, char := range text[count : count+len(multWord)] {
					if !unicode.IsSpace(char) {
						break
					}

					count += 1
					if count > 1 {
						paragraphStart = true
					}
				}
			}

			token := NewToken(mult, p.PunctStrings)
			token.Position = count
			token.ParaStart = paragraphStart
			token.LineStart = lineStart

			tokens = append(tokens, token)

			lineStart = false
			paragraphStart = false
			count += len(mult)
		}

		// check if next word starts with a newline
		if strings.HasSuffix(word, "\n") {
			lineStart = true
			if strings.Count(word, "\n") > 1 {
				paragraphStart = true
			}
		} else {
			lineStart = false
		}

		count += 1
	}

	return tokens
}
