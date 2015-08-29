package punkt

import (
	"regexp"
	"strings"
	"unicode"
)

type WToken interface {
	GetType(string) string
	TypeNoPeriod() string
	TypeNoSentPeriod() string
	FirstCase() bool
	FirstLower() bool
	FirstUpper() bool
	IsAlpha() bool
	IsEllipsis() bool
	IsInitial() bool
	IsNumber() bool
	IsNonPunct() bool
}

// Stores a token of text with annotations produced during sentence boundary detection.
type Token struct {
	WToken
	reEllipsis  *regexp.Regexp
	reNumeric   *regexp.Regexp
	reInitial   *regexp.Regexp
	reInitials  *regexp.Regexp
	reAlpha     *regexp.Regexp
	Tok         string
	Typ         string
	PeriodFinal bool
	SentBreak   bool
	ParaStart   bool
	LineStart   bool
	Ellipsis    bool
	Abbr        bool
}

func NewToken(token string) *Token {
	tok := Token{
		Tok:        token,
		reEllipsis: regexp.MustCompile(`\.\.+$`),
		reNumeric:  regexp.MustCompile(`-?[\.,]?\d[\d,\.-]*\.?$`),
		reInitial:  regexp.MustCompile(`^[A-Za-z]\.$`),
		reInitials: regexp.MustCompile(`[A-Za-z]\.[A-Za-z]\.$`),
		reAlpha:    regexp.MustCompile(`^[A-Za-z]+$`),
		Ellipsis:   false,
	}
	tok.Typ = tok.GetType(token)
	tok.PeriodFinal = strings.HasSuffix(token, ".")

	return &tok
}

// Returns a case-normalized representation of the token.
func (p *Token) GetType(tok string) string {
	return p.reNumeric.ReplaceAllString(strings.ToLower(tok), "##number##")
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

	/*// if the first character in word is not a letter or number,
	// then assume it is not upper cased
	alphanum := regexp.MustCompile(`^\s*[\W]`)
	if alphanum.MatchString(p.Tok) {
		return false
	}

	firstTok := string(p.Tok[0])
	return strings.ToUpper(firstTok) == firstTok*/
}

// True if the token's first character is lowercase
func (p *Token) FirstLower() bool {
	if p.Tok == "" {
		return false
	}

	runes := []rune(p.Tok)
	return unicode.IsLower(runes[0])
	/*	firstTok := string(p.Tok[0])
		return strings.ToLower(firstTok) == firstTok*/
}

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
	return ReNonPunct.MatchString(p.Typ)
}
