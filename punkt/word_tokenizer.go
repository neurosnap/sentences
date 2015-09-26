package punkt

import (
	"regexp"
	"strings"
)

type WordTokenizer interface {
	Tokenize(string) []*Token
}

type DefaultWordTokenizer struct {
	PunctStrings
}

// Helps tokenize the body of text into words while also providing context for
// the words and how they are used in the text.
func (p *DefaultWordTokenizer) Tokenize(text string) []*Token {
	lines := strings.Split(text, "\n")
	tokens := make([]*Token, 0, len(lines))
	parastart := false

	for _, line := range lines {
		if strings.Trim(line, " ") == "" || line == " " {
			parastart = true
		} else {
			lineToks := p.pairWordTokenizer(line)
			for index, lineTok := range lineToks {
				if index == 0 {
					tokens = p.addToken(tokens, lineTok, parastart, true)
					parastart = false
				} else {
					tokens = p.addToken(tokens, lineTok, parastart, false)
				}
			}
		}
	}

	return tokens
}

// temporary helper struct, not particularly useful
type pairTokens struct {
	First, Second string
}

// Adds a token to our list of tokens and provides some context for the token.
// Is the token a non-word multipuntuation character or does it end in a comma?
// Does the token start a paragraph?  Does it start a new line?
func (p *DefaultWordTokenizer) addToken(tokens []*Token, pairTok *pairTokens, parastart bool, linestart bool) []*Token {
	nonword := regexp.MustCompile(strings.Join([]string{p.NonWordChars(), p.MultiCharPunct()}, "|"))
	tok := strings.Join([]string{pairTok.First, pairTok.Second}, "")

	if nonword.MatchString(pairTok.Second) || strings.HasSuffix(pairTok.Second, ",") {
		tokOne := NewToken(pairTok.First)
		tokOne.ParaStart = parastart
		tokOne.LineStart = linestart

		tokTwo := NewToken(pairTok.Second)

		tokens = append(tokens, tokOne, tokTwo)
	} else {
		token := NewToken(tok)
		token.ParaStart = parastart
		token.LineStart = linestart
		tokens = append(tokens, token)
	}

	return tokens
}

// Breaks the text up into words and also splits the token into two distinct
// pieces that assist in determining what type of token we are dealing with
func (p *DefaultWordTokenizer) pairWordTokenizer(text string) []*pairTokens {
	endPuncts := []string{ /*`."`, `.'`, `.”`,*/ ":", ",", "?", `?"`, ".)"}
	words := strings.Fields(text)
	tokens := make([]*pairTokens, 0, len(words))

	multi := regexp.MustCompile(p.MultiCharPunct())

	for _, word := range words {
		// Skip one letter words
		if len(word) == 1 {
			continue
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
			if strings.HasSuffix(word, ".") && (multipunct[1] != len(word) || multipunct[0]+multipunct[1] == len(word)) {
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

		token := &pairTokens{first, second}
		tokens = append(tokens, token)
	}

	return tokens
}
