package punkt

import (
	"regexp"
	"strings"
	"unicode"
)

type PairToken interface {
	PairTokens([]*DefaultToken) [][2]*DefaultToken
}

type DefaultPairToken struct{}

// Groups two adjacent tokens together.
func (p *DefaultPairToken) PairTokens(tokens []*DefaultToken) [][2]*DefaultToken {
	pairTokens := make([][2]*DefaultToken, 0, len(tokens))

	prevToken := tokens[0]
	for _, tok := range tokens {
		if prevToken == tok {
			continue
		}
		pairTokens = append(pairTokens, [2]*DefaultToken{prevToken, tok})
		prevToken = tok
	}
	pairTokens = append(pairTokens, [2]*DefaultToken{prevToken, nil})

	return pairTokens
}

// Primary token interface that determines the context and type of a tokenized word.
type Token interface {
	GetType(string) string
	TypeNoPeriod() string
	TypeNoSentPeriod() string
	FirstCase() string
	FirstLower() bool
	FirstUpper() bool
	IsAlpha() bool
	IsEllipsis() bool
	IsInitial() bool
	IsNumber() bool
	IsNonPunct() bool
}

// Stores a token of text with annotations produced during sentence boundary detection.
type DefaultToken struct {
	Token
	*Language
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

func NewToken(token string) *DefaultToken {
	tok := DefaultToken{
		Tok:        token,
		Language:   NewLanguage(),
		reEllipsis: regexp.MustCompile(`\.\.+$`),
		reNumeric:  regexp.MustCompile(`-?[\.,]?\d[\d,\.-]*\.?$`),
		reInitial:  regexp.MustCompile(`^[A-Za-z]\.$`),
		reInitials: regexp.MustCompile(`[A-Za-z]\.[A-Za-z]\.$`),
		reAlpha:    regexp.MustCompile(`^[A-Za-z]+$`),
		Ellipsis:   false,
	}
	tok.Typ = tok.GetType(token)
	tok.PeriodFinal = strings.HasSuffix(token, ".")
	tok.Token = &tok

	return &tok
}

// Returns a case-normalized representation of the token.
func (p *DefaultToken) GetType(tok string) string {
	return p.reNumeric.ReplaceAllString(strings.ToLower(tok), "##number##")
}

// The type with its final period removed if it has one.
func (p *DefaultToken) TypeNoPeriod() string {
	if len(p.Typ) > 1 && string(p.Typ[len(p.Typ)-1]) == "." {
		return string(p.Typ[:len(p.Typ)-1])
	}
	return p.Typ
}

// The type with its final period removed if it is marked as a sentence break.
func (p *DefaultToken) TypeNoSentPeriod() string {
	if p == nil {
		return ""
	}

	if p.SentBreak {
		return p.TypeNoPeriod()
	}
	return p.Typ
}

// True if the token's first character is uppercase.
func (p *DefaultToken) FirstUpper() bool {
	if p == nil || p.Tok == "" {
		return false
	}

	runes := []rune(p.Tok)
	return unicode.IsUpper(runes[0])
}

// True if the token's first character is lowercase
func (p *DefaultToken) FirstLower() bool {
	if p.Tok == "" {
		return false
	}

	runes := []rune(p.Tok)
	return unicode.IsLower(runes[0])
}

// Text based case of the first letter in the token.
func (p *DefaultToken) FirstCase() string {
	if p.FirstLower() {
		return "lower"
	} else if p.FirstUpper() {
		return "upper"
	}
	return "none"
}

// True if the token text is that of an ellipsis.
func (p *DefaultToken) IsEllipsis() bool {
	return p.reEllipsis.MatchString(p.Tok)
}

// True if the token text is that of a number.
func (p *DefaultToken) IsNumber() bool {
	return strings.HasPrefix(p.Tok, "##number##")
}

// True if the token text is that of an initial.
func (p *DefaultToken) IsInitial() bool {
	return p.reInitial.MatchString(p.Tok)
}

// True if the token text is all alphabetic.
func (p *DefaultToken) IsAlpha() bool {
	return p.reAlpha.MatchString(p.Tok)
}

// True if the token is either a number or is alphabetic.
func (p *DefaultToken) IsNonPunct() bool {
	nonPunct := regexp.MustCompile(p.NonPunct())
	return nonPunct.MatchString(p.Typ)
}
