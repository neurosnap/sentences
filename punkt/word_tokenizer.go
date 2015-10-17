package punkt

import (
	"regexp"
	"strings"
	"unicode"
)

type WordTokenizer interface {
	TokenParser
	Tokenize(string, bool) []*Token
}

// Helpers to get the type of a token
type TokenType interface {
	Type(*Token) string
	// The type with its final period removed if it has one.
	TypeNoPeriod(*Token) string
	// The type with its final period removed if it is marked as a sentence break.
	TypeNoSentPeriod(*Token) string
}

// Helpers to determine the case of the token's first letter
type TokenFirst interface {
	// True if the token's first character is lowercase
	FirstLower(*Token) bool
	// True if the token's first character is uppercase.
	FirstUpper(*Token) bool
}

// Helpers to determine what type of token we are dealing with.
type TokenExistential interface {
	// True if the token text is all alphabetic.
	IsAlpha(*Token) bool
	// True if the token text is that of an ellipsis.
	IsEllipsis(*Token) bool
	// True if the token text is that of an initial.
	IsInitial(*Token) bool
	// True if the token text is that of a number.
	IsNumber(*Token) bool
	// True if the token is either a number or is alphabetic.
	IsNonPunct(*Token) bool
	// Does this token end with a period?
	HasPeriodFinal(*Token) bool
	// Does this token end with a punctuation and a quote?
	HasSentEndChars(*Token) bool
}

// Primary token interface that determines the context and type of a tokenized word.
type TokenParser interface {
	TokenType
	TokenFirst
	TokenExistential
}

type DefaultWordTokenizer struct {
	PunctStrings
}

func NewWordTokenizer(p PunctStrings) *DefaultWordTokenizer {
	return &DefaultWordTokenizer{p}
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

		hasSentencePunct := p.PunctStrings.HasSentencePunct(word)

		if onlyPeriodContext && !hasSentencePunct && !getNextWord {
			lastSpace = i
			continue
		}

		token := NewToken(word)
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

// Returns a case-normalized representation of the token.
func (p *DefaultWordTokenizer) Type(t *Token) string {
	typ := t.reNumeric.ReplaceAllString(strings.ToLower(t.Tok), "##number##")
	if len(typ) == 1 {
		return typ
	}

	// removing comma from typ
	return strings.Replace(typ, ",", "", -1)
}

// The type with its final period removed if it has one.
func (p *DefaultWordTokenizer) TypeNoPeriod(t *Token) string {
	typ := p.Type(t)
	if len(typ) > 1 && string(typ[len(typ)-1]) == "." {
		return string(typ[:len(typ)-1])
	}
	return typ
}

// The type with its final period removed if it is marked as a sentence break.
func (p *DefaultWordTokenizer) TypeNoSentPeriod(t *Token) string {
	if p == nil {
		return ""
	}

	if t.SentBreak {
		return p.TypeNoPeriod(t)
	}

	return p.Type(t)
}

// True if the token's first character is uppercase.
func (p *DefaultWordTokenizer) FirstUpper(t *Token) bool {
	if t.Tok == "" {
		return false
	}

	runes := []rune(t.Tok)
	return unicode.IsUpper(runes[0])
}

// True if the token's first character is lowercase
func (p *DefaultWordTokenizer) FirstLower(t *Token) bool {
	if t.Tok == "" {
		return false
	}

	runes := []rune(t.Tok)
	return unicode.IsLower(runes[0])
}

// True if the token text is that of an ellipsis.
func (p *DefaultWordTokenizer) IsEllipsis(t *Token) bool {
	return t.reEllipsis.MatchString(t.Tok)
}

// True if the token text is that of a number.
func (p *DefaultWordTokenizer) IsNumber(t *Token) bool {
	return strings.HasPrefix(t.Tok, "##number##")
}

// True if the token text is that of an initial.
func (p *DefaultWordTokenizer) IsInitial(t *Token) bool {
	return t.reInitial.MatchString(t.Tok)
}

// True if the token text is all alphabetic.
func (p *DefaultWordTokenizer) IsAlpha(t *Token) bool {
	return t.reAlpha.MatchString(t.Tok)
}

// True if the token is either a number or is alphabetic.
func (p *DefaultWordTokenizer) IsNonPunct(t *Token) bool {
	nonPunct := regexp.MustCompile(p.PunctStrings.NonPunct())
	return nonPunct.MatchString(p.Type(t))
}

func (p *DefaultWordTokenizer) HasPeriodFinal(t *Token) bool {
	return strings.HasSuffix(t.Tok, ".")
}

// Find any punctuation excluding the period final
func (p *DefaultWordTokenizer) HasSentEndChars(t *Token) bool {
	enders := []string{
		`."`, `.'`, `.)`, /*`.’`, `.”`,*/
		`?`, `?"`, `?'`, `?)`, /*`?’`, `?”`,*/
		`!`, `!"`, `!'`, `!)`, `!’`, `!”`,
	}

	for _, ender := range enders {
		if strings.HasSuffix(t.Tok, ender) {
			return true
		}
	}

	parens := []string{
		`.[`, `.(`, `."`, `.'`,
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
