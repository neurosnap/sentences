package punkt

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

type WordTokenizer interface {
	Tokenize(string) []*Token
}

type DefaultWordTokenizer struct {
	PunctStrings
}

// Helps tokenize the body of text into words while also providing context for
// the words and how they are used in the text.
/*func (p *DefaultWordTokenizer) Tokenize(text string) []*Token {
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
}*/

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

	finalTokens := make([]*Token, 0, len(tokens))
	for _, token := range tokens {
		splitTokens := p.splitToken(token)
		if splitTokens == nil {
			continue
		}

		finalTokens = append(finalTokens, splitTokens...)
	}

	return finalTokens
}

func (p *DefaultWordTokenizer) splitToken(token *Token) []*Token {
	word := strings.Fields(token.Tok)[0]
	endPuncts := []string{":", ",", "?", `?"`, ".)"}
	nonword := regexp.MustCompile(strings.Join([]string{p.NonWordChars(), p.MultiCharPunct()}, "|"))
	multi := regexp.MustCompile(p.MultiCharPunct())

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
		secondToken := NewToken(second, p.PunctStrings)
		tokens = append(tokens, token, secondToken)
	} else {
		token.Tok = word
		token.Typ = token.GetType(word)
		tokens = append(tokens, token)
	}

	return tokens
}

// temporary helper struct, not particularly useful
type pairTokens struct {
	First, Second string
}

func (p *pairTokens) String() string {
	return fmt.Sprintf("First: %q, Second: %q", p.First, p.Second)
}

// Adds a token to our list of tokens and provides some context for the token.
// Is the token a non-word multipuntuation character or does it end in a comma?
// Does the token start a paragraph?  Does it start a new line?
func (p *DefaultWordTokenizer) addToken(tokens []*Token, pairTok *pairTokens, parastart bool, linestart bool) []*Token {
	nonword := regexp.MustCompile(strings.Join([]string{p.NonWordChars(), p.MultiCharPunct()}, "|"))
	tok := strings.Join([]string{pairTok.First, pairTok.Second}, "")

	if nonword.MatchString(pairTok.Second) || strings.HasSuffix(pairTok.Second, ",") {
		tokOne := NewToken(pairTok.First, p.PunctStrings)
		tokOne.ParaStart = parastart
		tokOne.LineStart = linestart

		tokTwo := NewToken(pairTok.Second, p.PunctStrings)

		tokens = append(tokens, tokOne, tokTwo)
	} else {
		token := NewToken(tok, p.PunctStrings)
		token.ParaStart = parastart
		token.LineStart = linestart
		tokens = append(tokens, token)
	}

	return tokens
}

// Breaks the text up into words and also splits the token into two distinct
// pieces that assist in determining what type of token we are dealing with
func (p *DefaultWordTokenizer) pairWordTokenizer(text string) []*pairTokens {
	endPuncts := []string{":", ",", "?", `?"`, ".)"}
	words := strings.Fields(text)
	tokens := make([]*pairTokens, 0, len(words))

	multi := regexp.MustCompile(p.MultiCharPunct())

	for _, word := range words {
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

		token := &pairTokens{first, second}
		tokens = append(tokens, token)
	}

	return tokens
}
