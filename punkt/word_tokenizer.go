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

func (p *DefaultWordTokenizer) Tokenize(text string, onlyPunctuation bool) []*Token {
	tokens := []*Token{}
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

		if onlyPunctuation && !hasSentencePunct && !getNextWord {
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

func (p *DefaultWordTokenizer) OTokenize(text string, onlyPunctuation bool) []*Token {
	words := strings.Split(text, " ")
	tokens := make([]*Token, 0, len(words))

	paragraphStart := false
	lineStart := false
	getNextWord := false
	count := 0
	for _, word := range words {
		if word == "" {
			count += 1
			continue
		}

		if onlyPunctuation && !p.HasSentencePunct(word) && !getNextWord {
			count += len(word) + 1
			getNextWord = false
			continue
		}

		getNextWord = true

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

			count += len(mult)
			//logger.Println(count, mult)
			token := NewToken(mult, p.PunctStrings)
			token.Position = count
			token.ParaStart = paragraphStart
			token.LineStart = lineStart

			tokens = append(tokens, token)

			lineStart = false
			paragraphStart = false
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
