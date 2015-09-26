package punkt

import (
	"regexp"
	"strings"
	"unicode"
)

// Groups two adjacent tokens together.
type TokenGrouper interface {
	Group([]*Token) [][2]*Token
}

type DefaultTokenGrouper struct{}

func (p *DefaultTokenGrouper) Group(tokens []*Token) [][2]*Token {
	pairTokens := make([][2]*Token, 0, len(tokens))

	prevToken := tokens[0]
	for _, tok := range tokens {
		if prevToken == tok {
			continue
		}
		pairTokens = append(pairTokens, [2]*Token{prevToken, tok})
		prevToken = tok
	}
	pairTokens = append(pairTokens, [2]*Token{prevToken, nil})

	return pairTokens
}

// Helpers to get the type of a token
type TokenType interface {
	// The type with its final period removed if it has one.
	TypeNoPeriod() string
	// The type with its final period removed if it is marked as a sentence break.
	TypeNoSentPeriod() string
}

// Helpers to determine the case of the token's first letter
type TokenFirst interface {
	// Text based case of the first letter in the token.
	FirstCase() string
	// True if the token's first character is lowercase
	FirstLower() bool
	// True if the token's first character is uppercase.
	FirstUpper() bool
}

// Helpers to determine what type of token we are dealing with.
type TokenExistential interface {
	// True if the token text is all alphabetic.
	IsAlpha() bool
	// True if the token text is that of an ellipsis.
	IsEllipsis() bool
	// True if the token text is that of an initial.
	IsInitial() bool
	// True if the token text is that of a number.
	IsNumber() bool
	// True if the token is either a number or is alphabetic.
	IsNonPunct() bool
	// Does this token end with a period?
	HasPeriodFinal() bool
}

// Primary token interface that determines the context and type of a tokenized word.
type TokenParser interface {
	TokenType
	TokenFirst
	TokenExistential
}

// Stores a token of text with annotations produced during sentence boundary detection.
type Token struct {
	TokenParser
	PunctStrings
	Tok         string
	Typ         string
	SentBreak   bool
	ParaStart   bool
	LineStart   bool
	Abbr        bool
	periodFinal bool
	reEllipsis  *regexp.Regexp
	reNumeric   *regexp.Regexp
	reInitial   *regexp.Regexp
	reAlpha     *regexp.Regexp
}

func NewToken(token string) *Token {
	tok := Token{
		Tok:          token,
		PunctStrings: NewLanguage(),
		reEllipsis:   regexp.MustCompile(`\.\.+$`),
		reNumeric:    regexp.MustCompile(`-?[\.,]?\d[\d,\.-]*\.?$`),
		reInitial:    regexp.MustCompile(`^[A-Za-z]\.$`),
		reAlpha:      regexp.MustCompile(`^[A-Za-z]+$`),
	}
	tok.periodFinal = strings.HasSuffix(token, ".")
	tok.TokenParser = &tok

	return &tok
}

// The type with its final period removed if it has one.
func (p *Token) TypeNoPeriod() string {
	if len(p.Typ) > 1 && string(p.Typ[len(p.Typ)-1]) == "." {
		return string(p.Typ[:len(p.Typ)-1])
	}
	return p.Typ
}

// The type with its final period removed if it is marked as a sentence break.
func (p *Token) TypeNoSentPeriod() string {
	if p == nil {
		return ""
	}

	if p.SentBreak {
		return p.TypeNoPeriod()
	}

	return p.Typ
}

// True if the token's first character is uppercase.
func (p *Token) FirstUpper() bool {
	if p == nil || p.Tok == "" {
		return false
	}

	runes := []rune(p.Tok)
	return unicode.IsUpper(runes[0])
}

// True if the token's first character is lowercase
func (p *Token) FirstLower() bool {
	if p.Tok == "" {
		return false
	}

	runes := []rune(p.Tok)
	return unicode.IsLower(runes[0])
}

// Text based case of the first letter in the token.
func (p *Token) FirstCase() string {
	if p.FirstLower() {
		return "lower"
	} else if p.FirstUpper() {
		return "upper"
	}

	return "none"
}

// True if the token text is that of an ellipsis.
func (p *Token) IsEllipsis() bool {
	return p.reEllipsis.MatchString(p.Tok)
}

// True if the token text is that of a number.
func (p *Token) IsNumber() bool {
	return strings.HasPrefix(p.Tok, "##number##")
}

// True if the token text is that of an initial.
func (p *Token) IsInitial() bool {
	return p.reInitial.MatchString(p.Tok)
}

// True if the token text is all alphabetic.
func (p *Token) IsAlpha() bool {
	return p.reAlpha.MatchString(p.Tok)
}

// True if the token is either a number or is alphabetic.
func (p *Token) IsNonPunct() bool {
	nonPunct := regexp.MustCompile(p.NonPunct())
	return nonPunct.MatchString(p.Typ)
}

func (p *Token) HasPeriodFinal() bool {
	return p.periodFinal
}
